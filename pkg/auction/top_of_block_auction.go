package auction

import (
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/v24/app/params"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	auctiontypes "github.com/skip-mev/block-sdk/x/auction/types"
)

func SubmitTopOfBlockAuction(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	signerKey *secp256k1.PrivKey,
) error {
	// set up all clients
	bundle := [][]byte{}
	//	height := uint64(0)
	//	sequenceOffset := uint64(1)
	AuctionUSDCDenom := "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"

	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	//bankClient := banktypes.NewQueryClient(conn)

	//eleValAddr, _ := sdk.ValAddressFromBech32(valOpAddr)
	//eleValAddr, _ := sdk.ValAddressFromBech32("osmovaloper1tv9wnreg9z5qlxyte8526n7p3tjasndede2kj9")
	priv1 := signerKey
	accAddress := sdk.AccAddress(priv1.PubKey().Address())

	bidMsg := &auctiontypes.MsgAuctionBid{
		Bidder:       accAddress.String(),
		Bid:          sdk.NewCoin(AuctionUSDCDenom, sdk.NewInt(1000000)),
		Transactions: bundle,
	}

	err := chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{signerKey},
		[]cryptotypes.PrivKey{signerKey},
		nil,
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{bidMsg},
		[]uint64{},
	)
	if err != nil {
		return err
	}

	//	limit := uint64(1000)
	//	key := []byte{}
	//	validatorDelegations, err := stakingClient.ValidatorDelegations(
	//		context.Background(),
	//		&stakingtypes.QueryValidatorDelegationsRequest{
	//			ValidatorAddr: eleValAddr.String(),
	//			Pagination: &query.PageRequest{
	//				Key:   key,
	//				Limit: limit,
	//			},
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//fmt.Println(validatorDelegations)

	log.Printf("Finished auction %s", AuctionUSDCDenom)

	return nil
}
