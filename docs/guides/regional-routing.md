---
sidebar_position: 4
title: Regional Routing
description: Configure regional routing with option.WithArea, option.WithBaseURL, and Pool-based DNS selection.
---

# Regional Routing

The SDK supports routing API requests to regional endpoints for lower latency and data residency compliance.

## Quick Start: `option.WithArea`

The simplest way to enable regional routing is `option.WithArea`, which automatically creates a `Pool` that selects the best domain via DNS resolution:

```go
package main

import (
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
        option.WithArea(option.AreaUS),
    )
    _ = c
}
```

### Available Areas

| Constant | Region | Domain Prefixes |
|---|---|---|
| `option.AreaUS` | United States | `api-us-west-1`, `api-us-east-1` |
| `option.AreaEU` | Europe | `api-eu-west-1`, `api-eu-central-1` |
| `option.AreaAP` | Asia-Pacific | `api-ap-southeast-1`, `api-ap-northeast-1` |
| `option.AreaCN` | Chinese Mainland | `api-cn-east-1`, `api-cn-north-1` |

Each area has two regional prefixes for failover, and two domain suffixes (`agora.io` for overseas, `sd-rtn.com` for Chinese mainland).

## Manual Base URL Override

For direct control, use `option.WithBaseURL`:

```go
c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    option.WithBaseURL("https://api-us-west-1.agora.io"),
)
```

The default API endpoint is available as a constant:

```go
import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

// Agora.Environments.Default = "https://api.agora.io/api/conversational-ai-agent"
c := client.NewClient(
    option.WithBaseURL(Agora.Environments.Default),
)
```

## Pool-Based DNS Selection

For advanced use cases, manage the `Pool` lifecycle yourself using `option.WithPool`:

```go
package main

import (
    "context"
    "log"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    pool, err := core.NewPool(core.AreaUS)
    if err != nil {
        log.Fatalf("Failed to create pool: %v", err)
    }

    // Perform DNS-based domain selection
    ctx := context.Background()
    err = pool.SelectBestDomain(ctx)
    if err != nil {
        log.Printf("DNS selection failed, using default: %v", err)
    }

    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
        option.WithPool(pool),
    )
    _ = c
}
```

### Pool Methods

| Method | Description |
|---|---|
| `SelectBestDomain(ctx context.Context) error` | DNS-based selection of the fastest domain suffix. Caches result for 30 seconds. |
| `NextRegion()` | Cycle to the next regional prefix (e.g., `us-west-1` -> `us-east-1`) |
| `GetCurrentURL() string` | Returns the current URL (e.g., `https://api-us-west-1.agora.io`) |

## Region Failover with Context Timeouts

Use Go's `context.WithTimeout` to implement manual region failover — this is idiomatic Go for request-scoped timeouts:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func startWithFailover(pool *core.Pool, req *Agora.StartAgentsRequest) error {
    regions := []string{"primary", "secondary"}
    for _, region := range regions {
        c := client.NewClient(
            option.WithBasicAuth("<customer_id>", "<customer_secret>"),
            option.WithPool(pool),
        )

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        _, err := c.Agents.Start(ctx, req)
        cancel()

        if err == nil {
            fmt.Printf("Started on %s region\n", region)
            return nil
        }
        log.Printf("Region %s failed: %v, trying next...", region, err)
        pool.NextRegion()
    }
    return fmt.Errorf("all regions failed")
}

func main() {
    pool, err := core.NewPool(core.AreaUS)
    if err != nil {
        log.Fatalf("Failed to create pool: %v", err)
    }

    req := &Agora.StartAgentsRequest{
        Appid: "<app_id>",
        Name:  "failover-agent",
        Properties: &Agora.StartAgentsRequestProperties{
            Channel:     "my-channel",
            Token:       "<token>",
            AgentRtcUID: "1001",
        },
    }

    err = startWithFailover(pool, req)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
}
```

## How DNS Selection Works

When `SelectBestDomain` is called:

1. The pool resolves all domain suffixes for the current region prefix concurrently using `net.LookupHost`
2. The first domain that resolves successfully is selected
3. The result is cached for 30 seconds (`updateDuration`)
4. Subsequent calls within the cache window return immediately

This ensures the client uses the fastest-responding domain without repeated DNS lookups.
