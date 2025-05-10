package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"sortedstartup.com/simple-blockchain/backend/blockchain"
	"sortedstartup.com/simple-blockchain/backend/helpers"
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
	// tx.Timestamp = time.Now().Unix()

	if err := helpers.ValidateRawPublicKey(tx.Sender); err != nil {
		return &proto.SubmitTransactionResponse{
			Success: false,
			Message: "invalid sender public key: " + err.Error(),
		}, nil
	}
	if err := helpers.ValidateRawPublicKey(tx.Recipient); err != nil {
		return &proto.SubmitTransactionResponse{
			Success: false,
			Message: "invalid recipient public key: " + err.Error(),
		}, nil
	}

	raw := tx.Sender + tx.Recipient + fmt.Sprintf("%d%d", tx.Amount, tx.Timestamp)
	hash := sha256.Sum256([]byte(raw))
	fmt.Printf(" api hash: %x\n", hash[:])

	tx.Txid = hex.EncodeToString(hash[:])

	success, msg := b.blockchain.HandleTransaction(tx) //validates balances from AccountBalances Map

	if success {
		fmt.Println("Transaction accepted. Mempool now:")
		b.blockchain.PrintMemPool()
	}

	return &proto.SubmitTransactionResponse{
		Success: success,
		Message: msg,
	}, nil

}
