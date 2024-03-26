package validator

import (
	"context"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v23/app/params"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
)

func CreateValidator(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig,
	chainID string,
	valSigningKey *secp256k1.PrivKey,
	authKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
) error {
	pk1 := ed25519.GenPrivKey().PubKey()
	priv1 := valSigningKey

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	stakingClient := stakingtypes.NewQueryClient(conn)

	description := stakingtypes.Description{
		Moniker: "jail",
		Details: "slash me",
	}
	comRates := stakingtypes.CommissionRates{
		Rate:          osmomath.NewDec(1),
		MaxRate:       osmomath.NewDec(1),
		MaxChangeRate: osmomath.NewDec(1),
	}

	accAddress := sdk.AccAddress(priv1.PubKey().Address())
	valAddr := sdk.ValAddress(accAddress)

	pkAny, err := codectypes.NewAnyWithValue(pk1)
	if err != nil {
		return err
	}

	msg := &stakingtypes.MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  accAddress.String(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Value:             sdk.NewCoin("uosmo", osmomath.NewInt(2000)),
		Commission:        comRates,
		MinSelfDelegation: osmomath.NewInt(2000),
	}

	err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		encCfg,
		ac,
		txClient,
		chainID,
		[]sdk.Msg{msg},
		[]uint64{},
	)
	if err != nil {
		return err
	}
	validator, err := stakingClient.Validator(
		context.Background(),
		&stakingtypes.QueryValidatorRequest{
			ValidatorAddr: valAddr.String(),
		},
	)

	log.Printf("Creating a validator with address %s", validator.Validator.OperatorAddress)

	return nil
}
