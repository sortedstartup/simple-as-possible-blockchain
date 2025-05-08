package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sortedstartup.com/simple-blockchain/backend/api"
	"sortedstartup.com/simple-blockchain/backend/blockchain"
	"sortedstartup.com/simple-blockchain/backend/common/interceptors"
	"sortedstartup.com/simple-blockchain/backend/proto"
)

func main() {
	// Create a listener on port 8080
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	bc := blockchain.NewBlockChain()

	pubKeyHex := "a0caa95ac1b9a961b804f606e86e976d561fe08956f6e32f72b6a268304e59d795c75bf4c8ea238f0e74aaddee59a9c65e55dfed22c7ceb92a10ec630a1cbb5b"
	bc.AccountBalances[pubKeyHex] = 1000

	apiServer := api.NewBlockChainAPI(bc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.PanicRecoveryInterceptor()))

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	// Register your gRPC services here
	proto.RegisterBlockchainServiceServer(grpcServer, apiServer)

	fmt.Println("Starting gRPC server on port ", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
