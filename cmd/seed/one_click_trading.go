package seed

import (
	"log"

	"github.com/spf13/cobra"

	as "github.com/osmosis-labs/autenticator-test/pkg/authenticator"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	pm "github.com/osmosis-labs/autenticator-test/pkg/poolmanager"
)

func SeedCreateOneClickTradingAccount(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-one-click-trading-flow",
		Short: "this command goes through a series of tasks to test the one click trading flow",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]
			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			LuncIBCDenom := seedConfig.DenomMap["LuncIBCDenom"]
			osmoAtomClPool := uint64(1400)
			luncOsmoBalancerPool := uint64(561)
			selectedAuthenticator := []int32{1}

			//spendLimitContractAddress := "osmo1kr95hg7c2d0u40fa379nry94z3g3tfg7r37cvm3ulr2qwackh98qh3yfsn"
			spendLimitContractAddress := "osmo1lxmzejg7en07e0llnsc2jveymuulxjedm04j0lwgfujzrpst3gysvlf7rx"

			log.Printf("Starting spend limit authenticator flow")
			log.Printf("Adding spend limit authenticator")
			err := as.CreateOneClickTradingAccount(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				spendLimitContractAddress,
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
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				osmoAtomClPool,
				1000000,
			)
			if err != nil {
				log.Println("error", err.Error())
				//return err
			}

			log.Printf("Starting swappping to Lunc, should fail")
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				selectedAuthenticator,
				OsmoDenom,
				LuncIBCDenom,
				luncOsmoBalancerPool,
				100000000000,
			)
			if err != nil {
				// we expected this to fail
				log.Println("error", err.Error())
			}

			log.Printf("Removing spend limit authenticator")
			//			err = as.RemoveLatestAuthenticator(
			//				conn,
			//				encCfg,
			//				seedConfig.ChainID,
			//				alice,
			//				alice,
			//			)
			//			if err != nil {
			//				return err
			//			}

			return nil
		},
	}
	return cmd
}
