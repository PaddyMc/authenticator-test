package poolmanager

import (
	"context"
	"log"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	chaingrpc "github.com/osmosis-labs/autenticator-test/pkg/grpc"
	"github.com/osmosis-labs/osmosis/v26/app/params"
)

func ClaimLimitOrders(
	conn *grpc.ClientConn,
	encCfg params.EncodingConfig, chainID string,
	fromKey *secp256k1.PrivKey,
	signerKey *secp256k1.PrivKey,
	cosignersKeys map[int][]cryptotypes.PrivKey,
	selectedAuthenticator []int32,
	codeId uint64,
) error {
	accAddress := sdk.AccAddress(fromKey.PubKey().Address())
	log.Printf("Batch claiming limit orders from address %s", accAddress)

	// set up all clients
	txClient := txtypes.NewServiceClient(conn)
	ac := auth.NewQueryClient(conn)
	wasmClient := wasmtypes.NewQueryClient(conn)

	// Query the orderbook code for all contract addresses
	// https://celatone.osmosis.zone/osmosis-1/codes/885
	contractsByCodeResp, err := wasmClient.ContractsByCode(
		context.Background(),
		&wasmtypes.QueryContractsByCodeRequest{
			CodeId: codeId,
		},
		// TODO: maybe add some pagination
	)
	if err != nil {
		return err
	}

	for _, contractAddress := range contractsByCodeResp.Contracts {
		log.Printf("Claiming limit orders for contract %s", contractAddress)

		//	quertSmartResp, err := wasmClient.QuertSmart(
		//		context.Background(),
		//		&wasmtypes.QueryContractsByCodeRequest{
		//			CodeId: codeId,
		//		},
		//		// TODO: maybe add some pagination
		//	)
		if err != nil {
			return err
		}
		// TODO: get all orders by tick
		// NOTE: https://github.com/osmosis-labs/orderbook/blob/main/contracts/sumtree-orderbook/src/msg.rs#L69-L143
		// e.g
		//     #[returns(TicksResponse)]
		// 		 AllTicks {
		// 		     /// The tick id to start after for pagination (inclusive)
		// 		     start_from: Option<i64>,
		// 		     /// A max tick id to end at if limit is not reached/provided (inclusive)
		// 		     end_at: Option<i64>,
		// 		     /// The limit for amount of items to return
		// 		     limit: Option<usize>,
		// 		 },

		// e.g
		//     #[returns(TicksResponse)]
		// 		 TicksById { tick_ids: Vec<i64> },
		// e.g
		//    #[returns(OrdersResponse)]
		//    OrdersByTick {
		//        tick_id: i64,
		//        start_from: Option<u64>,
		//        end_at: Option<u64>,
		//        limit: Option<u64>,
		//    },

		// TODO: build the binary execute message
		// ExecuteMsg needed
		//    BatchClaim {
		//      orders: Vec<(i64, u64)>,
		// },

		contractMsg := []byte("{batch_claim: { orders: []}}")

		msg := &wasmtypes.MsgExecuteContract{
			Sender:   accAddress.String(),
			Contract: contractAddress,
			Msg:      contractMsg,
			Funds:    nil,
		}

		// Broadcast message
		err = chaingrpc.SignAndBroadcastAuthenticatorMsgMultiSigners(
			[]cryptotypes.PrivKey{signerKey},
			[]cryptotypes.PrivKey{signerKey},
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
	}

	return nil
}
