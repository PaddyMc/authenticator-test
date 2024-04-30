package grpc

import (
	"context"
	"log"

	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	key "github.com/osmosis-labs/autenticator-test/pkg/key"
	"github.com/osmosis-labs/osmosis/v24/app/params"
)

func SignAuthenticatorMsgMultiSignersBytes(
	senderPrivKeys []cryptotypes.PrivKey,
	signerPrivKeys []cryptotypes.PrivKey,
	cosignerPrivKeys map[int][]cryptotypes.PrivKey,
	encCfg params.EncodingConfig,
	tm tmservice.ServiceClient,
	ac authtypes.QueryClient,
	txClient txtypes.ServiceClient,
	chainID string,
	msgs []sdk.Msg,
	selectedAuthenticators []uint64,
	sequenceOffset uint64,
) ([]byte, error) {
	log.Println("Creating signed txn to include in bundle")

	var accNums []uint64
	var accSeqs []uint64

	for _, privKey := range senderPrivKeys {
		// Generate the account address from the private key
		addr := sdk.AccAddress(privKey.PubKey().Address()).String()

		// Get the account information
		res, err := ac.Account(
			context.Background(),
			&authtypes.QueryAccountRequest{Address: addr},
		)
		if err != nil {
			return nil, err
		}

		var acc authtypes.AccountI
		if err := encCfg.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
			return nil, err
		}

		log.Println("Signer account: " + acc.GetAddress().String())
		accNums = append(accNums, acc.GetAccountNumber())
		// XXX: here we return + 1 to offset the seq
		accSeqs = append(accSeqs, acc.GetSequence()+sequenceOffset)
	}

	block, err := tm.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}

	// Sign the message
	txBytes, _ := key.SignAuthenticatorMsgWithHeight(
		encCfg.TxConfig,
		msgs,
		sdk.Coins{sdk.NewInt64Coin("uosmo", 7000)},
		1700000,
		chainID,
		accNums,
		accSeqs,
		senderPrivKeys,
		signerPrivKeys,
		cosignerPrivKeys,
		selectedAuthenticators,
		uint64(block.Block.Header.Height)+1,
	)

	return txBytes, nil
}
