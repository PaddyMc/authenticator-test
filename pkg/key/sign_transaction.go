package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/osmosis-labs/osmosis/v26/app/params"
	authenticatortypes "github.com/osmosis-labs/osmosis/v26/x/smart-account/types"
)

// SignMsg signs an sdk.Message with a given private key, account number and account sequence
// this returns a signed byte array for broadcasting to the network
func SignMsg(
	encCfg params.EncodingConfig,
	chainID string,
	msg sdk.Msg,
	priv1 *secp256k1.PrivKey,
	accNum uint64,
	accSeq uint64,
) []byte {
	txBytes, err := SignAuthenticatorMsg(
		encCfg.TxConfig,
		[]sdk.Msg{msg},
		sdk.Coins{sdk.NewInt64Coin("uosmo", 1000)},
		300000,
		chainID,
		[]uint64{accNum},
		[]uint64{accSeq},
		[]cryptotypes.PrivKey{priv1},
		[]cryptotypes.PrivKey{priv1},
		make(map[int][]cryptotypes.PrivKey),
		[]uint64{},
	)
	if err != nil {
		panic(err)
	}

	return txBytes
}

// GenTx generates a signed mock transaction.
func SignAuthenticatorMsg(
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums, accSeqs []uint64,
	signers, signatures []cryptotypes.PrivKey,
	cosigners map[int][]cryptotypes.PrivKey,
	selectedAuthenticators []uint64,
) ([]byte, error) {
	txBytes, err := SignAuthenticatorMsgWithHeight(
		gen,
		msgs,
		sdk.Coins{sdk.NewInt64Coin("uosmo", 30_000)},
		12000000,
		chainID,
		accNums,
		accSeqs,
		signers,
		signatures,
		cosigners,
		selectedAuthenticators,
		0,
	)

	return txBytes, err
}

// GenTx generates a signed mock transaction.
func SignAuthenticatorMsgWithHeight(
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums, accSeqs []uint64,
	signers, signatures []cryptotypes.PrivKey,
	cosigners map[int][]cryptotypes.PrivKey,
	selectedAuthenticators []uint64,
	timeoutHeight uint64,
) ([]byte, error) {
	sigs := make([]signing.SignatureV2, len(signers))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode := gen.SignModeHandler().DefaultMode()
	protoSignMode, _ := authsigning.APISignModeToInternal(signMode)
	//signBytes, err := authsigning.GetSignBytesAdapter(ctx, sigModeHandler, signing.SignMode_SIGN_MODE_DIRECT, signerDataOld, tx)

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range signers {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: protoSignMode,
			},
			Sequence: accSeqs[i],
		}
	}

	baseTxBuilder := gen.NewTxBuilder()

	txBuilder, ok := baseTxBuilder.(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, fmt.Errorf("expected authtx.ExtensionOptionsTxBuilder, got %T", baseTxBuilder)
	}
	if len(selectedAuthenticators) > 0 {
		value, err := types.NewAnyWithValue(&authenticatortypes.TxExtension{
			SelectedAuthenticators: selectedAuthenticators,
		})
		if err != nil {
			return nil, err
		}
		txBuilder.SetNonCriticalExtensionOptions(value)
	}

	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = txBuilder.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(feeAmt)
	txBuilder.SetGasLimit(gas)
	txBuilder.SetTimeoutHeight(timeoutHeight)
	// TODO: set fee payer

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range signatures {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}

		tx := txBuilder.GetTx()

		// signBytes, err := gen.SignModeHandler().GetSignBytes(
		//
		//	context.Background(),
		//	signMode,
		//	signerData,
		//	tx,
		//
		// )
		signBytes, _ := authsigning.GetSignBytesAdapter(
			context.Background(),
			gen.SignModeHandler(),
			protoSignMode,
			signerData,
			tx,
		)
		if err != nil {
			return nil, err
		}

		// Assuming cosigners is already initialized and populated
		var compoundSignatures [][]byte
		if value, exists := cosigners[1]; exists {
			for _, item := range value {
				sig, err := item.Sign(signBytes)
				if err != nil {
					return nil, err
				}
				compoundSignatures = append(compoundSignatures, sig)
			}
		}

		// Marshalling the array of SignatureV2 for compound signatures
		if len(compoundSignatures) >= 1 {
			compoundSigData, err := json.Marshal(compoundSignatures)
			if err != nil {
				return nil, err
			}
			sigs[i] = signing.SignatureV2{
				PubKey: signers[i].PubKey(),
				Data: &signing.SingleSignatureData{
					SignMode:  protoSignMode,
					Signature: compoundSigData,
				},
				Sequence: accSeqs[i],
			}
		} else {
			sig, err := p.Sign(signBytes)
			if err != nil {
				return nil, err
			}
			sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		}

		err = txBuilder.SetSignatures(sigs...)
		if err != nil {
			return nil, err
		}
	}
	//fmt.Println(txBuilder.GetTx())

	txBytes, err := gen.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}
