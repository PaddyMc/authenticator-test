package gov

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/osmosis-labs/osmosis/v23/app/params"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func ChangeMaxAuthenticatorGas(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	voteKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	priv1 := voteKey
	//priv2 := seedConfig.Keys[1]
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	govClient := govv1beta1.NewQueryClient(conn)

	// change param for smart contract
	changes := proposal.ParamChange{
		Subspace: "authenticator",
		Key:      "MaximumUnauthenticatedGas",
		Value:    "100000",
	}
	paramChange := proposal.NewParameterChangeProposal(
		"Update the max gas for authenticators",
		"Updating the gas to 100000",
		[]proposal.ParamChange{changes},
	)
	changeParamMsg, err := govv1beta1.NewMsgSubmitProposal(
		paramChange,
		sdk.Coins{sdk.Coin{Denom: "uosmo", Amount: sdk.NewInt(2000000000)}},
		accAddress,
	)
	if err != nil {
		return err
	}
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{changeParamMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}
	proposals, err := govClient.Proposals(
		context.Background(),
		&govv1beta1.QueryProposalsRequest{},
	)
	latestProposalID := proposals.Proposals[len(proposals.Proposals)-1].ProposalId

	voteMsg := govv1beta1.NewMsgVote(accAddress, latestProposalID, govv1beta1.OptionYes)
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{voteMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	return nil
}
