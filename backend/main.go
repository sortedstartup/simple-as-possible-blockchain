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
