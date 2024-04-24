package authenticator

import (
	"context"
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v24/app/params"
	authenticatortypes "github.com/osmosis-labs/osmosis/v24/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func CreateSignatureVerificationAuthenticator(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	parentKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)

	priv1 := parentKey
	priv2 := authKey
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

	addAuthenticatorMsg := &authenticatortypes.MsgAddAuthenticator{
		Sender: accAddress.String(),
		Type:   "SignatureVerificationAuthenticator",
		Data:   priv2.PubKey().Bytes(),
	}

	log.Println("Adding authenticator for account", accAddress.String(), "first authenticator.")
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{addAuthenticatorMsg},
		[]uint64{},
	)

	allAuthenticatorsPostResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	log.Println("Number of authenticators post:", len(allAuthenticatorsResp.AccountAuthenticators))
	if len(allAuthenticatorsPostResp.AccountAuthenticators) == len(allAuthenticatorsResp.AccountAuthenticators) {
		log.Println("Error adding signature verification authenticator.")
	} else {
		log.Println("Success adding signature verification authenticator.")
	}

	log.Println("Add authenticator completed.")

	return nil
}
