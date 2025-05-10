package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"math/big"
)

func decodePrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	privBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %v", err)
	}

	privKey, err := x509.ParseECPrivateKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("parse EC private key: %v", err)
	}

	return privKey, nil
}

type ecdsaSignature struct {
	R, S *big.Int
}

func signTransaction(priv *ecdsa.PrivateKey, sender, recipient string, amount, timestamp int64) (string, error) {
	raw := fmt.Sprintf("%s%s%d%d", sender, recipient, amount, timestamp)
	hash := sha256.Sum256([]byte(raw))

	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", fmt.Errorf("signing failed: %v", err)
	}

	sigBytes, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", fmt.Errorf("failed to encode signature: %v", err)
	}

	return hex.EncodeToString(sigBytes), nil
}
