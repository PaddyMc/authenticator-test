package cl

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

	//govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	//query "github.com/cosmos/cosmos-sdk/types/query"

	//transfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v25/app/params"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v25/x/poolmanager/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func SwapTokens(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []uint64,
	fromToken string,
	toToken string,
	poolId uint64,
	tokenAmount int64,
) error {
	log.Println("Starting token swap...")

	// set up all grpc clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	bankClient := banktypes.NewQueryClient(conn)

	// decode the test private keys
	priv1 := fromKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	balancePre, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: toToken},
	)
	if err != nil {
		return err
	}

	log.Printf("Pre-swap balance of %s: %s\n", toToken, balancePre.GetBalance().Amount)

	swapTokenMsg := &poolmanagertypes.MsgSwapExactAmountOut{
		Sender: accAddress.String(),
		Routes: []poolmanagertypes.SwapAmountOutRoute{
			{
				PoolId:       poolId,
				TokenInDenom: fromToken,
			},
		},
		TokenInMaxAmount: osmomath.NewInt(tokenAmount),
		TokenOut:         sdk.NewCoin(toToken, osmomath.NewInt(10000000)),
	}

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{fromKey},
		[]cryptotypes.PrivKey{fromKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{swapTokenMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	balancePost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: toToken},
	)
	if err != nil {
		return err
	}

	log.Printf("Post-swap balance of %s: %s\n", toToken, balancePost.GetBalance().Amount)

	return nil
}
