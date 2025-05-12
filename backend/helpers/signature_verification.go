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

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("invalid signature format: %v", err)
	}

	result := ecdsa.VerifyASN1(&publicKey, hash[:], sigBytes)
	if result {
		fmt.Println("DEBUG VerifySignature: Signature verified successfully!")
		return nil
	} else {
		fmt.Println("DEBUG VerifySignature: Signature verification failed!")
		return errors.New("signature verification failed")
	}

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

func SignTransaction(priv *ecdsa.PrivateKey, sender, recipient string, amount, timestamp int64) (string, error) {
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

// ConvertHexToPrivateKey converts a hex-encoded string into an *ecdsa.PrivateKey.
func ConvertHexToPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %v", err)
	}

	privKey, err := x509.ParseECPrivateKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %v", err)
	}

	return privKey, nil
}
