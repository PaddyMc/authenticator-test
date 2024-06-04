package authenticator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v24/app/params"
	"github.com/osmosis-labs/osmosis/v24/x/smart-account/authenticator"
	authenticatortypes "github.com/osmosis-labs/osmosis/v24/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func CreateOneClickTradingAccount(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	parentKey *secp256k1.PrivKey,
	tradingKey *secp256k1.PrivKey,
	spendLimitContractAddress string,
) error {
	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)
	//wasmClient := wasmtypes.NewQueryClient(conn)

	priv1 := parentKey
	priv2 := tradingKey

	accAddress := sdk.AccAddress(priv1.PubKey().Address())
	//accAddress2 := sdk.AccAddress(priv2.PubKey().Address())

	allAuthenticatorsResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	log.Println("Querying authenticators for account", accAddress.String())
	log.Println("Number of authenticators:", len(allAuthenticatorsResp.AccountAuthenticators))

	// initialise spend limit authenticator
	initDataPrivKey0 := authenticator.SubAuthenticatorInitData{
		Type:   "SignatureVerification",
		Config: priv2.PubKey().Bytes(),
	}

	// Time limit for spend limits
	now := time.Now()
	future := now.Add(time.Hour * 3)

	jsonString := fmt.Sprintf(
		`{"time_limit": {"end": "%d"}, "reset_period": "day", "limit": "10000000000"}`, future.UnixNano())
	encodedParams := base64.StdEncoding.EncodeToString([]byte(jsonString))
	initDataSpendLimit := authenticator.SubAuthenticatorInitData{
		Type: "CosmwasmAuthenticatorV1",
		Config: []byte(
			`{"contract": "` + spendLimitContractAddress + `", "params": "` + encodedParams + `"}`),
	}

	initDataMessageFilter := authenticator.SubAuthenticatorInitData{
		Type:   "MessageFilter",
		Config: []byte(`{"@type":"/osmosis.poolmanager.v1beta1.MsgSwapExactAmountIn"}`),
	}
	compositeAuthData := []authenticator.SubAuthenticatorInitData{
		initDataPrivKey0,
		initDataSpendLimit,
		initDataMessageFilter,
	}

	dataAllOf, err := json.Marshal(compositeAuthData)
	addAllOfAuthenticatorMsg := &authenticatortypes.MsgAddAuthenticator{
		Sender: accAddress.String(),
		Type:   "AllOf",
		Data:   dataAllOf,
	}

	log.Println("Adding authenticator for account", accAddress.String(), "first authenticator")
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{addAllOfAuthenticatorMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	allAuthenticatorsPostResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	//auth := wasmClient.QuerySmart()

	log.Println("Number of authenticators post:", len(allAuthenticatorsPostResp.AccountAuthenticators))
	if len(allAuthenticatorsPostResp.AccountAuthenticators) == len(allAuthenticatorsResp.AccountAuthenticators) {
		log.Println("Error adding spend limit authenticator")
	} else {
		log.Println("Added authenticator")
	}

	log.Println("Add authenticator completed.")

	return nil
}
