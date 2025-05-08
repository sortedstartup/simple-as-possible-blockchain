package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

func main() {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	fmt.Println("Public Key:", hex.EncodeToString(pub))
	encodedPriv, _ := x509.MarshalECPrivateKey(priv)
	fmt.Println("Private Key:", hex.EncodeToString(encodedPriv))
}
