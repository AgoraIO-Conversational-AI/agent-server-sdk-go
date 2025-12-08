package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fern-demo/agoraio-go-sdk/v505/domain"
	"github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func main() {
	fmt.Println("Testing domain.NewClient function...")
	fmt.Println()

	// Test 1: Create client with default options (US region)
	fmt.Println("Test 1: Creating client with default options (US region)...")
	client1, err := domain.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client with default options: %v", err)
	}
	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Current URL: %s\n", client1.GetCurrentURL())
	fmt.Println()

	// Test 2: Create client with EU region
	fmt.Println("Test 2: Creating client with EU region...")
	client2, err := domain.NewClient(
		domain.WithArea(domain.EU),
	)
	if err != nil {
		log.Fatalf("Failed to create client with EU region: %v", err)
	}
	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Current URL: %s\n", client2.GetCurrentURL())
	fmt.Println()

	// Test 3: Create client with AP region
	fmt.Println("Test 3: Creating client with AP region...")
	client3, err := domain.NewClient(
		domain.WithArea(domain.AP),
	)
	if err != nil {
		log.Fatalf("Failed to create client with AP region: %v", err)
	}
	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Current URL: %s\n", client3.GetCurrentURL())
	fmt.Println()

	// Test 4: Create client with CN region
	fmt.Println("Test 4: Creating client with CN region...")
	client4, err := domain.NewClient(
		domain.WithArea(domain.CN),
	)
	if err != nil {
		log.Fatalf("Failed to create client with CN region: %v", err)
	}
	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Current URL: %s\n", client4.GetCurrentURL())
	fmt.Println()

	// Test 5: Create client with custom API path
	fmt.Println("Test 5: Creating client with custom API path...")
	customPath := "/api/custom-path"
	client5, err := domain.NewClient(
		domain.WithArea(domain.US),
		domain.WithAPIPath(customPath),
	)
	if err != nil {
		log.Fatalf("Failed to create client with custom API path: %v", err)
	}
	fmt.Printf("✓ Client created successfully\n")
	fmt.Printf("  Current URL: %s\n", client5.GetCurrentURL())
	fmt.Println()

	// Test 6: Create client with request options (basic auth)
	fmt.Println("Test 6: Creating client with request options (basic auth)...")
	client6, err := domain.NewClient(
		domain.WithArea(domain.US),
		domain.WithRequestOptions(
			option.WithBasicAuth("test-user", "test-password"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create client with request options: %v", err)
	}
	fmt.Printf("✓ Client created successfully with basic auth\n")
	fmt.Printf("  Current URL: %s\n", client6.GetCurrentURL())
	fmt.Println()

	// Test 7: Verify pool access
	fmt.Println("Test 7: Verifying pool access...")
	pool := client1.GetPool()
	if pool == nil {
		log.Fatal("Failed to get pool from client")
	}
	fmt.Printf("✓ Pool retrieved successfully\n")
	fmt.Println()

	// Test 8: Test NextRegion functionality
	fmt.Println("Test 8: Testing NextRegion functionality...")
	urlBefore := client1.GetCurrentURL()
	client1.NextRegion()
	urlAfter := client1.GetCurrentURL()
	fmt.Printf("✓ NextRegion called successfully\n")
	fmt.Printf("  URL before: %s\n", urlBefore)
	fmt.Printf("  URL after:  %s\n", urlAfter)
	fmt.Println()

	// Test 9: Test SelectBestDomain (optional, may require network)
	fmt.Println("Test 9: Testing SelectBestDomain (DNS resolution)...")
	ctx := context.Background()
	err = client1.SelectBestDomain(ctx)
	if err != nil {
		fmt.Printf("⚠ SelectBestDomain returned error (may be expected): %v\n", err)
	} else {
		fmt.Printf("✓ SelectBestDomain completed successfully\n")
	}
	fmt.Println()

	// Test 10: Verify embedded client is accessible
	fmt.Println("Test 10: Verifying embedded client is accessible...")
	if client1.Client == nil {
		log.Fatal("Embedded client is nil")
	}
	fmt.Printf("✓ Embedded Fern client is accessible\n")
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("All tests completed successfully!")
	fmt.Println("========================================")
}
