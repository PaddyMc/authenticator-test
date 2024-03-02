package authenticator

import (
	"context"
	"encoding/json"
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v23/app/params"
	"github.com/osmosis-labs/osmosis/v23/x/authenticator/authenticator"
	authenticatortypes "github.com/osmosis-labs/osmosis/v23/x/authenticator/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func CreateCosignerAccount(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	parentKey *secp256k1.PrivKey,
	cosignerKey *secp256k1.PrivKey,
) error {
	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)

	priv1 := parentKey
	priv2 := cosignerKey

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
		AuthenticatorType: "SignatureVerificationAuthenticator",
		Data:              priv1.PubKey().Bytes(),
	}

	initDataPrivKey1 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "SignatureVerificationAuthenticator",
		Data:              priv2.PubKey().Bytes(),
	}

	initDataMessageFilter1 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "MessageFilterAuthenticator",
		Data:              []byte(`{"@type":"/osmosis.poolmanager.v1beta1.MsgSwapExactAmountIn"}`),
	}

	initDataMessageFilter2 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "MessageFilterAuthenticator",
		Data:              []byte(`{"@type":"/cosmos.bqnk"}`),
	}

	compositeAuthDataAllOf := []authenticator.SubAuthenticatorInitData{
		initDataPrivKey0,
		initDataPrivKey1,
	}

	compositeAuthDataAnyOf := []authenticator.SubAuthenticatorInitData{
		initDataMessageFilter1,
		initDataMessageFilter2,
	}

	dataAnyOf, err := json.Marshal(compositeAuthDataAnyOf)
	initDataAnyOf := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "AnyOfAuthenticator",
		Data:              dataAnyOf,
	}

	dataAllOf, err := json.Marshal(compositeAuthDataAllOf)
	initDataAllOf := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "PartitionedAllOfAuthenticator",
		Data:              dataAllOf,
	}

	compositeComplex := []authenticator.SubAuthenticatorInitData{
		initDataAnyOf,
		initDataAllOf,
	}

	dataCompositeComplex, err := json.Marshal(compositeComplex)
	addAllOfAuthenticatorMsg := &authenticatortypes.MsgAddAuthenticator{
		Sender: accAddress.String(),
		Type:   "AllOfAuthenticator",
		Data:   dataCompositeComplex,
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

	allAuthenticatorsPostResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	log.Println("Number of authenticators post:", len(allAuthenticatorsPostResp.AccountAuthenticators))
	if len(allAuthenticatorsPostResp.AccountAuthenticators) == len(allAuthenticatorsResp.AccountAuthenticators) {
		log.Println("Error adding spend limit authenticator")
	} else {
		log.Println("Added authenticator")
	}

	log.Println("Add authenticator completed.")

	return nil
}
