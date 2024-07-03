package seed

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	gov "github.com/osmosis-labs/autenticator-test/pkg/gov"
)

func StartUpdateSmartAccountControllerFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-update-smart-account-controller-flow",
		Short: "this command does updates the smart account controller transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig
			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting governance flow")

			changes := proposal.ParamChange{
				Subspace: "smartaccount",
				Key:      "CircuitBreakerControllers",
				Value:    "[\"osmo1j0u8d5sl6fnp50tr79hw0agrql9qn0xuzyyxr5lned0xsp05znmqallwcw\"]",
			}

			paramChange := proposal.NewParameterChangeProposal(
				"Update CircuitBreakerControllers to DAODAO",
				"update CircuitBreakerControllers to testnet address",
				[]proposal.ParamChange{changes},
			)

			log.Printf("Creating param change proposal gov module")
			err := gov.ParameterChangeProposal(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				paramChange,
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
