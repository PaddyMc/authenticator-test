package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	staking "github.com/osmosis-labs/autenticator-test/pkg/staking"
)

func StartValidatorFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-validator-flow",
		Short: "this command creates a validator",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[1]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Querying all validator delegations")

			//log.Printf("Creating a validator")
			//			err := staking.CreateValidator(
			//				conn,
			//				encCfg,
			//				seedConfig.ChainID,
			//				alice,
			//				bob,
			//				alice,
			//			)
			//			if err != nil {
			//				return err
			//			}

			err := staking.GetValidatorDelegations(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
