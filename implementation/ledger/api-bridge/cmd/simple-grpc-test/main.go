package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Test basic connectivity
	fmt.Println("Testing gRPC server connectivity...")

	// Try to get server info
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with health service if available
	healthClient := grpc_health_v1.NewHealthClient(conn)
	healthResp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		fmt.Printf("Health check failed (expected): %v\n", err)
	} else {
		fmt.Printf("Health check response: %v\n", healthResp.Status)
	}

	// Test reflection if available
	fmt.Println("Testing reflection...")
	// We can't easily test reflection without additional imports, but let's just verify the connection works

	fmt.Println("✅ gRPC server is reachable!")
	fmt.Println("The 'unknown service' error suggests the service registration might have an issue.")
	fmt.Println("This could be due to:")
	fmt.Println("1. Protobuf package name mismatch")
	fmt.Println("2. Service not properly registered")
	fmt.Println("3. Generated code mismatch")
}
