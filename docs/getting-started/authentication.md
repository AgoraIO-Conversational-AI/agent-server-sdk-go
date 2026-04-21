---
sidebar_position: 2
title: Authentication
description: Configure the Go SDK with the recommended app-credentials flow and understand the supported auth modes.
---

# Authentication

The recommended production path is app credentials mode.

Create `AgoraClient` with `AppID` and `AppCertificate`, then let `AgentSession` generate the ConvoAI REST auth token and the RTC join token automatically.

## Recommended: app credentials

```go
package main

import (
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
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

## Why this is the default

- The SDK handles ConvoAI REST auth and RTC join token generation for you.
- Your onboarding code stays focused on agent behavior instead of auth plumbing.
- Your quick start code stays vendor-key free when you use presets.

## Other supported modes

The SDK also supports pre-built token auth and Basic Auth, but they are intentionally not the default onboarding path.

- Pre-built token auth exists for advanced or testing-specific cases.
- Basic Auth is supported for legacy integrations and account-level workflows.

## Inspecting the resolved auth mode

```go
import "fmt"

fmt.Println(client.AuthMode) // "app-credentials"
```
