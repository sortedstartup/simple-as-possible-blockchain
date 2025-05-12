package api

import (
	"context"
	"fmt"

	"sortedstartup.com/simple-blockchain/backend/blockchain"
	"sortedstartup.com/simple-blockchain/backend/helpers"
	proto "sortedstartup.com/simple-blockchain/backend/proto"
)

type BlockChainAPI struct {
	proto.UnimplementedBlockchainServiceServer
	blockchain *blockchain.Blockchain
}

func NewBlockChainAPI(bc *blockchain.Blockchain) *BlockChainAPI {
	return &BlockChainAPI{
		blockchain: bc,
	}
}

func (b *BlockChainAPI) SubmitTransaction(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {

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
