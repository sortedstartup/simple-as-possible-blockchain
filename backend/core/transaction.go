package core

import (
	"crypto/sha256"
	"fmt"
)

type TxInput struct {
	TxID      []byte
	OutIndex  int
	Signature string
	PulicKey  string
}

type TxOutput struct {
	Value      int64
	PubKeyHash string
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) Hash() []byte {
	h := sha256.Sum256([]byte(fmt.Sprintf("%v", tx)))
	return h[:]
}

func NewCoinbaseTx(to string, value int64) *Transaction {
	txOutput := TxOutput{Value: value, PubKeyHash: to}
	tx := &Transaction{
		Inputs:  []TxInput{{}}, //no input in coinbase transaction
		Outputs: []TxOutput{txOutput},
	}
	tx.ID = tx.Hash()

	return tx
}
