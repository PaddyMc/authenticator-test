package config

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v21/app"
	"github.com/osmosis-labs/osmosis/v21/app/params"

	grpc_chain "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	grpc "google.golang.org/grpc"
)

type SeedConfig struct {
	ChainID        string
	GRPCConnection *grpc.ClientConn
	EncodingConfig params.EncodingConfig
	Keys           []*secp256k1.PrivKey
	DenomMap       map[string]string
}

func SetUp(
	chainID string,
	grpcAddr string,
	keysHex []string,
	DefaultDenoms map[string]string,
) SeedConfig {
	// set up all clients
	conn := grpc_chain.CreateGRPCConnection(grpcAddr)
	encCfg := app.MakeEncodingConfig()
	keys := []*secp256k1.PrivKey{}

	// Decode and add keys to the slice
	for _, keyHex := range keysHex {
		bz, err := hex.DecodeString(keyHex)
		fmt.Println(len(bz))
		fmt.Println(bz)
		if err != nil {
			fmt.Printf("Error decoding hex string: %v\n", err)
			continue
		}

		privKey := &secp256k1.PrivKey{Key: bz}
		accAddress := sdk.AccAddress(privKey.PubKey().Address())
		log.Println("Account: ", accAddress.String())
		keys = append(keys, privKey)
	}

	seedConfig := SeedConfig{
		ChainID:        chainID,
		GRPCConnection: conn,
		EncodingConfig: encCfg,
		Keys:           keys,
		DenomMap:       DefaultDenoms,
	}

	return seedConfig
}
