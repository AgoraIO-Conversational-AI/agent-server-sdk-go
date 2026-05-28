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

## Why this is the default

- The SDK handles ConvoAI REST auth and RTC join token generation for you.
- Your onboarding code stays focused on agent behavior instead of auth plumbing.
- Your quick start code stays vendor-key free when you use supported Agora-managed models.

## Legacy auth modes

The generated low-level client still contains legacy auth hooks for prebuilt REST tokens and HTTP Basic Auth. Do not use those for new session integrations. Use app credentials so AgentKit can mint short-lived ConvoAI REST auth and RTC join tokens for each session.

## Inspecting the resolved auth mode

```go
import "fmt"

fmt.Println(client.AuthMode) // "app-credentials"
```
