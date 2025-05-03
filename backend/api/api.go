package api

import (
	"context"

	proto "sortedstartup.com/simple-blockchain/backend/proto"
)

type BlockChainAPI struct {
	// Add your fields here
	proto.UnimplementedBlockchainServiceServer
}

func NewBlockChainAPI() *BlockChainAPI {
	return &BlockChainAPI{}
}

func (b *BlockChainAPI) SubmitTransaction(context.Context, *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	// Implement your logic here
	return &proto.SubmitTransactionResponse{
		Success: false,
		Message: "not implemented",
	}, nil
}
