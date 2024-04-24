package validator

import (
	"context"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/osmosis-labs/osmosis/v24/app/params"
)

func GetValidatorDelegations(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	valSigningKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	// set up all clients
	stakingClient := stakingtypes.NewQueryClient(conn)
	valOpAddr := "osmovaloper1dnmz4yzv73lr3lmauuaa0wpwn8zm8s20lg3pg9"

	eleValAddr, _ := sdk.ValAddressFromBech32(valOpAddr)
	//eleValAddr, _ := sdk.ValAddressFromBech32("osmovaloper1tv9wnreg9z5qlxyte8526n7p3tjasndede2kj9")

	limit := uint64(1000)
	key := []byte{}
	validatorDelegations, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&stakingtypes.QueryValidatorDelegationsRequest{
			ValidatorAddr: eleValAddr.String(),
			Pagination: &query.PageRequest{
				Key:   key,
				Limit: limit,
			},
		},
	)
	if err != nil {
		return err
	}
	//fmt.Println(validatorDelegations)

	log.Printf("Number of delegators for validator %s %d", valOpAddr, len(validatorDelegations.DelegationResponses))

	return nil
}
