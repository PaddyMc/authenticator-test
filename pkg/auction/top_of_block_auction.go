package auction

import (
	"context"
	"log"

	grpc "google.golang.org/grpc"

	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v24/app/params"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v24/x/poolmanager/types"
	auctiontypes "github.com/skip-mev/block-sdk/x/auction/types"
)

func SubmitTopOfBlockAuction(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	signerKey *secp256k1.PrivKey,
	spoofKey *secp256k1.PrivKey,
) error {
	// set up all clients
	auctionModuleAccount := "osmo1j4yzhgjm00ch3h0p9kel7g8sp6g045qfnc9kmc"
	AuctionUSDCDenom := "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"
	osmoDenom := "uosmo"
	AtomIBCDenom := "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"

	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	tm := tmservice.NewServiceClient(conn)
	bankClient := banktypes.NewQueryClient(conn)

	//eleValAddr, _ := sdk.ValAddressFromBech32(valOpAddr)
	//eleValAddr, _ := sdk.ValAddressFromBech32("osmovaloper1tv9wnreg9z5qlxyte8526n7p3tjasndede2kj9")
	priv1 := signerKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	balancePre, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: AtomIBCDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Pre-swap balance of %s: %s\n", AtomIBCDenom, balancePre.GetBalance().Amount)

	swapTokenMsg := &poolmanagertypes.MsgSwapExactAmountOut{
		Sender: accAddress.String(),
		Routes: []poolmanagertypes.SwapAmountOutRoute{
			{
				PoolId:       1265,
				TokenInDenom: osmoDenom,
			},
		},
		TokenInMaxAmount: osmomath.NewInt(1000000000),
		TokenOut:         sdk.NewCoin(AtomIBCDenom, osmomath.NewInt(10000000)),
	}
	txBytes1, err := chaingrpc.SignAuthenticatorMsgMultiSignersBytes(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		nil,
		encCfg,
		tm,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{swapTokenMsg},
		[]uint64{},
		1,
	)

	swapTokenMsg2 := &poolmanagertypes.MsgSwapExactAmountOut{
		Sender: accAddress.String(),
		Routes: []poolmanagertypes.SwapAmountOutRoute{
			{
				PoolId:       1265,
				TokenInDenom: osmoDenom,
			},
		},
		TokenInMaxAmount: osmomath.NewInt(1000000000),
		TokenOut:         sdk.NewCoin(AtomIBCDenom, osmomath.NewInt(10000000)),
	}
	txBytes2, err := chaingrpc.SignAuthenticatorMsgMultiSignersBytes(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		nil,
		encCfg,
		tm,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{swapTokenMsg2},
		[]uint64{},
		2,
	)
	bundle := [][]byte{txBytes1, txBytes2}

	//	sequenceOffset := uint64(1)
	bidMsg := &auctiontypes.MsgAuctionBid{
		Bidder:       accAddress.String(),
		Bid:          sdk.NewCoin(AuctionUSDCDenom, sdk.NewInt(1000000)),
		Transactions: bundle,
	}

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSignersWithBlock(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		nil,
		encCfg,
		tm,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{bidMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	balancePost, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: accAddress.String(), Denom: AtomIBCDenom},
	)
	if err != nil {
		return err
	}

	balanceModuleAccount, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: auctionModuleAccount, Denom: AuctionUSDCDenom},
	)
	if err != nil {
		return err
	}

	log.Printf("Post-auction-swap balance of %s: %s\n", AtomIBCDenom, balancePost.GetBalance().Amount)
	log.Printf("Post-auction-swap balance of module account %s: %s\n", AuctionUSDCDenom, balanceModuleAccount.GetBalance().Amount)

	log.Printf("Finished auction %s", AuctionUSDCDenom)

	return nil
}
