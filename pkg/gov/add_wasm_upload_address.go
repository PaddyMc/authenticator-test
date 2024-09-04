package gov

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	//	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v26/app/params"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func AddWasmAddressProposal(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	voteKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	paramChange *proposal.ParameterChangeProposal,
) error {
	priv1 := voteKey
	//priv2 := seedConfig.Keys[1]
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	govClient := govv1beta1.NewQueryClient(conn)
	paramsClient := proposal.NewQueryClient(conn)

	subspaces, err := paramsClient.Subspaces(
		context.Background(),
		&proposal.QuerySubspacesRequest{},
	)
	fmt.Println(subspaces)

	changeParamMsg, err := govv1beta1.NewMsgSubmitProposal(
		paramChange,
		sdk.Coins{sdk.Coin{Denom: "uosmo", Amount: osmomath.NewInt(2000000000)}},
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
		&govv1beta1.QueryProposalsRequest{
			ProposalStatus: 2,
			Depositor:      accAddress.String(),
		},
	)
	latestProposalID := proposals.Proposals[len(proposals.Proposals)-1].ProposalId
	log.Printf("Voting yes on proposal id: %d", latestProposalID)

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

	proposal, err := govClient.Proposal(
		context.Background(),
		&govv1beta1.QueryProposalRequest{
			ProposalId: latestProposalID,
		},
	)
	if proposal.Proposal.Status == 2 {
		log.Printf("Proposal passed id: %d", latestProposalID)
	} else {
		log.Printf("Proposal failed id: %d", latestProposalID)
	}

	return nil
}
