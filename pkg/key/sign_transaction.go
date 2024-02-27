package chain

import (
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

	"github.com/osmosis-labs/osmosis/v21/app/params"
	authenticatortypes "github.com/osmosis-labs/osmosis/v21/x/authenticator/types"
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
		[]int32{},
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
	selectedAuthenticators []int32,
) ([]byte, error) {
	sigs := make([]signing.SignatureV2, len(signers))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode := gen.SignModeHandler().DefaultMode()

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range signers {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
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
	// TODO: set fee payer

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range signatures {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := gen.SignModeHandler().GetSignBytes(
			signMode,
			signerData,
			txBuilder.GetTx(),
		)
		if err != nil {
			return nil, err
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			return nil, err
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = txBuilder.SetSignatures(sigs...)
		if err != nil {
			return nil, err
		}
	}

	txBytes, err := gen.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}