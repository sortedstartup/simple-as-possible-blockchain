package merkle

import (
	"crypto/sha256"
	"encoding/hex"

	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

func computeMerkleRoot(transactions []*pb.Transaction) string {

	if len(transactions) == 0 {
		return ""
	}

	var txnHashes []string

	for _, txn := range transactions {
		data := txn.Txid + txn.Sender + txn.Recipient
		h := sha256.Sum256([]byte(data))
		txnHashes = append(txnHashes, hex.EncodeToString(h[:]))
	}

	for len(txnHashes) > 1 {
		var newLevel []string
		for i := 0; i < len(txnHashes); i += 2 {
			if i+1 < len(txnHashes) {
				combined := txnHashes[i] + txnHashes[i+1]
				h := sha256.Sum256([]byte(combined))
				newLevel = append(newLevel, hex.EncodeToString(h[:]))
			} else {
				newLevel = append(newLevel, txnHashes[i])
			}
		}
		txnHashes = newLevel
	}
	return txnHashes[0]
}
