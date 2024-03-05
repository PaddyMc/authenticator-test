package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	bank "github.com/osmosis-labs/autenticator-test/pkg/bank"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
)

func StartBankSendFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-bank-send-flow",
		Short: "this command does bank sends to multiple accounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[1]
			charlie := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}
			selectedAuthenticator := []uint64{1}

			OsmoDenom := seedConfig.DenomMap["OsmoDenom"]
			AtomIBCDenom := seedConfig.DenomMap["AtomIBCDenom"]

			log.Printf("Starting bank send flow")

			log.Printf("Bank send from alice to bob")
			err := bank.SendTokens(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				10000000,
			)
			if err != nil {
				return err
			}

			log.Printf("Bank send from alice to charlie")
			err = bank.SendTokens(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				charlie,
				alice,
				cosigners,
				selectedAuthenticator,
				OsmoDenom,
				AtomIBCDenom,
				10000000,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
