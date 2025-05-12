package api

import (
	"context"
	"testing"

	"sortedstartup.com/simple-blockchain/backend/blockchain"
	"sortedstartup.com/simple-blockchain/backend/helpers"
	proto "sortedstartup.com/simple-blockchain/backend/proto"
)

func TestSubmitTransaction(t *testing.T) {

	// Initialize a mock blockchain
	mockBlockchain := blockchain.NewBlockChain()
	api := NewBlockChainAPI(mockBlockchain)

	// ---------------------------------
	// Test case: valid transaction
	//------------------------------------
	req := &proto.SubmitTransactionRequest{
		Transaction: &proto.Transaction{
			Sender: blockchain.SatoshiPublicKey,
			//TODO: should error out on random string
			Recipient: "bob",
			Amount:    10,
		},
	}

	sig, err := signTransaction(blockchain.SatoshiPrivateKey, req.Transaction)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	req.Transaction.Signature = sig

	resp, err := api.SubmitTransaction(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// ---------------------------------=
	// Test case: nil request
	//------------------------------------
	resp, err = api.SubmitTransaction(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Success {
		t.Errorf("expected success to be false, got true")
	}
	if resp.Message != "empty transaction request" {
		t.Errorf("unexpected message: %v", resp.Message)
	}

	// Test case: invalid sender public key
	req = &proto.SubmitTransactionRequest{
		Transaction: &proto.Transaction{
			Sender:    "invalid_sender",
			Recipient: "valid_recipient",
			Amount:    10,
		},
	}
	resp, err = api.SubmitTransaction(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Success {
		t.Errorf("expected success to be false, got true")
	}

}

// Helper method to sign a transaction
func signTransaction(privateKeyHex string, tx *proto.Transaction) ([]byte, error) {
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
