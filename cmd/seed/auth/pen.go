package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	as "github.com/osmosis-labs/autenticator-test/pkg/authenticator"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	pm "github.com/osmosis-labs/autenticator-test/pkg/poolmanager"
)

func StartPenetrationTest(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-pen-test-flow",
		Short: "this command goes pen tests the authenticator module",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]
			cosigners := make(map[int][]cryptotypes.PrivKey)
			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			osmoAtomClPool := uint64(1400)
			selectedAuthenticator := []int32{1}

			log.Printf("Starting spend limit authenticator flow")
			log.Printf("Adding spend limit authenticator")
			err := as.CreateMassiveAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
			)
			if err != nil {
				return err
			}

			log.Printf("Starting swap flow")
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				100,
			)
			if err != nil {
				log.Println("error", err.Error())
				//return err
			}

			log.Printf("Removing spend limit authenticator")
			//err = as.RemoveLatestAuthenticator(
			//	conn,
			//	encCfg,
			//	seedConfig.ChainID,
			//	alice,
			//	alice,
			//)
			//if err != nil {
			//	return err
			//}

			return nil
		},
	}
	return cmd
}
