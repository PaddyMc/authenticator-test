package poolmanager

import (
	"context"
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v26/app/params"

	//govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v26/x/poolmanager/types"
	authenticatortypes "github.com/osmosis-labs/osmosis/v26/x/smart-account/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func SwapTokensWithLastestAuthenticator(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []int32,
	fromToken string,
	toToken string,
	poolId uint64,
	tokenAmount int64,
) error {
	log.Println("Starting token swap...")

	// use values from the context
	OsmoDenom := fromToken
	AtomIBCDenom := toToken
	osmoAtomClPool := poolId

	// set up all grpc clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	bankClient := banktypes.NewQueryClient(conn)
	authenticatorClient := authenticatortypes.NewQueryClient(conn)

	// decode the test private keys
	priv1 := fromKey
	//priv2 := toKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	allAuthenticatorsResp, err := authenticatorClient.GetAuthenticators(
		context.Background(),
		&authenticatortypes.GetAuthenticatorsRequest{Account: accAddress.String()},
	)
	if err != nil {
		return err
	}

	balanceAtom, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: AtomIBCDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Initial balance of %s: %s\n", AtomIBCDenom, balanceAtom.GetBalance().Amount)

	swapTokenMsg := &poolmanagertypes.MsgSwapExactAmountIn{
		Sender: accAddress.String(),
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{
				PoolId:        osmoAtomClPool,
				TokenOutDenom: AtomIBCDenom,
			},
		},
		TokenOutMinAmount: osmomath.NewInt(1),
		TokenIn:           sdk.NewCoin(OsmoDenom, osmomath.NewInt(tokenAmount)),
	}

	authenticatorId := allAuthenticatorsResp.AccountAuthenticators[len(allAuthenticatorsResp.AccountAuthenticators)-1].Id
	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{fromKey},
		[]cryptotypes.PrivKey{signerKey},
		cosignersKeys,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{swapTokenMsg},
		[]uint64{authenticatorId},
	)
	if err != nil {
		return err
	}

	balanceAtomPost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: AtomIBCDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Post-swap balance of %s: %s\n", AtomIBCDenom, balanceAtomPost.GetBalance().Amount)

	// Check if the balance has increased
	if balanceAtomPost.GetBalance().Amount.GT(balanceAtom.GetBalance().Amount) {
		log.Println("Balance of toToken has increased after the swap.")
	} else {
		log.Println("Balance of toToken did not increase as expected.")
	}

	log.Println("Token swap completed.")

	return nil
}
