package seed

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	wasm "github.com/osmosis-labs/autenticator-test/pkg/wasm"
)

func StartUploadLevanaContractsFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-upload-levana-flow",
		Short: "this command uploads levana contracts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting wasm upload flow")

			// Define the path to the JSON file
			initMsgPath := ""

			err := wasm.UploadAndInstantiateContract(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				alice,
				"./cmd/contracts/spendlimit/spend_limit.wasm",
				initMsgPath,
			)
			if err != nil {
				return err
			}
			log.Printf("uploaded spendlimit successfully wasm upload flow")

			return nil
		},
	}
	return cmd
}
