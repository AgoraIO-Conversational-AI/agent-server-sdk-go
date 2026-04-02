---
sidebar_position: 2
title: MLLM Flow
description: Use multimodal models (OpenAI Realtime, Vertex AI) for real-time audio processing.
---

# MLLM Flow (Multimodal)

The MLLM flow uses a single multimodal model to process audio input and generate audio output directly, bypassing the ASR -> LLM -> TTS chain. This provides lower latency and more natural conversational behavior.

## Enabling MLLM Mode

MLLM mode requires setting `EnableMllm: Agora.Bool(true)` in the advanced features:

```go
import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

agent := agentkit.NewAgent(
    agentkit.WithName("realtime-agent"),
    agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
        EnableMllm: Agora.Bool(true),
    }),
)
```

Note the use of `Agora.Bool(true)` — this is a pointer helper that returns `*bool`. The Agora API uses pointer types for optional fields, so you cannot pass a bare `true` literal.

## OpenAI Realtime Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    agent := agentkit.NewAgent(
        agentkit.WithName("openai-realtime"),
        agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
            EnableMllm: Agora.Bool(true),
        }),
    ).WithMllm(
        vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
            APIKey: "<openai_key>",
            Model:  "gpt-4o-realtime-preview",
            Params: map[string]interface{}{
                "voice": "alloy",
            },
        }),
    )

    session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
        Client:         c.Agents,
        Agent:          agent,
        AppID:          "<app_id>",
        AppCertificate: "<app_cert>",
        Channel:        "realtime-channel",
        AgentUID:       "1001",
        RemoteUIDs:     []string{"1002"},
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

## Vertex AI (Google Gemini Live) Example

```go
agent := agentkit.NewAgent(
    agentkit.WithName("gemini-live"),
    agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
        EnableMllm: Agora.Bool(true),
    }),
).WithMllm(
    vendors.NewVertexAI(vendors.VertexAIOptions{
        ProjectID:           "<gcp_project_id>",
        Location:            "us-central1",
        Model:               "gemini-2.0-flash-exp",
        ADCredentialsString: "<adc_json>",
        Instructions:        "You are a helpful assistant.",
        AdditionalParams: map[string]interface{}{
            "voice": "Puck",
        },
    }),
)
```

## MLLM with Turn Detection

Server-side VAD works with MLLM mode. The preferred approach uses the SOS/EOS (Start of Speech / End of Speech) model via `Config.StartOfSpeech` and `Config.EndOfSpeech` — see the [Agent Reference](../reference/agent.md) for full type definitions.

Legacy format:

```go
agent := agentkit.NewAgent(
    agentkit.WithName("realtime-vad"),
    agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
        EnableMllm: Agora.Bool(true),
    }),
    agentkit.WithTurnDetectionConfig(&agentkit.TurnDetectionConfig{
        Type:              agentkit.TurnDetectionTypeServerVad.Ptr(), // deprecated; use Config.EndOfSpeech instead
        Threshold:         Agora.Float64(0.5),
        SilenceDurationMs: Agora.Int(500),
    }),
).WithMllm(
    vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
        APIKey: "<openai_key>",
    }),
)
```

## Using the Raw Client

You can also use MLLM mode directly with the Fern-generated client without the agentkit package:

```go
c.Agents.Start(
    context.Background(),
    &Agora.StartAgentsRequest{
        Appid: "<app_id>",
        Name:  "mllm_agent",
        Properties: &Agora.StartAgentsRequestProperties{
            Channel:       "channel_name",
            Token:         "<token>",
            AgentRtcUID:   "1001",
            RemoteRtcUIDs: []string{"1002"},
            IdleTimeout:   Agora.Int(120),
            AdvancedFeatures: &Agora.StartAgentsRequestPropertiesAdvancedFeatures{
                EnableMllm: Agora.Bool(true),
            },
            Mllm: &Agora.StartAgentsRequestPropertiesMllm{
                URL:    Agora.String("wss://api.openai.com/v1/realtime"),
                APIKey: Agora.String("<openai_key>"),
                Vendor: Agora.StartAgentsRequestPropertiesMllmVendorOpenai,
                Params: map[string]any{
                    "model": "gpt-4o-realtime-preview",
                    "voice": "alloy",
                },
                InputModalities:  []string{"audio"},
                OutputModalities: []string{"text", "audio"},
            },
        },
    },
)
```

## Pointer Helper Functions

MLLM configuration makes heavy use of pointer helpers for optional fields:

| Helper | Type | Example |
|---|---|---|
| `Agora.Bool(true)` | `*bool` | `EnableMllm: Agora.Bool(true)` |
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
| `AdvancedFeatures.EnableMllm` | Not set or `false` | Must be `Agora.Bool(true)` |
