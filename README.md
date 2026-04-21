# Agora Agent Server SDK for Go

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=https%3A%2F%2Fgithub.com%2FAgoraIO-Conversational-AI%2Fagent-server-sdk-go)

The Agora Conversational AI SDK provides convenient access to the Agora Conversational AI APIs, enabling you to build voice-powered AI agents with support for both cascading flows (ASR -> LLM -> TTS) and multimodal flows (MLLM) for real-time audio processing.

## Requirements

- Go 1.21+

## Installation

```sh
go mod init example.com/voice-agent
go get github.com/AgoraIO-Conversational-AI/agent-server-sdk-go
```

## Quick Start

Minimal builder-based example using supported preset-backed models with no vendor API keys:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          "your-app-id",
        AppCertificate: "your-app-certificate",
    })

    agent := agentkit.NewAgent(
        agentkit.WithInstructions("You are a concise voice assistant."),
        agentkit.WithGreeting("Hello! How can I help you today?"),
    ).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
        Model: "nova-3",
    })).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
        Model: "gpt-5-mini",
    })).WithTts(vendors.NewOpenAITTS(vendors.OpenAITTSOptions{
        Voice: "alloy",
    }))

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:    "support-room-123",
        AgentUID:   "1",
        RemoteUIDs: []string{"100"},
    })

    agentID, err := session.Start(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(agentID)
}
```

### Why no token or vendor key in the example?

The SDK-managed path is the recommended path. `AgoraClient` generates the required ConvoAI REST auth and RTC join tokens automatically, and AgentKit infers the matching supported presets from the vendor configs when you omit vendor API keys.

## BYOK

If you want to bring your own vendor credentials instead of using Agora-managed presets, use the BYOK guide:

- [BYOK Guide](./docs/guides/byok.md)

## Documentation

- [Overview](./docs/index.md)
- [Authentication](./docs/getting-started/authentication.md)
- [Quick Start](./docs/getting-started/quick-start.md)
- [BYOK Guide](./docs/guides/byok.md)
- [MLLM Flow](./docs/guides/mllm-flow.md)
- [Low-Level API](./docs/guides/low-level-api.md)

## Reference

- [SDK Reference](./reference.md)
- [Agora Conversational AI Docs](https://docs.agora.io/en/conversational-ai/overview)
