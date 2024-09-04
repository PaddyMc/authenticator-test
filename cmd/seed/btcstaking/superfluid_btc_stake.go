package takerfee

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/osmosis-labs/autenticator-test/pkg/config"
	gov "github.com/osmosis-labs/autenticator-test/pkg/gov"
	"github.com/osmosis-labs/osmosis/osmomath"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v26/x/poolmanager/types"
)

// https://forum.osmosis.zone/t/nbtc-revenue-share-proposal/2791
func StartTakerFeeActivationFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-add-taker-fee-rev-share-address",
		Short: "this command adds a wasm upload address transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting rev share taker fee flow")
			govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()
			localOsmosisKey10 := "osmo14gs9zqh8m49yy9kscjqu9h72exyf295afg6kgk"

			//			alloyedBTCPoolContract := "osmo1z6r6qdknhgsc0zeracktgpcxf43j6sekq07nw8sxduc9lg0qjjlqfu25e3"
			//			wBTCFactoryDenom := seedConfig.DenomMap["wBTCFactoryDenom"]
			//			wBTCAXLDenom := seedConfig.DenomMap["wBTCAXLDenom"]

			nBTCDenom := seedConfig.DenomMap["nBTCIBCDenom"]
			alloyedBTCPoolId := uint64(1868)

			setRegisteredAlloyedPoolMsg := &poolmanagertypes.MsgSetRegisteredAlloyedPool{
				Sender: govAddr,
				PoolId: alloyedBTCPoolId,
			}

			//			setTakerFeeRevShareWBTCMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
			//				Sender:      govAddr,
			//				Denom:       wBTCFactoryDenom,
			//				SkimPercent: osmomath.MustNewDecFromStr("0.1"),
			//				SkimAddress: smartAccountAddr,
			//			}
			//
			//			setTakerFeeRevShareWBTCaxlMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
			//				Sender:      govAddr,
			//				Denom:       wBTCAXLDenom,
			//				SkimPercent: osmomath.MustNewDecFromStr("0.1"),
			//				SkimAddress: smartAccountAddr,
			//			}

			setTakerFeeRevSharenBTCMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
				Sender:      govAddr,
				Denom:       nBTCDenom,
				SkimPercent: osmomath.MustNewDecFromStr("0.1"),
				SkimAddress: localOsmosisKey10,
			}

			log.Printf("Setting taker fee rev share for nBTC")
			err := gov.GovMessageProposal(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				[]sdk.Msg{setRegisteredAlloyedPoolMsg, setTakerFeeRevSharenBTCMsg},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
