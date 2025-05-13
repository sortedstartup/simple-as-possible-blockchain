package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"sortedstartup.com/simple-blockchain/backend/helpers"
	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

// TODO
// - Transfer to self ?

// Helper method to sign a transaction
func signTransaction(privateKeyHex string, tx *pb.Transaction) ([]byte, error) {
	privKey, err := helpers.ConvertHexToPrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}

	signature, err := helpers.SignTransaction(privKey, tx.Sender, tx.Recipient, int64(tx.Amount), tx.Timestamp)
	if err != nil {
		return nil, err
	}

	return []byte(signature), nil
}

func generateRecipientKeyHex() (priv *ecdsa.PrivateKey, pubHex string, err error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, "", err
	}
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return privKey, fmt.Sprintf("%x", pubKey), nil
}
func TestUTXOTransactionSuccess(t *testing.T) {
	bc := NewBlockChain()

	_, recipientHex, err := generateRecipientKeyHex()
	if err != nil {
		t.Fatalf("failed to generate recipient key: %v", err)
	}

	bc.UTXOSet = map[string]UTXO{
		"init:0": {
			Txid:      "init",
			Index:     0,
			Amount:    70,
			Recipient: SatoshiPublicKey,
		},
	}

	tx := &pb.Transaction{
		Sender:    SatoshiPublicKey,
		Recipient: recipientHex,
		Amount:    50,
		Timestamp: time.Now().Unix(),
	}

	sig, err := signTransaction(SatoshiPrivateKey, tx)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	tx.Signature = sig

	success, msg := bc.HandleTransaction(tx)
	if !success {
		t.Fatalf("expected transaction to succeed, got error: %s", msg)
	}

	if len(bc.MemoryPool) != 1 {
		t.Errorf("expected 1 transaction in mempool, got %d", len(bc.MemoryPool))
	}

	var gotRecipient, gotChange bool
	for _, utxo := range bc.UTXOSet {
		if utxo.Recipient == tx.Recipient && utxo.Amount == 50 {
			gotRecipient = true
		}
		if utxo.Recipient == tx.Sender && utxo.Amount == 20 {
			gotChange = true
		}
	}
	if !gotRecipient || !gotChange {
		t.Errorf("expected both recipient and change UTXOs, got: %+v", bc.UTXOSet)
	}
}

func TestUTXOTransactionInsufficientBalance(t *testing.T) {
	bc := NewBlockChain()

	_, recipientHex, err := generateRecipientKeyHex()
	if err != nil {
		t.Fatalf("failed to generate recipient key: %v", err)
	}

	bc.UTXOSet = map[string]UTXO{
		"init:0": {
			Txid:      "init",
			Index:     0,
			Amount:    30,
			Recipient: SatoshiPublicKey,
		},
	}

	tx := &pb.Transaction{
		Sender:    SatoshiPublicKey,
		Recipient: recipientHex,
		Amount:    50,
		Timestamp: time.Now().Unix(),
	}

	sig, err := signTransaction(SatoshiPrivateKey, tx)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	tx.Signature = sig

	success, msg := bc.HandleTransaction(tx)
	if success {
		t.Fatal("expected transaction to fail due to insufficient UTXO balance")
	}
	if msg != "insufficient balance (from UTXOs)" {
		t.Errorf("unexpected error message: %s", msg)
	}

	if len(bc.MemoryPool) != 0 {
		t.Errorf("mempool should be empty, got %d", len(bc.MemoryPool))
	}
}
