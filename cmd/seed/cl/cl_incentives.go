package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	cl "github.com/osmosis-labs/autenticator-test/pkg/cl"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
)

func StartClIncentiveFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-cl-incentive-data",
		Short: "this command prints and exports incentive data for cl pools",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[1]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting cl incentive flow")

			log.Printf("Check cl pools for incentive acculmators")
			err := cl.CheckIncentiveAccumulatorForPools(
				conn,
				encCfg,
				seedConfig.ChainID,
				"accumulmator.csv",
			)
			if err != nil {
				return err
			}
			log.Printf("Incentive Accumulators exported to acculmator.csv")

			//log.Printf("Check cl pools for incentive records")
			//err = cl.CheckAllIncentivesForPools(
			//	conn,
			//	encCfg,
			//	seedConfig.ChainID,
			//	"incentives.csv",
			//)
			//if err != nil {
			//	return err
			//}
			//log.Printf("Incentives exported to incentives.csv")

			log.Printf("Finished CL incentive flow")

			return nil
		},
	}
	return cmd
}
