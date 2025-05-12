package blockchain

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"time"

	"sortedstartup.com/simple-blockchain/backend/helpers"
	pb "sortedstartup.com/simple-blockchain/backend/proto"
)

//go:embed keys/satoshi.publickey
var SatoshiPublicKey string

//go:embed keys/satoshi.privatekey
var SatoshiPrivateKey string

type Block struct {
	index        int
	timestamp    int64
	transactions []*pb.Transaction
	prevHash     string
	hash         string
	MerkleRoot   string
	Nonce        int
}

type UTXO struct {
	Txid      string
	Index     int
	Amount    uint64
	Recipient string
}

type Blockchain struct {
	Blocks     []Block
	MemoryPool []*pb.Transaction
	// AccountBalances map[string]uint64
	UTXOSet map[string]UTXO
}

func NewBlockChain() *Blockchain {
	return createGenesisBlock()
}

func createGenesisBlock() *Blockchain {
	bc := &Blockchain{
		Blocks:     []Block{},
		MemoryPool: []*pb.Transaction{},
		// AccountBalances: make(map[string]uint64),
		UTXOSet: make(map[string]UTXO),
	}

	// satoshiPubKey := SatoshiPublicKey
	// bc.AccountBalances[satoshiPubKey] = 100000 // coinbase transaction

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

	coinbaseUTXO := UTXO{
		Txid:      "genesis",
		Index:     0,
		Amount:    100000,
		Recipient: SatoshiPublicKey,
	}
	bc.UTXOSet["genesis:0"] = coinbaseUTXO
	return bc
}

func computeHash(block Block) string {
	record := fmt.Sprintf("%d%s%d", block.index, block.MerkleRoot, block.timestamp)
	h := sha256.Sum256([]byte(record))

	return hex.EncodeToString(h[:])
}

func (bc *Blockchain) HandleTransaction(tx *pb.Transaction) (bool, string) {

	// Generate transaction id
	raw := tx.Sender + tx.Recipient + fmt.Sprintf("%d%d", tx.Amount, tx.Timestamp)
	hash := sha256.Sum256([]byte(raw))
	tx.Txid = hex.EncodeToString(hash[:])

	// Verify Signature
	err := helpers.VerifySignature(tx.Sender, tx.Recipient, int64(tx.Amount), tx.Timestamp, string(tx.Signature))
	if err != nil {
		return false, "signature verification failed: " + err.Error()
	}

	var selected []UTXO
	var total uint64 = 0
	for _, utxo := range bc.UTXOSet {
		if utxo.Recipient == tx.Sender {
			selected = append(selected, utxo)
			total += utxo.Amount
			if total >= tx.Amount {
				break
			}

		}
	}

	if total < tx.Amount {
		return false, "insufficient balance (from UTXOs)"
	}

	for _, utxo := range selected {
		key := fmt.Sprintf("%s:%d", utxo.Txid, utxo.Index)
		delete(bc.UTXOSet, key)
	}

	recipientUTXO := UTXO{
		Txid:      tx.Txid,
		Index:     0,
		Amount:    tx.Amount,
		Recipient: tx.Recipient,
	}
	bc.UTXOSet[fmt.Sprintf("%s:%d", tx.Txid, 0)] = recipientUTXO

	change := total - tx.Amount
	if change > 0 {
		changeUTXO := UTXO{
			Txid:      tx.Txid,
			Index:     1,
			Amount:    change,
			Recipient: tx.Sender,
		}
		bc.UTXOSet[fmt.Sprintf("%s:%d", tx.Txid, 1)] = changeUTXO
	}

	// check account balanece
	// if bc.AccountBalances[tx.Sender] < tx.Amount {
	// 	return false, "insufficient balance"
	// }

	//will update this logic with UTXO
	// bc.AccountBalances[tx.Sender] -= tx.Amount
	// bc.AccountBalances[tx.Recipient] += tx.Amount

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
