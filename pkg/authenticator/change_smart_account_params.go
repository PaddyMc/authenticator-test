package authenticator

import (
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v25/app/params"
	authenticatortypes "github.com/osmosis-labs/osmosis/v25/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func ChangeSmartAccountParams(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	senderKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	priv1 := senderKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())
	//accAddress2 := sdk.AccAddress(priv2.PubKey().Address())

	log.Println("Deactivating smartaccounts", accAddress.String())

	setActiveStateMsg := &authenticatortypes.MsgSetActiveState{
		Sender: accAddress.String(),
		Active: false,
	}

	log.Println("Removing authenticator for account", accAddress.String())
	err := chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{setActiveStateMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	log.Println("Remove authenticator completed.")

	return nil
}
