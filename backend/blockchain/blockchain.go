package blockchain

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"time"

	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

//go:embed keys/satoshi.publickey
var satoshiPublicKey string

//go:embed keys/satoshi.privatekey
var satoshiPrivateKey string

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
	return createGenesisBlock()
}

func createGenesisBlock() *Blockchain {
	bc := &Blockchain{
		Blocks:          []Block{},
		MemoryPool:      []*pb.Transaction{},
		AccountBalances: make(map[string]uint64),
	}

	satoshiPubKey := satoshiPublicKey
	bc.AccountBalances[satoshiPubKey] = 1000 // coinbase transaction

	genesis := Block{
		index:        0,
		timestamp:    time.Now().Unix(),
		transactions: []*pb.Transaction{},
		prevHash:     "",
		hash:         "",
		MerkleRoot:   "",
	}
	genesis.hash = computeHash(genesis)
	bc.Blocks = append(bc.Blocks, genesis)
	return bc
}

func computeHash(block Block) string {
	record := fmt.Sprintf("%d%s%d", block.index, block.MerkleRoot, block.timestamp)
	h := sha256.Sum256([]byte(record))

	return hex.EncodeToString(h[:])
}

func (bc *Blockchain) HandleTransaction(tx *pb.Transaction) (bool, string) {
	if bc.AccountBalances[tx.Sender] < tx.Amount {
		return false, "insufficient balance"
	}

	//will update this logic with UTXO
	bc.AccountBalances[tx.Sender] -= tx.Amount
	bc.AccountBalances[tx.Recipient] += tx.Amount

	bc.MemoryPool = append(bc.MemoryPool, tx)

	return true, "transaction added to mempool"
}

func (bc *Blockchain) PrintMemPool() {
	fmt.Println("Current Mempool:")
	if len(bc.MemoryPool) == 0 {
		fmt.Println("  (empty)")
		return
	}

	for i, tx := range bc.MemoryPool {
		fmt.Printf("  #%d → TxID: %s | From: %s | To: %s | Amount: %d\n",
			i+1, tx.Txid, tx.Sender, tx.Recipient, tx.Amount)
	}
}
