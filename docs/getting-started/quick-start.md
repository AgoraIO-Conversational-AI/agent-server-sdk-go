---
sidebar_position: 3
title: Quick Start
description: Build and run your first Agora Conversational AI agent in Go with app credentials and presets.
---

# Quick Start

This guide uses the recommended onboarding path:

- `AppID`, `AppCertificate`, and `Area` on `AgoraClient`
- `Preset` for Agora-managed ASR, LLM, and TTS
- automatic ConvoAI REST auth and RTC join token generation
- no vendor API keys in application code

## Full example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    ctx := context.Background()
    idleTimeout := 120

    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          "your-app-id",
        AppCertificate: "your-app-certificate",
    })

    // Agent-level behavior lives here. Vendor selection comes from presets below.
    agent := agentkit.NewAgent(
        agentkit.WithName("support-assistant"),
        agentkit.WithInstructions("You are a concise support voice assistant."),
        agentkit.WithGreeting("Hello! How can I help you today?"),
        agentkit.WithMaxHistory(10),
    )

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:     "support-room-123",
        AgentUID:    "1",
        RemoteUIDs:  []string{"100"},
        IdleTimeout: &idleTimeout,
        Preset: []string{
            agentkit.AgentPresets.Asr.DeepgramNova3,
            agentkit.AgentPresets.Llm.OpenAIGpt5Mini,
            agentkit.AgentPresets.Tts.OpenAITts1,
        },
    })

    agentSessionID, err := session.Start(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Agent started:", agentSessionID)

    if err := session.Say(ctx, "Thanks for calling Agora support.", nil, nil); err != nil {
        log.Fatal(err)
    }
    if err := session.Stop(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## What this does

1. `AgoraClient` runs in app-credentials mode when you pass `AppID` and `AppCertificate` only.
2. `Agent` holds reusable behavior such as instructions, greeting, and history settings.
3. `Preset` tells Agora which managed ASR, LLM, and TTS vendors to run.
4. `session.Start(...)` lets the SDK generate the required auth tokens automatically.
5. `session.Start(...)` returns the unique agent session ID.

## When to use BYOK instead

Use presets when you want the fastest path to a working agent.

Use BYOK when you need to:

- supply your own vendor API keys
- use models outside the preset catalog
- point at custom vendor endpoints
- manage vendor-specific parameters directly

See [BYOK Guide](../guides/byok.md).

## Next steps

- [Authentication](./authentication.md)
- [BYOK Guide](../guides/byok.md)
- [MLLM Flow](../guides/mllm-flow.md)
- [Agent Reference](../reference/agent.md)
