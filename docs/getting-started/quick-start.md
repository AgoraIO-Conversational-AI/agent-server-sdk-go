---
sidebar_position: 3
title: Quick Start
description: Build your first conversational AI agent with the Agora Go SDK in under 5 minutes.
---

# Quick Start

This guide walks through creating a conversational AI agent using the cascading ASR -> LLM -> TTS flow with the `agentkit` package.

## Full Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
)

func main() {
    // 1. Create the API client
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    // 2. Build an agent with functional options
    agent := agentkit.NewAgent(
        agentkit.WithName("my-assistant"),
        agentkit.WithInstructions("You are a helpful voice assistant. Keep responses concise."),
        agentkit.WithGreeting("Hello! How can I help you today?"),
        agentkit.WithFailureMessage("Sorry, something went wrong. Please try again."),
        agentkit.WithMaxHistory(10),
    ).WithLlm(
        vendors.NewOpenAI(vendors.OpenAIOptions{
            APIKey: "<your_openai_key>",
            Model:  "gpt-4o-mini",
        }),
    ).WithTts(
        vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
            Key:     "<your_elevenlabs_key>",
            ModelID: "eleven_turbo_v2_5",
            VoiceID: "<your_voice_id>",
        }),
    ).WithStt(
        vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
            APIKey: "<your_deepgram_key>",
        }),
    )

    // 3. Create a session
    session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
        Client:         c.Agents,
        Agent:          agent,
        AppID:          "<your_app_id>",
        AppCertificate: "<your_app_certificate>",
        Channel:        "my-channel",
        AgentUID:       "1001",
        RemoteUIDs:     []string{"1002"},
    })

    // 4. Register event handlers
    session.On("started", func(data interface{}) {
        fmt.Println("Agent started:", data)
    })
    session.On("stopped", func(data interface{}) {
        fmt.Println("Agent stopped:", data)
    })
    session.On("error", func(data interface{}) {
        log.Println("Agent error:", data)
    })

    // 5. Start the session
    ctx := context.Background()

    agentID, err := session.Start(ctx)
    if err != nil {
        log.Fatalf("Failed to start session: %v", err)
    }
    fmt.Println("Agent running with ID:", agentID)

    // 6. Interact with the agent
    err = session.Say(ctx, "Welcome to the demo!", nil, nil)
    if err != nil {
        log.Fatalf("Failed to say: %v", err)
    }

    // 7. Stop when done
    err = session.Stop(ctx)
    if err != nil {
        log.Fatalf("Failed to stop session: %v", err)
    }
    fmt.Println("Session stopped")
}
```

## Step-by-Step Breakdown

### 1. Create the Client

`client.NewClient` accepts variadic `option.RequestOption` arguments. At minimum, provide authentication:

```go
c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
)
```

### 2. Build an Agent

`agentkit.NewAgent` uses Go's functional options pattern. Pass configuration via `With*` option functions:

```go
agent := agentkit.NewAgent(
    agentkit.WithName("my-assistant"),
    agentkit.WithInstructions("You are a helpful voice assistant."),
)
```

Then chain vendor methods. Each `With*` vendor method returns a new `*Agent` (immutable cloning):

```go
agent = agent.
    WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{APIKey: "..."})).
    WithTts(vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{Key: "...", ModelID: "...", VoiceID: "..."})).
    WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{APIKey: "..."}))
```

### 3. Create and Start a Session

```go
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
```

### 4. Clean Up

Always stop the session when you're done:

```go
err = session.Stop(ctx)
if err != nil {
    log.Fatalf("Failed to stop: %v", err)
}
```

## Next Steps

- [Architecture](../concepts/architecture.md) — understand the two-layer design
- [Agent](../concepts/agent.md) — deep dive into the functional options pattern
- [Cascading Flow Guide](../guides/cascading-flow.md) — explore different vendor combinations
- [MLLM Flow Guide](../guides/mllm-flow.md) — use multimodal models instead of ASR -> LLM -> TTS
