package authenticator

import (
	"encoding/json"
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v26/app/params"
	"github.com/osmosis-labs/osmosis/v26/x/smart-account/authenticator"
	authenticatortypes "github.com/osmosis-labs/osmosis/v26/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func CreateMassiveAuthenticator(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	parentKey *secp256k1.PrivKey,
	tradingKey *secp256k1.PrivKey,
) error {
	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	//	authenticatorClient := authenticatortypes.NewQueryClient(conn)
	//wasmClient := wasmtypes.NewQueryClient(conn)

	priv1 := parentKey
	//priv2 := tradingKey

	accAddress := sdk.AccAddress(priv1.PubKey().Address())
	//accAddress2 := sdk.AccAddress(priv2.PubKey().Address())

	//	allAuthenticatorsResp, err := authenticatorClient.GetAuthenticators(
	//		context.Background(),
	//		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	//	)
	//	if err != nil {
	//		return err
	//	}

	//	log.Println("Querying authenticators for account", accAddress.String())
	//	log.Println("Number of authenticators:", len(allAuthenticatorsResp.AccountAuthenticators))

	compositeFoulAll := []authenticator.SubAuthenticatorInitData{}

	for i := 0; i < 3000; i++ {
		initDataMessageFilter := authenticator.SubAuthenticatorInitData{
			Type:   "MessageFilter",
			Config: []byte(`{"@type":"/osmosis.poolmanager.v1beta1.MsgSwapExactAmountIn"}`),
		}
		compositeFoulAll = append(compositeFoulAll, initDataMessageFilter)
	}
	initDataMessageFilterRoute := authenticator.SubAuthenticatorInitData{
		Type:   "MessageFilter",
		Config: []byte(`{"@type":"/osmosis.poolmanager.v1beta1.MsgSplitRouteSwapExactAmountIn"}`),
	}
	compositeFoulAll = append(compositeFoulAll, initDataMessageFilterRoute)
	//	dataAnyOf, err := json.Marshal(compositeAuthDataAll)
	//	initDataAnyOfAuthenticator := authenticator.SubAuthenticatorInitData{
	//		Type:   "AnyOf",
	//		Config: dataAnyOf,
	//	}
	//	compositeAuthData := []authenticator.SubAuthenticatorInitData{
	//		initDataPrivKey0,
	//		initDataSpendLimit,
	//		initDataAnyOfAuthenticator,
	//	}

	dataAllOf, err := json.Marshal(compositeFoulAll)
	addAllOfAuthenticatorMsg := &authenticatortypes.MsgAddAuthenticator{
		Sender:            accAddress.String(),
		AuthenticatorType: "AllOf",
		Data:              dataAllOf,
	}

	log.Println("Adding authenticator for account", accAddress.String(), "first authenticator")
	for i := 0; i < 50; i++ {
		err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
			[]cryptotypes.PrivKey{
				priv1,
			},
			[]cryptotypes.PrivKey{
				priv1,
			},
			make(map[int][]cryptotypes.PrivKey),
			encCfg,
			ac,
			txClient,
			chainID,
			[]sdk.Msg{
				addAllOfAuthenticatorMsg,
			},
			[]uint64{},
		)
		if err != nil {
			return err
		}
	}

	//	allAuthenticatorsPostResp, err := authenticatorClient.GetAuthenticators(
	//		context.Background(),
	//		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//
	//	log.Println("Number of authenticators post:", len(allAuthenticatorsPostResp.AccountAuthenticators))
	//	if len(allAuthenticatorsPostResp.AccountAuthenticators) == len(allAuthenticatorsResp.AccountAuthenticators) {
	//		log.Println("Error adding spend limit authenticator")
	//	} else {
	//		log.Println("Added authenticator")
	//	}

	log.Println("Add authenticator completed.")

	return nil
}
