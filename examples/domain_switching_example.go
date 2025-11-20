package examples

import (
	"context"
	"fmt"
	"time"

	agora "github.com/fern-demo/agoraio-go-sdk/v505"
	"github.com/fern-demo/agoraio-go-sdk/v505/client"
	"github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func ExampleBasicUsage() {
	baseURL, err := agora.GetBaseURLForRegion(agora.RegionUS)
	if err != nil {
		panic(err)
	}

	c := client.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBasicAuth("username", "password"),
	)

	fmt.Printf("Client configured with base URL: %s\n", baseURL)
	_ = c
}

func ExampleAdvancedUsageWithDomainSwitching() {
	ds, err := agora.NewDomainSwitcher(agora.RegionEU)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ds.SelectBestDomain(ctx); err != nil {
		fmt.Printf("Warning: Could not select best domain: %v\n", err)
	}

	baseURL := ds.GetBaseURL()

	c := client.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBasicAuth("username", "password"),
	)

	fmt.Printf("Client configured with base URL: %s\n", baseURL)
	_ = c
}

func ExampleRegionFailover() {
	ds, err := agora.NewDomainSwitcher(agora.RegionAP)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		baseURL := ds.GetBaseURL()
		fmt.Printf("Attempt %d: Using base URL: %s\n", i+1, baseURL)

		c := client.NewClient(
			option.WithBaseURL(baseURL),
			option.WithBasicAuth("username", "password"),
		)

		_ = c

		ds.NextRegion()
	}
}

func ExampleAllRegions() {
	regions := []struct {
		name   string
		region agora.Region
	}{
		{"US", agora.RegionUS},
		{"EU", agora.RegionEU},
		{"AP", agora.RegionAP},
		{"CN", agora.RegionCN},
	}

	for _, r := range regions {
		baseURL, err := agora.GetBaseURLForRegion(r.region)
		if err != nil {
			fmt.Printf("Error getting base URL for %s: %v\n", r.name, err)
			continue
		}
		fmt.Printf("%s region base URL: %s\n", r.name, baseURL)
	}
}
