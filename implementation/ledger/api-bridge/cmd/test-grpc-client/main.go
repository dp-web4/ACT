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
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	// Test GetComponent
	if registerResp != nil && registerResp.ComponentId != "" {
		fmt.Println("\nTesting GetComponent...")
		componentResp, err := client.GetComponent(ctx, &pb.GetComponentRequest{
			ComponentId: registerResp.ComponentId,
		})
		if err != nil {
			log.Printf("GetComponent failed: %v", err)
		} else {
			fmt.Printf("Component data: %s\n", componentResp.ComponentData)
			fmt.Printf("Creator: %s\n", componentResp.Creator)
		}
	}

	// Test CreateLCT
	fmt.Println("\nTesting CreateLCT...")
	lctResp, err := client.CreateLCT(ctx, &pb.CreateLCTRequest{
		Creator:    "alice",
		ComponentA: "battery-001",
		ComponentB: "motor-001",
		Context:    "race-car-pairing",
		ProxyId:    "proxy-001",
	})
	if err != nil {
		log.Printf("CreateLCT failed: %v", err)
	} else {
		fmt.Printf("LCT created: %s\n", lctResp.LctId)
		fmt.Printf("Transaction hash: %s\n", lctResp.Txhash)
		fmt.Printf("LCT Key Half: %s\n", lctResp.LctKeyHalf)
		fmt.Printf("Device Key Half: %s\n", lctResp.DeviceKeyHalf)
	}

	// Test InitiatePairing
	fmt.Println("\nTesting InitiatePairing...")
	pairingResp, err := client.InitiatePairing(ctx, &pb.InitiatePairingRequest{
		Creator:            "alice",
		ComponentA:         "battery-001",
		ComponentB:         "motor-001",
		OperationalContext: "race-car-operation",
		ProxyId:            "proxy-001",
		ForceImmediate:     false,
	})
	if err != nil {
		log.Printf("InitiatePairing failed: %v", err)
	} else {
		fmt.Printf("Pairing initiated: %s\n", pairingResp.ChallengeId)
		fmt.Printf("Transaction hash: %s\n", pairingResp.Txhash)
	}

	// Test CompletePairing
	if pairingResp != nil && pairingResp.ChallengeId != "" {
		fmt.Println("\nTesting CompletePairing...")
		completeResp, err := client.CompletePairing(ctx, &pb.CompletePairingRequest{
			Creator:        "alice",
			ChallengeId:    pairingResp.ChallengeId,
			ComponentAAuth: "battery-auth-token",
			ComponentBAuth: "motor-auth-token",
			SessionContext: "race-session-001",
		})
		if err != nil {
			log.Printf("CompletePairing failed: %v", err)
		} else {
			fmt.Printf("Pairing completed: %s\n", completeResp.LctId)
			fmt.Printf("Transaction hash: %s\n", completeResp.Txhash)
			fmt.Printf("Split Key A: %s\n", completeResp.SplitKeyA)
			fmt.Printf("Split Key B: %s\n", completeResp.SplitKeyB)
		}
	}

	fmt.Println("\n✅ gRPC client test completed successfully!")
}
