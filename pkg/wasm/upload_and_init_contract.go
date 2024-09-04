package seed

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	"github.com/osmosis-labs/osmosis/v26/app/params"
)

func UploadAndInstantiateContract(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	uploadKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	contractLocation string,
	instantiateMsgPath string,
) error {
	priv1 := uploadKey
	//	priv2 := signerKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	wasmClient := wasmtypes.NewQueryClient(conn)

	// "./cw_authenticators/cosigner_authenticator.wasm"
	wasm, err := os.ReadFile(contractLocation)
	if err != nil {
		panic(err)
	}

	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			panic(err)
		}
	} else if !ioutils.IsGzip(wasm) {
		panic(fmt.Errorf("invalid input file. Use wasm binary or gzip"))
	}

	msg := &wasmtypes.MsgStoreCode{
		Sender:                accAddress.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: nil,
	}
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{msg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	codes, err := wasmClient.Codes(
		context.Background(),
		&wasmtypes.QueryCodesRequest{
			Pagination: &query.PageRequest{
				Reverse: true,
			},
		},
	)
	codeID := codes.CodeInfos[0].CodeID

	// init contract
	//	instantiateMsg := CosignerInstantiateMsg{PubKeys: [][]byte{priv2.PubKey().Bytes()}}
	//	instantiateMsgBz, err := json.Marshal(instantiateMsg)

	if instantiateMsgPath != "" {
		instantiateMsgBz, err := os.ReadFile(instantiateMsgPath)
		if err != nil {
			panic(err)
		}
		initMsg := &wasmtypes.MsgInstantiateContract{
			Sender: accAddress.String(),
			CodeID: codeID,
			Label:  "contract",
			Msg:    instantiateMsgBz,
		}
		err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
			[]cryptotypes.PrivKey{signerKey},
			[]cryptotypes.PrivKey{signerKey},
			make(map[int][]cryptotypes.PrivKey),
			encCfg,
			ac,
			txClient,
			chainID,
			[]sdk.Msg{initMsg},
			[]uint64{},
		)

		contracts, err := wasmClient.ContractsByCode(
			context.Background(),
			&wasmtypes.QueryContractsByCodeRequest{CodeId: codeID},
		)
		contractAddress := contracts.Contracts[len(contracts.Contracts)-1]
		log.Println("Contract address: ", contractAddress)

	}

	return nil
}
