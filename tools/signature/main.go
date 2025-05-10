package main

import (
	"fmt"
	"os"
)

const satoshiPrivateKey = "307702010104207271b06c11a64ca51810b5634b1cd79893d7b9b83ca64a9d5b9eae5b62115245a00a06082a8648ce3d030107a1440342000407fe15daf25dc69615abeca85494003e662dcb913bec1a7abe667be3f3f384de1d776c239d585d83965b7d634516fc9951567eb23b9ef03dc20d4e85b8707fdf"
const senderPublicKey = "07fe15daf25dc69615abeca85494003e662dcb913bec1a7abe667be3f3f384de1d776c239d585d83965b7d634516fc9951567eb23b9ef03dc20d4e85b8707fdf"

func main() {
	sender := "07fe15daf25dc69615abeca85494003e662dcb913bec1a7abe667be3f3f384de1d776c239d585d83965b7d634516fc9951567eb23b9ef03dc20d4e85b8707fdf"
	recipient := "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e32f72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b"
	amount := int64(1000)
	timestamp := int64(1746890553)

	privKey, err := decodePrivateKey(satoshiPrivateKey)
	if err != nil {
		fmt.Println("Error decoding private key:", err)
		os.Exit(1)
	}

	signature, err := signTransaction(privKey, sender, recipient, amount, timestamp)
	if err != nil {
		fmt.Println("Error signing transaction:", err)
		os.Exit(1)
	}

	fmt.Println("Signature:", signature)
}
