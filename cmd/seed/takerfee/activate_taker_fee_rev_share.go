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
	poolmanagertypes "github.com/osmosis-labs/osmosis/v25/x/poolmanager/types"
	smartaccounttypes "github.com/osmosis-labs/osmosis/v25/x/smart-account/types"
)

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
			smartAccountAddr := authtypes.NewModuleAddress(smartaccounttypes.ModuleName).String()

			//alloyedBTCPoolContract := "osmo1z6r6qdknhgsc0zeracktgpcxf43j6sekq07nw8sxduc9lg0qjjlqfu25e3"
			nBTCDenom := seedConfig.DenomMap["nBTCIBCDenom"]
			//			wBTCFactoryDenom := seedConfig.DenomMap["wBTCFactoryDenom"]
			//			wBTCAXLDenom := seedConfig.DenomMap["wBTCAXLDenom"]
			//			alloyedBTCPoolId := uint64(1868)

			//		setRegisteredAlloyedPoolMsg := &poolmanagertypes.MsgSetRegisteredAlloyedPool{
			//			Sender: govAddr,
			//			PoolId: alloyedBTCPoolId,
			//		}

			//		setTakerFeeRevShareWBTCMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
			//			Sender:      govAddr,
			//			Denom:       wBTCFactoryDenom,
			//			SkimPercent: osmomath.MustNewDecFromStr("0.1"),
			//			SkimAddress: smartAccountAddr,
			//		}

			//		setTakerFeeRevShareWBTCaxlMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
			//			Sender:      govAddr,
			//			Denom:       wBTCAXLDenom,
			//			SkimPercent: osmomath.MustNewDecFromStr("0.1"),
			//			SkimAddress: smartAccountAddr,
			//		}

			setTakerFeeRevSharenBTCMsg := &poolmanagertypes.MsgSetTakerFeeShareAgreementForDenom{
				Sender:      govAddr,
				Denom:       nBTCDenom,
				SkimPercent: osmomath.MustNewDecFromStr("0.1"),
				SkimAddress: smartAccountAddr,
			}

			log.Printf("Setting taker fee rev share for nBTC")
			err := gov.GovMessageProposal(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				[]sdk.Msg{setTakerFeeRevSharenBTCMsg /*setRegisteredAlloyedPoolMsg, setTakerFeeRevShareWBTCMsg, setTakerFeeRevShareWBTCaxlMsg*/},
			)
			if err != nil {
				return err
			}

			log.Printf("Please wait for the gov prop to pass then query the node with:")
			log.Printf("osmosisd query poolmanager all-taker-fee-share-agreements --node %s", seedConfig.GRPCConnection.Target())
			log.Printf("osmosisd query poolmanager taker-fee-share-agreement-from-denom %s --node %s", nBTCDenom, seedConfig.GRPCConnection.Target())
			log.Printf("osmosisd query poolmanager all-taker-fee-share-accumulators --node %s", seedConfig.GRPCConnection.Target())
			log.Printf("osmosisd query poolmanager all-registered-alloyed-pools --node %s", seedConfig.GRPCConnection.Target())

			return nil
		},
	}
	return cmd
}
