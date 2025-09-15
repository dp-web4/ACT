package main

import (
	"context"
	"log"
	"net"

	pb "api-bridge/proto"

	"google.golang.org/grpc"
)

type DebugServer struct {
	pb.UnimplementedAPIBridgeServiceServer
}

func (s *DebugServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	log.Println("GetAccounts called!")
	return &pb.GetAccountsResponse{
		Accounts: []*pb.Account{
			{
				Name:    "debug-account",
				Address: "debug-address",
				KeyType: "debug-key",
			},
		},
		Count: 1,
	}, nil
}

func (s *DebugServer) RegisterComponent(ctx context.Context, req *pb.RegisterComponentRequest) (*pb.RegisterComponentResponse, error) {
	log.Println("RegisterComponent called!")
	return &pb.RegisterComponentResponse{
		ComponentId:       "debug-component-id",
		ComponentIdentity: "debug-identity",
		ComponentData:     req.ComponentData,
		Context:           req.Context,
		Creator:           req.Creator,
		LctId:             "debug-lct-id",
		Status:            "registered",
		Txhash:            "debug-txhash",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAPIBridgeServiceServer(grpcServer, &DebugServer{})

	log.Printf("Debug gRPC server starting on port 9091")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
