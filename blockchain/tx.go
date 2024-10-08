package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/tensor-programming/golang-blockchain/wallet"
)

// the owner of transaction when output is yours or you have been mentioned in one of the inputs

// TxOutput are invisible, so you can't split the value. so If 5 out of 10 is needed. Public key will be hashed.
type TxOutput struct {
	Value      int    // value in token
	PubKeyHash []byte // public key is a value needed to unlock tokens that stored in value.
}

type TxOutputs struct {
	Outputs []TxOutput
}

// TxInput are just reference to given TxOutput
type TxInput struct {
	ID        []byte // It's transactionID that have given output.
	Out       int    // it's the index where this output appears. So if you want to display ID with index 2 then the value will be 2.
	Signature []byte // it's a signature to use for the pubKey. for now Signature and pub key are same value.
	PubKey    []byte
}

// NewTxOutput is a new command as caller will pass amount and the address.
func NewTxOutput(amount int, address string) *TxOutput {
	txOut := &TxOutput{Value: amount}
	txOut.Lock([]byte(address))

	return txOut
}

// UsesKey is validation to check if tx input could be unlocked.
func (in *TxInput) UsesKey(pubKeyHas []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Equal(lockingHash, pubKeyHas)
}

// Lock set publicKeyHashed after trimming version, and checksum.
func (out *TxOutput) Lock(address []byte) {
	out.PubKeyHash = PubKeyHash(address)
}

func PubKeyHash(address []byte) []byte {
	pubKeyHash := wallet.Base58Encode(address)
	return pubKeyHash[1 : len(pubKeyHash)-4]
}

// IsLockedWithHash is validation to check if tx output could be unlocked.
func (out *TxOutput) IsLockedWithHash(pubKeyHash []byte) bool {
	given := string(out.PubKeyHash)
	passed := string(pubKeyHash)

	return given == passed
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)

	return outputs
}
