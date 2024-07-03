package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
			cosigners := make(map[int][]cryptotypes.PrivKey)
			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]
			LuncIBCDenom := seedConfig.DenomMap["LuncIBCDenom"]
			osmoAtomClPool := uint64(1400)
			luncOsmoBalancerPool := uint64(561)
			selectedAuthenticator := []int32{1}

			//spendLimitContractAddress := "osmo1kr95hg7c2d0u40fa379nry94z3g3tfg7r37cvm3ulr2qwackh98qh3yfsn"
			//spendLimitContractAddress := "osmo1lxmzejg7en07e0llnsc2jveymuulxjedm04j0lwgfujzrpst3gysvlf7rx"
			//spendLimitContractAddress := "osmo1nvz4e7duzmchyk4res6tdlpxpxrl9nps6vl2htlevu0a59chdarsvds5d8"
			//spendLimitContractAddress := "osmo1mf5dnx0wqv7s4v9r4ykkr7tr249pctjwxs4me5n96tch37p95ccsc3zehq"
			//spendLimitContractAddress := "osmo133290w5vgjttasqhmqcw2x9g598vg6re8lw0q0jwljlep44xvfhsp04ev8"
			//spendLimitContractAddress := "osmo19meazu70q77tzt9vzrp7d8pqf7wupcvfjqtmdtnkjqe9e0f2r5ds2er9l9"
			//spendLimitContractAddress := "osmo1f7drukwape7d320sjvp6trmdk7wju908373pwjn979nazjj20cwqc3rdm8"
			//spendLimitContractAddress := "osmo1heq8u26kn0vgf8rltxj5cqtfrwu5eggsncztjt0560mmj0ak2rrqga6rek"
			//spendLimitContractAddress := "osmo1rjq8g7mzeu99f7vlsg2c3htxnjnue5zv02j4xynl4yhe2q8r7ycscfur62"
			//spendLimitContractAddress := "osmo13j8kuxnszx9mcl5lkl92eusnx2229krlfqzpuzzt4tqvmaphzpzq6le6ge"
			spendLimitContractAddress := "osmo1rjq8g7mzeu99f7vlsg2c3htxnjnue5zv02j4xynl4yhe2q8r7ycscfur62"

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

			log.Printf("Starting swappping to Lunc, should fail")
			err = pm.SwapTokensWithLastestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				LuncIBCDenom,
				luncOsmoBalancerPool,
				// 10_000 osmo
				10000000000,
			)
			if err != nil {
				// we expected this to fail
				log.Println("error", err.Error())
			}

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
			}

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
