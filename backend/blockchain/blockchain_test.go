package blockchain

import (
	"testing"

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

func TestSufficientBalance(t *testing.T) {

	bc := NewBlockChain()

	tx := &pb.Transaction{
		Sender:    SatoshiPublicKey,
		Recipient: "bob",
		Amount:    500,
	}

	// Use helper method to sign the transaction
	signature, err := signTransaction(SatoshiPrivateKey, tx)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature

	success, msg := bc.HandleTransaction(tx)

	if !success {
		t.Errorf("expected transaction to succeed, got message: %s", msg)
	}

	if len(bc.MemoryPool) != 1 {
		t.Errorf("expected memory pool to have 1 transaction, got %d", len(bc.MemoryPool))
	}
}

func TestInsufficientBalance(t *testing.T) {
	bc := NewBlockChain()

	tx := &pb.Transaction{
		Sender:    SatoshiPublicKey,
		Recipient: "bob",
		Amount:    100000000,
	}

	// Use helper method to sign the transaction
	signature, err := signTransaction(SatoshiPrivateKey, tx)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature

	success, msg := bc.HandleTransaction(tx)

	if success {
		t.Errorf("expected transaction to fail due to insufficient balance")
	}

	if msg != "insufficient balance" {
		t.Errorf("unexpected error message: %s", msg)
	}

	if len(bc.MemoryPool) != 0 {
		t.Errorf("memory pool should be empty, got %d transactions", len(bc.MemoryPool))
	}
}
