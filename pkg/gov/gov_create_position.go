package gov

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v23/app/params"
	clmodel "github.com/osmosis-labs/osmosis/v23/x/concentrated-liquidity/model"
	cltypes "github.com/osmosis-labs/osmosis/v23/x/concentrated-liquidity/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	poolmanagergrpc "github.com/osmosis-labs/osmosis/v23/x/poolmanager/client/queryproto"
)

func GovCreateClPosition(
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
	bankClient := banktypes.NewQueryClient(conn)
	pmClient := poolmanagergrpc.NewQueryClient(conn)
	//clClient := clproto.NewQueryClient(conn)

	govModuleAcc := "osmo10d07y265gmmuvt4z0w9aw880jnsr700jjeq4qp"

	// Get the pool from the pool id
	resp, err := pmClient.Pool(
		context.Background(),
		&poolmanagergrpc.PoolRequest{
			PoolId: 1251,
		},
	)
	if err != nil {
		return err
	}

	var pool clmodel.Pool
	if err := encCfg.Marshaler.Unmarshal(resp.Pool.Value, &pool); err != nil {
		return err
	}

	balancePost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: govModuleAcc, Denom: pool.Token0},
	)
	if err != nil {
		return err
	}

	//	sendTokenMsg := &banktypes.MsgSend{
	//		FromAddress: accAddress.String(),
	//		ToAddress:   govModuleAcc,
	//		Amount:      sdk.NewCoins(sdk.NewCoin(pool.Token0, osmomath.NewInt(1000000))),
	//	}
	//
	//	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
	//		[]cryptotypes.PrivKey{priv1},
	//		[]cryptotypes.PrivKey{priv1},
	//		nil,
	//		encCfg,
	//		ac,
	//		txClient,
	//		chainID,
	//		[]sdk.Msg{sendTokenMsg},
	//		[]uint64{},
	//	)
	//	if err != nil {
	//		return err
	//	}

	// Need to be divisable by tick spacing
	down := (pool.CurrentTick - (10 * int64(pool.TickSpacing)))
	roundedDown := ((down / 100) + 1) * 100
	up := (pool.CurrentTick + (10 * int64(pool.TickSpacing)))
	roundedUp := ((up / 100) + 1) * 100

	log.Printf("Gov module balance being used for CL position %d:  %s %s\n", pool.Id, pool.Token0, balancePost.GetBalance().Amount.String())

	createCLPositionMsg := &cltypes.MsgCreatePosition{
		PoolId:    pool.Id,
		Sender:    govModuleAcc,
		LowerTick: roundedDown,
		UpperTick: roundedUp,
		TokensProvided: sdk.Coins{
			sdk.NewCoin(pool.Token0, osmomath.NewInt(1000000)),
			// 100 osmo
			sdk.NewCoin(pool.Token1, osmomath.NewInt(1000000)),
		}.Sort(),
		TokenMinAmount0: osmomath.NewInt(0),
		TokenMinAmount1: osmomath.NewInt(0),
	}
	p, err := govv1.NewProposal(
		[]sdk.Msg{createCLPositionMsg},
		1000,
		time.Now(),
		time.Now(),
		"metadata",
		"title",
		"summary",
		accAddress,
		false,
	)

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{createCLPositionMsg},
		sdk.Coins{sdk.Coin{Denom: "uosmo", Amount: sdk.NewInt(2000000000)}},
		accAddress.String(),
		p.Metadata,
		p.Title,
		p.Summary,
		p.Expedited,
	)

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{msg},
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
	fmt.Println(proposal)

	//createTransferPositionMsg := &cltypes.MsgTransferPositions{
	//	PositionIds: []uint64{position.Position.PositionId},
	//	Sender:      accAddress.String(),
	//	NewOwner:    newOwner,
	//}

	return nil
}
