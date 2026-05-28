---
sidebar_position: 2
title: MLLM Flow
description: Use multimodal models (OpenAI Realtime, Gemini Live) for real-time audio processing.
---

# MLLM Flow (Multimodal)

The MLLM flow uses a single multimodal model to process audio input and generate audio output directly, bypassing the ASR -> LLM -> TTS chain. This provides lower latency and more natural conversational behavior.

## Enabling MLLM Mode

Call `WithMllm(vendor)` to enable MLLM mode. The builder sets `mllm.enable = true` automatically.

```go
agent := agentkit.NewAgent(
    agentkit.WithName("realtime-agent"),
).WithMllm(vendor)
```

## OpenAI Realtime Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    Agora "github.com/AgoraIO/agora-agents-go"
    "github.com/AgoraIO/agora-agents-go/agentkit"
    "github.com/AgoraIO/agora-agents-go/agentkit/vendors"
    "github.com/AgoraIO/agora-agents-go/option"
)

func main() {
    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          "<app_id>",
        AppCertificate: "<app_cert>",
    })

    agent := agentkit.NewAgent(
        agentkit.WithName("openai-realtime"),
    ).WithMllm(
        vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
            APIKey: "<openai_key>",
            Model:  "gpt-4o-realtime-preview",
            Params: map[string]interface{}{
                "voice": "alloy",
            },
        }),
    )

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:    "realtime-channel",
        AgentUID:   "1001",
        RemoteUIDs: []string{"1002"},
    })

    ctx := context.Background()

    agentID, err := session.Start(ctx)
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    fmt.Println("Realtime agent running:", agentID)

    err = session.Stop(ctx)
    if err != nil {
        log.Fatalf("Failed to stop: %v", err)
    }
}
```

## Gemini Live Example

```go
agent := agentkit.NewAgent(
    agentkit.WithName("gemini-live"),
).WithMllm(
    vendors.NewGeminiLive(vendors.GeminiLiveOptions{
        APIKey:       "<google_ai_api_key>",
        Model:        "gemini-live-2.5-flash",
        Instructions: "You are a helpful assistant.",
        Voice:        "Puck",
    }),
)
```

## MLLM with Turn Detection

Configure MLLM turn detection on the MLLM vendor with `TurnDetection`. When set, `mllm.turn_detection` overrides the top-level `turn_detection` object.

Example:

```go
agent := agentkit.NewAgent(
    agentkit.WithName("realtime-vad"),
).WithMllm(
    vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
        APIKey: "<openai_key>",
        TurnDetection: &Agora.StartAgentsRequestPropertiesMllmTurnDetection{
            Mode: Agora.StartAgentsRequestPropertiesMllmTurnDetectionModeServerVad.Ptr(),
            ServerVadConfig: &Agora.StartAgentsRequestPropertiesMllmTurnDetectionServerVadConfig{
                IdleTimeoutMs: Agora.Int(5000),
            },
        },
    }),
)
```

## Low-Level Usage

Use AgentKit for MLLM session starts. It keeps the app-credentials flow intact and generates the required REST and RTC tokens from `AppID` and `AppCertificate`. See [Low-Level API](./low-level-api.md) for generated-client escape hatches that are outside the realtime session lifecycle.

## Pointer Helper Functions

MLLM configuration makes heavy use of pointer helpers for optional fields:

| Helper | Type | Example |
|---|---|---|
| `Agora.Bool(true)` | `*bool` | `Enable: Agora.Bool(true)` |
| `Agora.String("...")` | `*string` | `APIKey: Agora.String("<key>")` |
| `Agora.Int(120)` | `*int` | `IdleTimeout: Agora.Int(120)` |
| `Agora.Float64(0.5)` | `*float64` | `Threshold: Agora.Float64(0.5)` |

These exist because Go does not allow taking the address of a literal value (`&true` is invalid). The helpers return pointers to the given values.

## Key Differences from Cascading Flow

| Aspect | Cascading | MLLM |
|---|---|---|
| Vendors required | LLM + TTS (STT optional) | MLLM only |
| Audio processing | Three-step chain | Single model, end-to-end |
| Latency | Higher (3 network hops) | Lower (1 network hop) |
| `mllm.enable` | Not set or `false` | Must be `Agora.Bool(true)` |
