package orderbook

import (
	"log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	orderbook "github.com/osmosis-labs/autenticator-test/pkg/orderbook"

	"github.com/spf13/cobra"
)

func SeedBatchClaimOrdersFromAllOrderbooks(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-batch-claim-orderbook",
		Short: "this command batch claims orders from defined orderbooks",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			// NOTE: These are for the smart account module and can be ignored
			selectedAuthenticator := []int32{}
			alice := seedConfig.Keys[2]
			bob := seedConfig.Keys[3]
			cosigners := make(map[int][]cryptotypes.PrivKey)

			// sumtree orderbook contract code id: 885
			// https://celatone.osmosis.zone/osmosis-1/codes/885
			sumtreeOrderbookCodeId := uint64(885)
			log.Printf("Starting claim outstanging orders for codeID: %d", sumtreeOrderbookCodeId)
			err := orderbook.ClaimLimitOrders(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				cosigners,
				selectedAuthenticator,
				sumtreeOrderbookCodeId,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
