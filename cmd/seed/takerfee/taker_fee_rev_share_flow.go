package takerfee

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	as "github.com/osmosis-labs/autenticator-test/pkg/authenticator"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	pm "github.com/osmosis-labs/autenticator-test/pkg/poolmanager"

	"github.com/spf13/cobra"
)

func SeedSwapToEarnTakerFeeCmd(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-taker-fee-rev-share-nbtc",
		Short: "this command creates swaps for nbtc and checks the accumulator for rev share",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig
			UsdcIBCDenom := seedConfig.DenomMap["UsdcIBCDenom"]
			nBTCIBCDenom := seedConfig.DenomMap["nBTCIBCDenom"]
			//wBTCFactoryDenom := seedConfig.DenomMap["wBTCFactoryDenom"]
			//wBTCAXLDenom := seedConfig.DenomMap["wBTCAXLDenom"]
			usdcNBTCClPool := uint64(1253)

			selectedAuthenticator := []int32{}

			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]
			cosigners := make(map[int][]cryptotypes.PrivKey)

			log.Printf("Starting swap flow to test take fees from %s to %s", UsdcIBCDenom, nBTCIBCDenom)
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
				cosigners,
				selectedAuthenticator,
				UsdcIBCDenom,
				nBTCIBCDenom,
				usdcNBTCClPool,
				200_000_000,
			)
			if err != nil {
				log.Println("Transaction Failed...", err.Error())
				return err
			}

			log.Printf("Removing spend limit authenticator")
			err = as.RemoveLatestAuthenticator(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				alice,
			)
			if err != nil {
				log.Println("Transaction Failed...", err.Error())
				return err
			}

			return nil
		},
	}
	return cmd
}
