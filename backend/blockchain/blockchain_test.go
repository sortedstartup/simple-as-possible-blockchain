package blockchain

import (
	"testing"

	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

func SufficientBalance(t *testing.T) {

	bc := NewBlockChain()
	bc.AccountBalances["Alice"] = 1000

	tx := &pb.Transaction{
		Sender:    "alice",
		Recipient: "bob",
		Amount:    500,
	}

	success, msg := bc.HandleTransaction(tx)

	if !success {
		t.Errorf("expected transaction to succeed, got message: %s", msg)
	}

	if len(bc.MemoryPool) != 1 {
		t.Errorf("expected memory pool to have 1 transaction, got %d", len(bc.MemoryPool))
	}
}

func InsufficientBalance(t *testing.T) {
	bc := NewBlockChain()
	bc.AccountBalances["alice"] = 1000

	tx := &pb.Transaction{
		Sender:    "alice",
		Recipient: "bob",
		Amount:    2000,
	}

	success, msg := bc.HandleTransaction(tx)

	if !success {
		t.Errorf("expected transaction to fail due to insufficient balance")
	}

	if msg != "insufficient balance" {
		t.Errorf("unexpected error message: %s", msg)
	}

	if len(bc.MemoryPool) != 0 {
		t.Errorf("memory pool should be empty, got %d transactions", len(bc.MemoryPool))
	}
}
