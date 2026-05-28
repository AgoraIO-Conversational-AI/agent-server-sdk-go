---
sidebar_position: 4
title: Regional Routing
description: Configure regional routing for app-credentials AgentKit clients.
---

# Regional Routing

Use `Area` on `agentkit.NewAgoraClient` to route session requests to the desired Agora region while keeping the app-credentials auth path.

```go
package main

import (
    "github.com/AgoraIO/agora-agents-go/agentkit"
    "github.com/AgoraIO/agora-agents-go/option"
)

func main() {
    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          "your-app-id",
        AppCertificate: "your-app-certificate",
    })

    _ = client
}
```

## Available Areas

| Constant | Region | Domain Prefixes |
|---|---|---|
| `option.AreaUS` | United States | `api-us-west-1`, `api-us-east-1` |
| `option.AreaEU` | Europe | `api-eu-west-1`, `api-eu-central-1` |
| `option.AreaAP` | Asia-Pacific | `api-ap-southeast-1`, `api-ap-northeast-1` |
| `option.AreaCN` | Chinese Mainland | `api-cn-east-1`, `api-cn-north-1` |

Each area has two regional prefixes for failover. The SDK uses the selected area when constructing the Conversational AI API base URL.

## When To Override Lower-Level Routing

The generated `client` and `option` packages also expose lower-level routing primitives such as `option.WithBaseURL`, `option.WithPool`, and `core.NewPool`. Use those only when you are building direct generated-client integrations. For realtime agent sessions, prefer `AgoraClientOptions.Area` so auth, token generation, and regional routing stay together in AgentKit.
