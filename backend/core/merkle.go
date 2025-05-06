package core

import "crypto/sha256"

func NewMerkleRoot(transactions []*Transaction) []byte {
	var hashes [][]byte
	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}
	for len(hashes) > 1 {
		var nextLevel [][]byte
		for i := 0; i < len(hashes); i += 2 {
			var combined []byte
			if i+1 == len(hashes) {
				combined = append(hashes[i], hashes[i]...)
			} else {
				combined = append(hashes[i], hashes[i+1]...)
			}
			h := sha256.Sum256(combined)
			nextLevel = append(nextLevel, h[:])
		}
		hashes = nextLevel
	}
	return hashes[0]
}
