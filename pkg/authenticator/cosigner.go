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
	"github.com/osmosis-labs/osmosis/v24/app/params"
	"github.com/osmosis-labs/osmosis/v24/x/smart-account/authenticator"
	authenticatortypes "github.com/osmosis-labs/osmosis/v24/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

// CreateCosignerAccount creates a composit authenticator that can be used for cosigning
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

	allAuthenticatorsResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	log.Println("Querying authenticators for account", accAddress.String())
	log.Println("Number of authenticators:", len(allAuthenticatorsResp.AccountAuthenticators))

	// Create the first AllOfAuthenticator
	initDataPrivKey0 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "SignatureVerificationAuthenticator",
		Data:              priv1.PubKey().Bytes(),
	}

	initDataPrivKey1 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "SignatureVerificationAuthenticator",
		Data:              priv2.PubKey().Bytes(),
	}

	// Create the message filter AnyOf authenticator
	initDataMessageFilter1 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "MessageFilterAuthenticator",
		Data:              []byte(`{"@type":"/osmosis.poolmanager.v1beta1.MsgSwapExactAmountIn"}`),
	}

	initDataMessageFilter2 := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "MessageFilterAuthenticator",
		Data:              []byte(`{"@type":"/cosmos.bqnk"}`),
	}

	// Compose the data for cosigner authenticator
	compositeAuthDataAllOf := []authenticator.SubAuthenticatorInitData{
		initDataPrivKey0,
		initDataPrivKey1,
	}

	// Compose the data for anyof messagefilter authenticator
	compositeAuthDataAnyOf := []authenticator.SubAuthenticatorInitData{
		initDataMessageFilter1,
		initDataMessageFilter2,
	}

	dataAnyOf, err := json.Marshal(compositeAuthDataAnyOf)
	initDataAnyOf := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "AnyOfAuthenticator",
		Data:              dataAnyOf,
	}

	// This is the cosigner authenticator
	dataAllOf, err := json.Marshal(compositeAuthDataAllOf)
	initDataAllOf := authenticator.SubAuthenticatorInitData{
		AuthenticatorType: "PartitionedAllOfAuthenticator",
		Data:              dataAllOf,
	}

	// Compose the AllOf and the AnyOf
	compositeComplex := []authenticator.SubAuthenticatorInitData{
		initDataAnyOf,
		initDataAllOf,
	}

	// This is the final state of the authenticator that will be sent to the chain
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
