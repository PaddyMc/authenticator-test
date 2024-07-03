package seed

import (
	"log"

	"github.com/spf13/cobra"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/osmosis-labs/autenticator-test/pkg/config"
	gov "github.com/osmosis-labs/autenticator-test/pkg/gov"
)

func StartWasmUploadAddressFlow(seedConfig config.SeedConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-add-wasm-upload-address-flow",
		Short: "this command adds a wasm upload address transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := seedConfig.GRPCConnection
			encCfg := seedConfig.EncodingConfig

			alice := seedConfig.Keys[0]
			bob := seedConfig.Keys[2]

			cosigners := make(map[int][]cryptotypes.PrivKey)
			cosigners[1] = []cryptotypes.PrivKey{alice, bob}

			log.Printf("Starting wasm upload flow")
			govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

			updateParamsMsg := &wasmtypes.MsgAddCodeUploadParamsAddresses{
				Authority: govAddr,
				Addresses: []string{"osmo1j0u8d5sl6fnp50tr79hw0agrql9qn0xuzyyxr5lned0xsp05znmqallwcw"},
			}

			log.Printf("Creating wasm upload address proposal gov module")
			err := gov.GovMessageProposal(
				conn,
				encCfg,
				seedConfig.ChainID,
				alice,
				bob,
				alice,
				[]sdk.Msg{updateParamsMsg},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
