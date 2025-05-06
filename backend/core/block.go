package core

import (
	"crypto/sha256"
	"encoding/json"
)

type Block struct {
	PrevHash     []byte
	Hash         []byte
	Transactions []*Transaction
	MerkleRoot   []byte
	Height       int //like block number -> for longest chain rule
}

func (b *Block) CalculateHash() []byte {
	data, _ := json.Marshal(b.Transactions)
	hash := sha256.Sum256(append(b.PrevHash, data...))
	return hash[:]
}

func NewBlock(transactions []*Transaction, prevHash []byte, height int) *Block {
	merkle := NewMerkleRoot(transactions)

	block := &Block{
		PrevHash:     prevHash,
		Transactions: transactions,
		MerkleRoot:   merkle,
		Height:       height,
	}
	block.Hash = block.CalculateHash()
	return block
}
