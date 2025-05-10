package helpers

import (
	"encoding/hex"
	"errors"
)

/*
*
Public key by ECDSA
-> 32 bytes(X) + 32 bytes(y) = 64 bytes = 128 chars
-> x and y parts of that public should be on the elliptic curve,
*/
func ValidateRawPublicKey(hexPub string) error {
	if len(hexPub) == 128 {
		hexPub = "04" + hexPub
	} else if len(hexPub) != 130 {
		return errors.New("invalid public key lenght")
	}

	_, err := hex.DecodeString(hexPub)
	if err != nil {
		return errors.New("invalid hex string")
	}

	return nil
}
