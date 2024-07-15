package cl

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/osmosis-labs/osmosis/v25/app/params"
	clproto "github.com/osmosis-labs/osmosis/v25/x/concentrated-liquidity/client/queryproto"
	clmodel "github.com/osmosis-labs/osmosis/v25/x/concentrated-liquidity/model"
)

// CheckAllIncentiveAccumulatorForPools prints all incentive accumulmator info for all cl pools
// and exports it to a csv file
func CheckIncentiveAccumulatorForPools(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	filename string,
) error {
	log.Println("Starting checking all pool incentive accumulators...")

	clClient := clproto.NewQueryClient(conn)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	headerCSV := []string{
		"Pool Id", "Incentive Accumulator",
	}
	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headerCSV); err != nil {
		panic(err)
	}

	//
	sfile, err := os.Create("spread-rewards.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	headerCSV = []string{
		"Pool Id", "Spread Reward Denom", "Spread Reward Amount",
	}
	// Create a CSV writer
	swriter := csv.NewWriter(sfile)
	defer swriter.Flush()

	if err := swriter.Write(headerCSV); err != nil {
		panic(err)
	}

	limit := uint64(100)
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

		for _, binaryAny := range resp.Pools {
			var pool clmodel.Pool
			if err := encCfg.Marshaler.Unmarshal(binaryAny.Value, &pool); err != nil {
				return err
			}

			var incentiveAccumValue string

			poolAccumulatorRewardsResp, err := clClient.PoolAccumulatorRewards(
				context.Background(),
				&clproto.PoolAccumulatorRewardsRequest{
					PoolId: pool.Id,
				},
			)
			if err != nil {
				return err
			}

			var incRecord []string
			var record []string

			for _, uptimeG := range poolAccumulatorRewardsResp.UptimeGrowthGlobal {
				if uptimeG.UptimeGrowthOutside != nil {
					log.Printf("Incentive accumulmator coins for pool %d: %s\n", pool.Id, uptimeG.UptimeGrowthOutside)
					incentiveAccumValue = uptimeG.UptimeGrowthOutside.String()
					incRecord = []string{
						fmt.Sprintf("%d", pool.Id),
						incentiveAccumValue,
					}
					if err := writer.Write(incRecord); err != nil {
						panic(err) // Handle error
					}
				}
			}
			if len(incRecord) == 0 {
				incRecord = []string{
					fmt.Sprintf("%d", pool.Id),
					"0",
				}
				if err := writer.Write(incRecord); err != nil {
					panic(err) // Handle error
				}
			}

			for _, spreadReward := range poolAccumulatorRewardsResp.SpreadRewardGrowthGlobal {
				log.Printf("Spread rewards coins for pool %d: %s %s\n", pool.Id, spreadReward.Denom, spreadReward.Amount.String())
				record = []string{
					fmt.Sprintf("%d", pool.Id),
					spreadReward.Denom,
					spreadReward.Amount.String(),
				}
				if err := swriter.Write(record); err != nil {
					panic(err) // Handle error
				}
			}

			if len(record) == 0 {
				record = []string{
					fmt.Sprintf("%d", pool.Id),
					"0",
				}
				if err := writer.Write(record); err != nil {
					panic(err) // Handle error
				}
			}
		}

		if len(resp.Pagination.NextKey) == 0 {
			break
		}

		key = resp.Pagination.NextKey
	}

	return nil
}

// CheckAllIncentiveForPools prints all incentive records for all cl pools
// and exports it to a csv file
func CheckAllIncentivesForPools(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fileName string,
) error {
	log.Println("Starting checking all pool incentives...")

	clClient := clproto.NewQueryClient(conn)

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	headerCSV := []string{
		"Pool Id", "RecordNumner", "RemainingCoins", "EmissionRate", "MinUptime",
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headerCSV); err != nil {
		panic(err)
	}

	limit := uint64(100)
	key := []byte{}
	for {
		resp, err := clClient.Pools(
			context.Background(),
			&clproto.PoolsRequest{
				Pagination: &query.PageRequest{
					Key:        key,
					Limit:      limit,
					CountTotal: true, // Set to true if you need the total count
				},
			},
		)
		if err != nil {
			return err
		}

		for _, binaryAny := range resp.Pools {
			var pool clmodel.Pool
			if err := encCfg.Marshaler.Unmarshal(binaryAny.Value, &pool); err != nil {
				return err
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

			if len(incentiveRecordsResp.IncentiveRecords) > 0 {
				for i, incentiveRecord := range incentiveRecordsResp.IncentiveRecords {
					record := []string{
						fmt.Sprintf("%d", pool.Id),
						fmt.Sprintf("%d", i),
						incentiveRecord.IncentiveRecordBody.RemainingCoin.String(),
						incentiveRecord.IncentiveRecordBody.EmissionRate.String(),
						incentiveRecord.MinUptime.String(),
					}
					log.Printf("Incentives Record %d in pool %d: remainingCoin %s emissionRate %s minUptime %s\n",
						i,
						pool.Id,
						incentiveRecord.IncentiveRecordBody.RemainingCoin.String(),
						incentiveRecord.IncentiveRecordBody.EmissionRate.String(),
						incentiveRecord.MinUptime.String(),
					)
					if err := writer.Write(record); err != nil {
						panic(err) // Handle error
					}
				}
			} else {
				record := []string{
					fmt.Sprintf("%d", pool.Id),
					fmt.Sprintf("%d", -1),
					"No remaining coins",
					"No emission rate",
					"No uptime",
				}
				if err := writer.Write(record); err != nil {
					panic(err) // Handle error
				}

			}
		}

		if len(resp.Pagination.NextKey) == 0 {
			break
		}

		key = resp.Pagination.NextKey
	}

	return nil
}
