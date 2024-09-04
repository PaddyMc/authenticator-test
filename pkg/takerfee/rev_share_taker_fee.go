package poolmanager

import (
	"log"

	grpc "google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/osmosis-labs/osmosis/v26/app/params"
)

func SetUpTakerFeeRevShare(
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

	return nil
}
