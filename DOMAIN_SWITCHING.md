# Domain Auto-Switching for Agora Go SDK

This document describes the domain auto-switching feature for the Agora Go SDK, which allows you to automatically select the best API endpoint based on your geographic region.

## Overview

The domain auto-switching feature provides a helper function that returns the appropriate BaseURL for the Agora API based on your deployment region. This ensures optimal performance and reliability by routing requests to the nearest available API gateway.

## Supported Regions

The following regions are supported:

- **RegionUS**: United States (West and East)
- **RegionEU**: Europe (West and Central)
- **RegionAP**: Asia-Pacific (Southeast and Northeast)
- **RegionCN**: Chinese Mainland (East and North)

## Basic Usage

### Simple Helper Function

The easiest way to use domain switching is with the `GetBaseURLForRegion` helper function:

```go
package main

import (
    agora "github.com/fern-demo/agoraio-go-sdk/v505"
    "github.com/fern-demo/agoraio-go-sdk/v505/client"
    "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func main() {
    // Get the base URL for the US region
    baseURL, err := agora.GetBaseURLForRegion(agora.RegionUS)
    if err != nil {
        panic(err)
    }

    // Create a client with the region-specific base URL
    c := client.NewClient(
        option.WithBaseURL(baseURL),
        option.WithBasicAuth("username", "password"),
    )

    // Use the client as normal
    // ...
}
```

## Advanced Usage

### Domain Switcher with Automatic Selection

For more advanced use cases, you can use the `DomainSwitcher` type to automatically select the best domain based on DNS resolution:

```go
package main

import (
    "context"
    "time"

    agora "github.com/fern-demo/agoraio-go-sdk/v505"
    "github.com/fern-demo/agoraio-go-sdk/v505/client"
    "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func main() {
    // Create a domain switcher for the EU region
    ds, err := agora.NewDomainSwitcher(agora.RegionEU)
    if err != nil {
        panic(err)
    }

    // Automatically select the best domain based on DNS resolution
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := ds.SelectBestDomain(ctx); err != nil {
        // Handle error - falls back to default domain
        // You can still use ds.GetBaseURL() which will return the default
    }

    // Get the selected base URL
    baseURL := ds.GetBaseURL()

    // Create a client with the selected base URL
    c := client.NewClient(
        option.WithBaseURL(baseURL),
        option.WithBasicAuth("username", "password"),
    )

    // Use the client as normal
    // ...
}
```

### Region Failover

You can manually switch between regions within the same geographic area:

```go
package main

import (
    agora "github.com/fern-demo/agoraio-go-sdk/v505"
    "github.com/fern-demo/agoraio-go-sdk/v505/client"
    "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func main() {
    ds, err := agora.NewDomainSwitcher(agora.RegionAP)
    if err != nil {
        panic(err)
    }

    // Try the first region
    baseURL := ds.GetBaseURL()
    c := client.NewClient(
        option.WithBaseURL(baseURL),
        option.WithBasicAuth("username", "password"),
    )

    // If the request fails, switch to the next region
    // err := c.SomeMethod()
    // if err != nil {
    //     ds.NextRegion()
    //     baseURL = ds.GetBaseURL()
    //     c = client.NewClient(
    //         option.WithBaseURL(baseURL),
    //         option.WithBasicAuth("username", "password"),
    //     )
    // }
}
```

## How It Works

The domain switcher works by:

1. **Region Selection**: You specify your deployment region (US, EU, AP, or CN)
2. **Domain Mapping**: Each region has multiple domain prefixes and suffixes
3. **DNS Resolution**: The `SelectBestDomain` method performs DNS lookups to find the fastest responding domain
4. **Automatic Caching**: Domain selection is cached for 30 seconds to avoid excessive DNS lookups
5. **Failover**: You can manually switch to alternative regions using `NextRegion()`

## Domain Mapping

Each region maps to specific domain prefixes and suffixes:

### US Region
- Prefixes: `api-us-west-1`, `api-us-east-1`
- Suffixes: `agora.io`, `sd-rtn.com`
- Default: `https://api-us-west-1.agora.io`

### EU Region
- Prefixes: `api-eu-west-1`, `api-eu-central-1`
- Suffixes: `agora.io`, `sd-rtn.com`
- Default: `https://api-eu-west-1.agora.io`

### AP Region
- Prefixes: `api-ap-southeast-1`, `api-ap-northeast-1`
- Suffixes: `agora.io`, `sd-rtn.com`
- Default: `https://api-ap-southeast-1.agora.io`

### CN Region
- Prefixes: `api-cn-east-1`, `api-cn-north-1`
- Suffixes: `sd-rtn.com`, `agora.io`
- Default: `https://api-cn-east-1.sd-rtn.com`

## Thread Safety

The `DomainSwitcher` type is thread-safe and can be used concurrently from multiple goroutines.

## Best Practices

1. **Choose the Right Region**: Select the region closest to your deployment for best performance
2. **Use Simple Helper for Most Cases**: The `GetBaseURLForRegion` function is sufficient for most use cases
3. **Use Advanced Features for Resilience**: Use `DomainSwitcher` with `SelectBestDomain` for automatic failover
4. **Handle Errors Gracefully**: Always check for errors when selecting domains
5. **Reuse Domain Switcher**: Create one `DomainSwitcher` instance and reuse it across your application

## Examples

See the `examples/domain_switching_example.go` file for complete working examples.
