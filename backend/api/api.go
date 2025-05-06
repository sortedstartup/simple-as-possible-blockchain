package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"sortedstartup.com/simple-blockchain/backend/blockchain"
	proto "sortedstartup.com/simple-blockchain/backend/proto"
)

type BlockChainAPI struct {
	// Add your fields here
	proto.UnimplementedBlockchainServiceServer
	blockchain *blockchain.Blockchain
}

func NewBlockChainAPI(bc *blockchain.Blockchain) *BlockChainAPI {
	return &BlockChainAPI{
		blockchain: bc,
	}
}

func (b *BlockChainAPI) SubmitTransaction(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	// Implement your logic here

	if req == nil || req.Transaction == nil {
		return &proto.SubmitTransactionResponse{
			Success: false,
			Message: "empty transaction request",
		}, nil
	}

	tx := req.Transaction
	tx.Timestamp = time.Now().Unix()

	raw := tx.Sender + tx.Recipient + fmt.Sprintf("%d%d", tx.Amount, tx.Timestamp)
	hash := sha256.Sum256([]byte(raw))
	tx.Txid = hex.EncodeToString(hash[:])

	success, msg := b.blockchain.HandleTransaction(tx) //validates balances from AccountBalances Map

	return &proto.SubmitTransactionResponse{
		Success: success,
		Message: msg,
	}, nil

}
