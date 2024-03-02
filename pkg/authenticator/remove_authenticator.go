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
	"github.com/osmosis-labs/osmosis/v23/app/params"
	authenticatortypes "github.com/osmosis-labs/osmosis/v23/x/authenticator/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func RemoveLatestAuthenticator(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	senderKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)
	priv1 := senderKey
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

	removeMsg := &authenticatortypes.MsgRemoveAuthenticator{
		Sender: accAddress.String(),
		Id:     allAuthenticatorsResp.AccountAuthenticators[len(allAuthenticatorsResp.AccountAuthenticators)-1].Id,
	}

	log.Println("Removing authenticator for account", accAddress.String())
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{removeMsg},
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
		log.Println("Error removing spend limit authenticator")
	} else {
		log.Println("Removed authenticator")
	}

	log.Println("Remove authenticator completed.")

	return nil
}
