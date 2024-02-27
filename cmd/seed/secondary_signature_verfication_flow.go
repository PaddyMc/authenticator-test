package seed

import (
	"log"

	as "github.com/osmosis-labs/autenticator-test/pkg/authenticator"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	pm "github.com/osmosis-labs/autenticator-test/pkg/poolmanager"

	"github.com/spf13/cobra"
)

func SeedSwapCmd(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-swap-with-signature-authenticator-flow",
		Short: "this command creates SignatureVerificationAuthenticator and swaps in a pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig
			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			osmoAtomClPool := uint64(1400)
			selectedAuthenticator := []int32{1}

			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]

			log.Println("Starting swap flow")
			err := as.CreateSignatureVerificationAuthenticator(
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
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100000000,
			)
			if err != nil {
				log.Println("Transaction Failed...", err.Error())
				// return err
			}

			log.Printf("Removing spend limit authenticator")
			err = as.RemoveLatestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				alice,
			)

			return nil
		},
	}
	return cmd
}
