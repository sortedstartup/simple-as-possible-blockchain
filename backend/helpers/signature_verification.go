package helpers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

func VerifySignature(senderHex, recipientHex string, amount, timestamp int64, signatureHex string) error {
	fmt.Printf("Signature length: %d\n", len(signatureHex))

	if signatureHex == "" {
		return fmt.Errorf("empty signature")
	}

	pubKeyBytes, err := hex.DecodeString(senderHex)
	if err != nil {
		return fmt.Errorf("invalid sender public key: %v", err)
	}

	if len(pubKeyBytes) != 64 {
		return fmt.Errorf("invalid sender public key length: %d, expected 64", len(pubKeyBytes))
	}

	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])
	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	raw := senderHex + recipientHex + fmt.Sprintf("%d%d", amount, timestamp)
	hash := sha256.Sum256([]byte(raw))
	fmt.Printf(" Raw message: %s...\n", raw[:40])
	fmt.Printf(" Hash: %x\n", hash[:])

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("invalid signature format: %v", err)
	}
	fmt.Printf(" sigBytes: %x\n", sigBytes)

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])
	fmt.Printf("DEBUG VerifySignature: r=%s, s=%s\n", r.String(), s.String())

	if !ecdsa.Verify(&publicKey, hash[:], r, s) {
		return errors.New("signature verification failed")
	}

	fmt.Println("DEBUG VerifySignature: Signature verified successfully!")
	return nil
}

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
