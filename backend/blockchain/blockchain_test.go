package blockchain

import (
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

func TestUTXOTransactionSuccess(t *testing.T) {
	bc := NewBlockChain()

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
		Recipient: "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e32f72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b",
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
		Recipient: "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e32f72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b",
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
