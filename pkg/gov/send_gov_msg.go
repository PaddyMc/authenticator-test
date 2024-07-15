package gov

import (
	"context"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v25/app/params"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func GovMessageProposal(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	voteKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	msgs []sdk.Msg,
) error {
	priv1 := voteKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	govClient := govv1.NewQueryClient(conn)

	govMsg, err := govv1.NewMsgSubmitProposal(
		msgs,
		sdk.Coins{sdk.Coin{Denom: "uosmo", Amount: osmomath.NewInt(500000000)}},
		accAddress.String(),
		"This is a gov message",
		"Adding generic messages for gov execution",
		"Gov me",
		false,
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
		[]sdk.Msg{govMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	proposals, err := govClient.Proposals(
		context.Background(),
		&govv1.QueryProposalsRequest{
			ProposalStatus: 2,
			Depositor:      accAddress.String(),
		},
	)

	latestProposalID := proposals.Proposals[len(proposals.Proposals)-1].Id

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
		&govv1.QueryProposalRequest{
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
