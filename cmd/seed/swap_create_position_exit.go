package seed

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cl "github.com/osmosis-labs/autenticator-test/pkg/cl"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
)

func StartClSwapAndTransferPositionFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-cl-swap-flow",
		Short: "this command swaps in a cl pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[1]
			//charlie := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}
			selectedAuthenticator := []uint64{1}

			log.Printf("Starting cl swap flow")

			log.Printf("Start swapping and checking incentives are distributed")
			err := cl.SwapAndCreatePositionInCLPool(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				1325,
				100000000,
			)
			if err != nil {
				return err
			}

			log.Printf("Start create position and transfer position")
			accAddress := sdk.AccAddress(bob.PubKey().Address())
			err = cl.CreatePositionAndTransfer(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				1325,
				100000000,
				accAddress.String(),
			)
			if err != nil {
				return err
			}

			log.Printf("Finished CL create remove position flow")

			return nil
		},
	}
	return cmd
}
