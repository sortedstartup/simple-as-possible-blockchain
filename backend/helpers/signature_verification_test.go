package helpers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestSignatureVerification(t *testing.T) {

	// ----------------------------------------------------
	// Generate a new key pair for the sender
	// ----------------------------------------------------
	senderPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Error generating sender key pair: %v", err)
	}
	senderPubKey := append(senderPrivKey.PublicKey.X.Bytes(), senderPrivKey.PublicKey.Y.Bytes()...)
	sender := fmt.Sprintf("%x", senderPubKey)

	// ----------------------------------------------------
	// Recipient key pair for reciever
	// ----------------------------------------------------
	recipientPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Error generating recipient key pair: %v", err)
	}
	recipientPubKey := append(recipientPrivKey.PublicKey.X.Bytes(), recipientPrivKey.PublicKey.Y.Bytes()...)
	recipient := fmt.Sprintf("%x", recipientPubKey)

	// TODO: definetly not hardcode this, do more generic tests
	amount := int64(1000)
	timestamp := int64(1746890553)

	signature, err := SignTransaction(senderPrivKey, sender, recipient, amount, timestamp)
	if err != nil {
		t.Fatalf("Error signing transaction: %v", err)
	}

	err = VerifySignature(sender, recipient, amount, timestamp, signature)
	if err != nil {
		t.Errorf("could not verify signature")
	}
}
