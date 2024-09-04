package cl

import (
	"context"
	"log"
	"time"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v26/app/params"
	clproto "github.com/osmosis-labs/osmosis/v26/x/concentrated-liquidity/client/queryproto"
	clmodel "github.com/osmosis-labs/osmosis/v26/x/concentrated-liquidity/model"
	cltypes "github.com/osmosis-labs/osmosis/v26/x/concentrated-liquidity/types"
	poolmanagergrpc "github.com/osmosis-labs/osmosis/v26/x/poolmanager/client/queryproto"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

// CreateIncentivicedPosition create a position in a given pool and swaps tokens
// in the pool to generate incentives for that position
func CreateIncentivicedPosition(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
	pool clmodel.Pool,
) error {
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	bankClient := banktypes.NewQueryClient(conn)
	clClient := clproto.NewQueryClient(conn)

	accAddress := sdk.AccAddress(fromKey.PubKey().Address())

	log.Printf("Creating position and swapping in pool id %d token0 %s token1 %s\n", pool.Id, pool.Token0, pool.Token1)
	incentiveRecordsResp, err := clClient.IncentiveRecords(
		context.Background(),
		&clproto.IncentiveRecordsRequest{
			PoolId:     pool.Id,
			Pagination: &query.PageRequest{},
		},
	)
	if err != nil {
		return err
	}

	if len(incentiveRecordsResp.IncentiveRecords) < 0 {
		for i, incentiveRecord := range incentiveRecordsResp.IncentiveRecords {
			log.Printf("Incentives Record %d in pool %d: remainingCoin %s emissionRate %s minUptime %s\n",
				i,
				pool.Id,
				incentiveRecord.IncentiveRecordBody.RemainingCoin.String(),
				incentiveRecord.IncentiveRecordBody.EmissionRate.String(),
				incentiveRecord.MinUptime.String(),
			)
		}
	}

	poolAccumulatorRewardsResp, err := clClient.PoolAccumulatorRewards(
		context.Background(),
		&clproto.PoolAccumulatorRewardsRequest{
			PoolId: pool.Id,
		},
	)
	if err != nil {
		return err
	}

	for _, uptimeG := range poolAccumulatorRewardsResp.UptimeGrowthGlobal {
		if uptimeG.UptimeGrowthOutside != nil {
			log.Printf("Incentive accumulmator coins for pool %d: %s\n", pool.Id, uptimeG.UptimeGrowthOutside)
		}
	}

	balancePool, err := bankClient.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{Address: pool.Address},
	)
	if err != nil {
		return err
	}

	// Ensure the pool has liquidity
	if len(balancePool.Balances) < 1 {
		return nil
	}

	err = SwapTokens(
		conn,
		encCfg,
		chainID,
		fromKey,
		fromKey,
		nil,
		selectedAuthenticator,
		pool.Token0,
		pool.Token1,
		pool.Id,
		10000000000,
	)
	if err != nil {
		return err
	}

	balancePost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: pool.Token0},
	)
	if err != nil {
		return err
	}

	// Need to be divisable by tick spacing
	down := (pool.CurrentTick - (10 * int64(pool.TickSpacing)))
	roundedDown := ((down / 100) + 1) * 100
	up := (pool.CurrentTick + (10 * int64(pool.TickSpacing)))
	roundedUp := ((up / 100) + 1) * 100

	log.Printf("Liquidity for CL position %d:  %s%d\n", pool.Id, pool.Token0, balancePost.GetBalance().Amount)

	createCLPositionMsg := &cltypes.MsgCreatePosition{
		PoolId:    pool.Id,
		Sender:    accAddress.String(),
		LowerTick: roundedDown,
		UpperTick: roundedUp,
		TokensProvided: sdk.Coins{
			sdk.NewCoin(pool.Token0, balancePost.GetBalance().Amount),
			// 100 osmo
			sdk.NewCoin(pool.Token1, osmomath.NewInt(1000000)),
		}.Sort(),
		TokenMinAmount0: osmomath.NewInt(0),
		TokenMinAmount1: osmomath.NewInt(0),
	}

	log.Printf("Creating CL position in pool %d: with current tick %d upperTick %d lowerTick %d\n",
		pool.Id,
		pool.CurrentTick,
		up,
		down,
	)

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{fromKey},
		[]cryptotypes.PrivKey{fromKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{createCLPositionMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	err = SwapTokens(
		conn,
		encCfg,
		chainID,
		fromKey,
		fromKey,
		nil,
		selectedAuthenticator,
		pool.Token0,
		pool.Token1,
		pool.Id,
		100000000,
	)
	if err != nil {
		return err
	}

	// Wait incentives to be distributed
	time.Sleep(30 * time.Second)

	// Check rewards for user in pool
	_, err = DisplayRewardsForPosition(
		clClient,
		pool,
		accAddress,
	)
	if err != nil {
		return err
	}

	return nil
}

// DisplayRewardsForPosition TODO: rename to something better, it returns a position id ;(
func DisplayRewardsForPosition(
	clClient clproto.QueryClient,
	pool clmodel.Pool,
	accAddress sdk.Address,
) (clmodel.FullPositionBreakdown, error) {
	resp, err := clClient.UserPositions(
		context.Background(),
		&clproto.UserPositionsRequest{
			Address:    accAddress.String(),
			PoolId:     pool.Id,
			Pagination: &query.PageRequest{},
		},
	)
	position := clmodel.FullPositionBreakdown{}
	if err != nil {
		return position, err
	}
	if len(resp.Positions) > 0 {
		position = resp.Positions[len(resp.Positions)-1]

		log.Printf("Rewards for position in pool %d: user position %d\n",
			pool.Id,
			position.Position.PositionId,
		)

		for _, incentive := range position.ClaimableIncentives {
			log.Printf("Position Claimable Incentives in pool %d: amount %d denom %s\n",
				pool.Id,
				incentive.Amount,
				incentive.Denom,
			)
		}

		for _, incentive := range position.ForfeitedIncentives {
			log.Printf("Position Forfeited Incentives in pool %d: amount %d denom %s\n",
				pool.Id,
				incentive.Amount,
				incentive.Denom,
			)
		}
	}

	return position, nil
}

// CreatesPositionAndTransfer creates a cl position and transfers that cl position to another user
func CreatePositionAndTransfer(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
	poolId uint64,
	tokenAmount int64,
	newOwner string,
) error {
	log.Println("Starting swapping and creating position in single pool...")

	// Set up all grpc clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	clClient := clproto.NewQueryClient(conn)
	pmClient := poolmanagergrpc.NewQueryClient(conn)

	accAddress := sdk.AccAddress(fromKey.PubKey().Address())

	// Get the pool from the pool id
	resp, err := pmClient.Pool(
		context.Background(),
		&poolmanagergrpc.PoolRequest{
			PoolId: poolId,
		},
	)
	if err != nil {
		return err
	}

	var pool clmodel.Pool
	if err := encCfg.Marshaler.Unmarshal(resp.Pool.Value, &pool); err != nil {
		return err
	}

	err = CreateIncentivicedPosition(
		conn,
		encCfg,
		chainID,
		fromKey,
		signerKey,
		cosignersKeys,
		selectedAuthenticator,
		pool,
	)
	if err != nil {
		return err
	}

	// Wait for incentives
	log.Println("Wait for incentives to be distributed...")
	time.Sleep(30 * time.Second)

	// check rewards for user in pool
	position, err := DisplayRewardsForPosition(
		clClient,
		pool,
		accAddress,
	)
	if err != nil {
		return err
	}

	log.Println("Transfering position...")
	createTransferPositionMsg := &cltypes.MsgTransferPositions{
		PositionIds: []uint64{position.Position.PositionId},
		Sender:      accAddress.String(),
		NewOwner:    newOwner,
	}
	_, err = DisplayRewardsForPosition(
		clClient,
		pool,
		accAddress,
	)
	if err != nil {
		return err
	}

	// remove liquidity from the pool
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{fromKey},
		[]cryptotypes.PrivKey{fromKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{createTransferPositionMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	poolAccumulatorRewardsResp, err := clClient.PoolAccumulatorRewards(
		context.Background(),
		&clproto.PoolAccumulatorRewardsRequest{
			PoolId: pool.Id,
		},
	)
	if err != nil {
		return err
	}

	for i, uptimeG := range poolAccumulatorRewardsResp.UptimeGrowthGlobal {
		if uptimeG.UptimeGrowthOutside != nil {
			log.Printf("Incentive accumulmator %d coins for pool %d: %s\n", i, pool.Id, uptimeG.UptimeGrowthOutside)
		}
	}

	return nil
}

func SwapAndCreatePositionInCLPool(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
	poolId uint64,
	tokenAmount int64,
) error {
	log.Println("Starting swapping and creating position in single pool...")

	// Set up all grpc clients
	//	txClient := txtypes.NewServiceClient(conn)
	//	ac := auth.NewQueryClient(conn)
	clClient := clproto.NewQueryClient(conn)
	pmClient := poolmanagergrpc.NewQueryClient(conn)

	accAddress := sdk.AccAddress(fromKey.PubKey().Address())

	resp, err := pmClient.Pool(
		context.Background(),
		&poolmanagergrpc.PoolRequest{
			PoolId: poolId,
		},
	)
	if err != nil {
		return err
	}

	var pool clmodel.Pool
	if err := encCfg.Marshaler.Unmarshal(resp.Pool.Value, &pool); err != nil {
		return err
	}

	err = CreateIncentivicedPosition(
		conn,
		encCfg,
		chainID,
		fromKey,
		signerKey,
		cosignersKeys,
		selectedAuthenticator,
		pool,
	)
	if err != nil {
		return err
	}

	// wait for uptime
	log.Println("Wait for incentives to be distributed...")
	time.Sleep(30 * time.Second)

	// check rewards for user in pool
	_, err = DisplayRewardsForPosition(
		clClient,
		pool,
		accAddress,
	)
	if err != nil {
		return err
	}

	log.Println("Withdrawing position...")
	//	createWithdrawPositionMsg := &cltypes.MsgWithdrawPosition{
	//		PositionId:      position.Position.PositionId,
	//		Sender:          accAddress.String(),
	//		LiquidityAmount: position.Position.Liquidity,
	//	}
	//
	//	// remove liquidity from the pool
	//	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
	//		[]cryptotypes.PrivKey{fromKey},
	//		[]cryptotypes.PrivKey{fromKey},
	//		nil,
	//		encCfg,
	//		ac,
	//		txClient,
	//		chainID,
	//		[]sdk.Msg{createWithdrawPositionMsg},
	//		[]uint64{},
	//	)
	//	if err != nil {
	//		return err
	//	}

	poolAccumulatorRewardsResp, err := clClient.PoolAccumulatorRewards(
		context.Background(),
		&clproto.PoolAccumulatorRewardsRequest{
			PoolId: pool.Id,
		},
	)
	if err != nil {
		return err
	}

	for i, uptimeG := range poolAccumulatorRewardsResp.UptimeGrowthGlobal {
		if uptimeG.UptimeGrowthOutside != nil {
			log.Printf("Incentive accumulator %d coins for pool %d: %s\n", i, pool.Id, uptimeG.UptimeGrowthOutside)
		}
	}

	return nil
}

func SwapAndCreatePositionInAllCLPoolsWithLiquidity(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
) error {
	log.Println("Starting swapping and creating positions in all incentivised pools...")

	clClient := clproto.NewQueryClient(conn)

	limit := uint64(2)
	key := []byte{}
	for {
		resp, err := clClient.Pools(
			context.Background(),
			&clproto.PoolsRequest{
				Pagination: &query.PageRequest{
					Key:        key,
					Limit:      limit,
					CountTotal: true,
				},
			},
		)
		if err != nil {
			return err
		}

		key = resp.Pagination.NextKey

		for _, binaryAny := range resp.Pools {
			var pool clmodel.Pool
			if err := encCfg.Marshaler.Unmarshal(binaryAny.Value, &pool); err != nil {
				return err
			}
			//log.Printf("Iterating pool id %d incentive address %s \n", pool.Id, pool.IncentivesAddress)

			if pool.Token1 != "uosmo" {
				continue
			}

			if pool.CurrentTickLiquidity.Equal(osmomath.ZeroDec()) {
				continue
			}

			incentiveRecordsResp, err := clClient.IncentiveRecords(
				context.Background(),
				&clproto.IncentiveRecordsRequest{
					PoolId:     pool.Id,
					Pagination: &query.PageRequest{},
				},
			)
			if err != nil {
				return err
			}

			if len(incentiveRecordsResp.IncentiveRecords) <= 0 {
				continue
			}

			log.Printf("Pool id %d has liquidity %s and incentive of %s with emission rate of %s \n",
				pool.Id,
				pool.CurrentTickLiquidity,
				incentiveRecordsResp.IncentiveRecords[0].IncentiveRecordBody.RemainingCoin,
				incentiveRecordsResp.IncentiveRecords[0].IncentiveRecordBody.EmissionRate,
			)

			err = SwapAndCreatePositionInCLPool(
				conn,
				encCfg,
				chainID,
				fromKey,
				fromKey,
				nil,
				selectedAuthenticator,
				pool.Id,
				10000000000,
			)

			if err != nil {
				return err
			}

		}

		if len(key) == 0 {
			break
		}

	}

	return nil
}
