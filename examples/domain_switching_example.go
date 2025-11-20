package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/fern-demo/agoraio-go-sdk/v505/domain"
	"github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func ExampleSimpleUsage() {
	c, err := domain.NewClientWithRegion(
		domain.RegionUS,
		option.WithBasicAuth("username", "password"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client created with US region\n")
	_ = c
}

func ExampleWithBaseURL() {
	baseURL, err := domain.GetBaseURLForRegion(domain.RegionUS)
	if err != nil {
		panic(err)
	}

	c, err := domain.NewClientWithRegion(
		domain.RegionEU,
		option.WithBaseURL(baseURL),
		option.WithBasicAuth("username", "password"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client configured with explicit base URL: %s\n", baseURL)
	_ = c
}

func ExampleAdvancedUsageWithDomainSwitching() {
	ds, err := domain.NewDomainSwitcher(domain.RegionEU)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ds.SelectBestDomain(ctx); err != nil {
		fmt.Printf("Warning: Could not select best domain: %v\n", err)
	}

	baseURL := ds.GetBaseURL()

	c, err := domain.NewClientWithRegion(
		domain.RegionEU,
		option.WithBaseURL(baseURL),
		option.WithBasicAuth("username", "password"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client configured with base URL: %s\n", baseURL)
	_ = c
}

func ExampleRegionFailover() {
	ds, err := domain.NewDomainSwitcher(domain.RegionAP)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		baseURL := ds.GetBaseURL()
		fmt.Printf("Attempt %d: Using base URL: %s\n", i+1, baseURL)

		c, err := domain.NewClientWithRegion(
			domain.RegionAP,
			option.WithBaseURL(baseURL),
			option.WithBasicAuth("username", "password"),
		)
		if err != nil {
			panic(err)
		}

		_ = c

		ds.NextRegion()
	}
}

func ExampleAllRegions() {
	regions := []struct {
		name   string
		region domain.Region
	}{
		{"US", domain.RegionUS},
		{"EU", domain.RegionEU},
		{"AP", domain.RegionAP},
		{"CN", domain.RegionCN},
	}

	for _, r := range regions {
		baseURL, err := domain.GetBaseURLForRegion(r.region)
		if err != nil {
			fmt.Printf("Error getting base URL for %s: %v\n", r.name, err)
			continue
		}
		fmt.Printf("%s region base URL: %s\n", r.name, baseURL)
	}
}
