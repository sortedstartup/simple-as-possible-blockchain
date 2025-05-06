package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

type Block struct {
	index        int
	timestamp    int64
	transactions []*pb.Transaction
	prevHash     string
	hash         string
	MerkleRoot   string
	Nonce        int
}

type Blockchain struct {
	Blocks          []Block
	MemoryPool      []*pb.Transaction
	AccountBalances map[string]uint64
}

func NewBlockChain() *Blockchain {
	genesis := Block{
		index:        0,
		timestamp:    time.Now().Unix(),
		transactions: []*pb.Transaction{},
		prevHash:     "",
		hash:         "",
		MerkleRoot:   "",
	}
	genesis.hash = computeHash(genesis)
	return &Blockchain{
		Blocks:          []Block{genesis},
		MemoryPool:      []*pb.Transaction{},
		AccountBalances: make(map[string]uint64),
	}
}

func computeHash(block Block) string {
	record := string(block.index) + block.MerkleRoot + string(block.timestamp)
	h := sha256.Sum256([]byte(record))

	return hex.EncodeToString(h[:])
}

func (bc *Blockchain) HandleTransaction(tx *pb.Transaction) (bool, string) {
	if bc.AccountBalances[tx.Sender] < tx.Amount {
		return false, "insufficient balance"
	}

	bc.MemoryPool = append(bc.MemoryPool, tx)

	return true, "transaction added to mempool"
}
