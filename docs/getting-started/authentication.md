---
sidebar_position: 2
title: Authentication
description: Configure authentication for the Agora Conversational AI Go SDK — Basic Auth, Token Auth, and App Credentials.
---

# Authentication

The SDK supports three authentication modes. Choose the one that fits your deployment model.

## 1. Basic Auth (Customer ID + Customer Secret)

Use `option.WithBasicAuth` to authenticate with your Agora Customer ID and Customer Secret from the Agora Console RESTful API page.

```go
package main

import (
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )
    // Use c.Agents, c.Telephony, c.PhoneNumbers...
    _ = c
}
```

## 2. Token Auth

Use `option.WithToken` to authenticate with a bearer token.

```go
c := client.NewClient(
    option.WithToken("<your_api_token>"),
)
```

## 3. App Credentials (Agentkit Layer)

When using the `agentkit` package, pass `AppID` and `AppCertificate` to `AgentSessionOptions`. The agentkit package automatically generates a combined RTC+RTM token for the Agora channel using HMAC-SHA256.

```go
package main

import (
    "context"
    "log"

    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit"
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    // agent := agentkit.NewAgent(...)

    session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
        Client:         c.Agents,
        Agent:          agent,
        AppID:          "<your_app_id>",
        AppCertificate: "<your_app_certificate>",
        Channel:        "my-channel",
        AgentUID:       "1001",
        RemoteUIDs:     []string{"1002"},
    })

    agentID, err := session.Start(context.Background())
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    _ = agentID
}
```

You can also pass a pre-generated token instead of app credentials:

```go
session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
    Client:     c.Agents,
    Agent:      agent,
    AppID:      "<your_app_id>",
    Token:      "<pre_generated_rtc_rtm_token>",
    Channel:    "my-channel",
    AgentUID:   "1001",
    RemoteUIDs: []string{"1002"},
})
```

## Comparison

| Mode | When to Use | Credential Source | Layer |
|---|---|---|---|
| Basic Auth | Direct API calls | Agora Console > RESTful API | `client` + `option` |
| Token Auth | Direct API calls with bearer token | Your token provider | `client` + `option` |
| App Credentials | Agentkit layer, auto token generation | Agora Console > Project Management | `agentkit` |

## Token Generation Details

When using app credentials, the agentkit package calls `agentkit.GenerateConvoAIToken` internally (RTC + RTM combined token) with these defaults:

- **Role:** Publisher (`RolePublisher = 1`)
- **Expiry:** `86400` seconds (24 hours, Agora maximum), configurable via `AgentSessionOptions.ExpiresIn`
- **Algorithm:** HMAC-SHA256

You can also call `GenerateConvoAIToken` directly if you need a token outside of a session:

```go
token, err := agentkit.GenerateConvoAIToken(agentkit.GenerateConvoAITokenOptions{
    AppID:          "<your_app_id>",
    AppCertificate: "<your_app_certificate>",
    ChannelName:    "my-channel",
    Account:        "1001",
    TokenExpire:    86400,
})
if err != nil {
    log.Fatalf("Token generation failed: %v", err)
}
```

## Token Expiry

To customise the token lifetime, set `ExpiresIn` on `AgentSessionOptions`. Use the `ExpiresInHours` or `ExpiresInMinutes` helpers to avoid raw second values:

```go
expiresIn, err := agentkit.ExpiresInHours(12) // 12-hour token
if err != nil {
    log.Fatalf("Invalid expiry: %v", err)
}

session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
    Client:         c.Agents,
    Agent:          agent,
    AppID:          "<your_app_id>",
    AppCertificate: "<your_app_certificate>",
    Channel:        "my-channel",
    AgentUID:       "1001",
    RemoteUIDs:     []string{"1002"},
    ExpiresIn:      expiresIn,
})
```

`ExpiresInHours` and `ExpiresInMinutes` return an error if the value is ≤ 0, and log a warning + cap at 86400 if it exceeds 24 hours. Valid range: **1–86400 seconds**.
