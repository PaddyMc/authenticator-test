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
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v24/app/params"

	authenticatortypes "github.com/osmosis-labs/osmosis/v24/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func SendTokens(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	toKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
	fromToken string,
	toToken string,
	tokenAmount int64,
) error {
	log.Println("Starting token send...")

	// use values from the context
	OsmoDenom := fromToken

	// set up all grpc clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	bankClient := banktypes.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)

	// decode the test private keys
	priv1 := fromKey
	priv2 := toKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())
	accAddressTo := sdk.AccAddress(priv2.PubKey().Address())
	log.Println("Info", "osmosisd query bank balances", accAddress.String(), "--node http://164.92.247.225:26657 --output json | jq")

	_, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	balanceAtom, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: OsmoDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Initial balance of %s: %s\n", OsmoDenom, balanceAtom.GetBalance().Amount)

	sendTokenMsg := &banktypes.MsgSend{
		FromAddress: accAddress.String(),
		ToAddress:   "osmo1adz3qee0s03zw257mfgtl7yzmjqlkj2mxcyhfn",
		Amount:      sdk.NewCoins(sdk.NewCoin(OsmoDenom, osmomath.NewInt(tokenAmount))),
	}

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{fromKey},
		[]cryptotypes.PrivKey{fromKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{sendTokenMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	balanceAtomPost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddressTo.String(), Denom: OsmoDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Post-send balance of %s: %s\n", OsmoDenom, balanceAtomPost.GetBalance().Amount)

	// Check if the balance has increased
	if balanceAtomPost.GetBalance().Amount.GTE(balanceAtom.GetBalance().Amount) {
		log.Println("Send failed token balance did not reduct")
	} else {
		log.Println("Successful send token.")
	}

	sendBackTokenMsg := &banktypes.MsgSend{
		FromAddress: accAddressTo.String(),
		ToAddress:   accAddress.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(OsmoDenom, osmomath.NewInt(1000))),
	}

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{toKey},
		[]cryptotypes.PrivKey{toKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{sendBackTokenMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	log.Println("Token send completed.")

	return nil
}
