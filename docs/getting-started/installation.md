---
sidebar_position: 1
title: Installation
description: Install the Agora Conversational AI Go SDK and configure your project.
---

# Installation

## Prerequisites

- **Go 1.21** or later

## Install

```sh
go get github.com/AgoraIO-Conversational-AI/agent-server-sdk-go
```

## Import Paths

The SDK is organized into several packages. Import the ones you need:

```go
// Root package — type definitions, pointer helpers (Agora.String(), Agora.Bool(), etc.), environments
Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

// Fern-generated client — low-level API access
"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"

// Request options — authentication, base URL, retries, HTTP client
"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"

// Wrapper layer — Agent builder with functional options, session lifecycle
"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"

// Vendor constructors — LLM, TTS, STT, MLLM, and Avatar vendors
"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
```

## Verify Installation

```go
package main

import (
    "fmt"

    Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
)

func main() {
    fmt.Println("Default API endpoint:", Agora.Environments.Default)
}
```

```sh
go run main.go
# Output: Default API endpoint: https://api.agora.io/api/conversational-ai-agent
```

## Next Steps

- [Authentication](./authentication.md) — configure your credentials
- [Quick Start](./quick-start.md) — build your first conversational agent
