package seed

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	gov "github.com/osmosis-labs/autenticator-test/pkg/gov"
)

func StartGovernanceFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-param-change-flow",
		Short: "this command does lots of gov transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting governance flow")

			changes := proposal.ParamChange{
				Subspace: "smartaccount",
				Key:      "MaximumUnauthenticatedGas",
				Value:    `"105000"`,
			}

			paramChange := proposal.NewParameterChangeProposal(
				"Update the max gas for authenticators",
				"Updating the gas to 100000",
				[]proposal.ParamChange{changes},
			)

			//		fmt.Println(auctiontypes.KeyParams)
			//			govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()
			//			AuctionParams := auctiontypes.Params{
			//				MaxBundleSize:          5,
			//				ReserveFee:             sdk.NewCoin("uosmo", osmomath.NewInt(1000000)),
			//				MinBidIncrement:        sdk.NewCoin("uosmo", osmomath.NewInt(1000000)),
			//				EscrowAccountAddress:   auctiontypes.DefaultEscrowAccountAddress,
			//				FrontRunningProtection: true,
			//				ProposerFee:            osmomath.MustNewDecFromStr("0.05"),
			//			}
			//
			//			updateParamsMsg := &auctiontypes.MsgUpdateParams{
			//				Authority: govAddr,
			//				Params:    AuctionParams,
			//			}

			log.Printf("Creating param change proposal gov module")
			err := gov.ParameterChangeProposal(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				//[]sdk.Msg{paramChange},
				paramChange,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
