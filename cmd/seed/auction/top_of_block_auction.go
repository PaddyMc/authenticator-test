package auction

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	auction "github.com/osmosis-labs/autenticator-test/pkg/auction"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
)

func StartAuctionFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-auction-flow",
		Short: "this command submits multiple auction txns",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[1]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Start top of block auction flow")

			err := auction.SubmitTopOfBlockAuction(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
