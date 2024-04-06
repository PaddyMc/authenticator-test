package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	gov "github.com/osmosis-labs/autenticator-test/pkg/gov"
)

func StartGovernanceFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-gov-flow",
		Short: "this command does lots of gov transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting governance flow")

			log.Printf("Creating cl postion for gov module")
			err := gov.GovCreateClPosition(
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
