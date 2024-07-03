package seed

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	auth "github.com/osmosis-labs/autenticator-test/pkg/authenticator"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
)

func StartSmartAccountDeactivatedFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-set-smart-account-deactive-flow",
		Short: "this command does activates smart account through a gov transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting governance flow")

			log.Printf("Creating param change proposal gov module")
			err := auth.ChangeSmartAccountParams(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				alice,
			)
			if err != nil {
				return err
			}
			log.Printf("osmosisd q smartaccount params --node %s", seedConfig.GRPCConnection.Target())

			return nil
		},
	}
	return cmd
}
