package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "api-bridge/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to debug gRPC server
	conn, err := grpc.Dial("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewAPIBridgeServiceClient(conn)

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test GetAccounts
	fmt.Println("Testing GetAccounts...")
	accountsResp, err := client.GetAccounts(ctx, &pb.GetAccountsRequest{})
	if err != nil {
		log.Printf("GetAccounts failed: %v", err)
	} else {
		fmt.Printf("Found %d accounts\n", accountsResp.Count)
		for _, account := range accountsResp.Accounts {
			fmt.Printf("  - %s (%s): %s\n", account.Name, account.KeyType, account.Address)
		}
	}

	// Test RegisterComponent
	fmt.Println("\nTesting RegisterComponent...")
	registerResp, err := client.RegisterComponent(ctx, &pb.RegisterComponentRequest{
		Creator:       "alice",
		ComponentData: "test-battery-module",
		Context:       "test-context",
	})
	if err != nil {
		log.Printf("RegisterComponent failed: %v", err)
	} else {
		fmt.Printf("Component registered: %s\n", registerResp.ComponentId)
		fmt.Printf("Transaction hash: %s\n", registerResp.Txhash)
	}

	fmt.Println("\n✅ Debug gRPC client test completed successfully!")
}
