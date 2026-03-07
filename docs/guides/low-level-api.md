---
sidebar_position: 9
title: Low-Level API
description: Direct client.Agents.Start() usage without the builder pattern.
---

# Low-Level API

For direct control over the REST API, use `client.Agents.Start()` with raw request objects. See the [API Reference](../../reference.md) for full details.

## Direct Client Usage

```go
package main

import (
    "context"
    "log"

    Agora "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk"
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    req := &Agora.StartAgentsRequest{
        Appid: "your_app_id",
        Name:  "unique_name",
        Properties: &Agora.StartAgentsRequestProperties{
            Channel:       "channel_name",
            Token:         "your_token",
            AgentRtcUID:   "1001",
            RemoteRtcUIDs: []string{"1002"},
            IdleTimeout:   Agora.Int(120),
            Asr: &Agora.StartAgentsRequestPropertiesAsr{
                Language: Agora.String("en-US"),
                Vendor:   Agora.StartAgentsRequestPropertiesAsrVendorDeepgram.Ptr(),
                Params: map[string]interface{}{
                    "api_key": "your-deepgram-key",
                },
            },
            Tts: &Agora.Tts{
                Elevenlabs: &Agora.ElevenLabsTts{
                    Params: &Agora.ElevenLabsTtsParams{
                        Key:     "your-elevenlabs-key",
                        ModelID: "eleven_flash_v2_5",
                        VoiceID: "your-voice-id",
                    },
                },
            },
            Llm: &Agora.StartAgentsRequestPropertiesLlm{
                URL:    "https://api.openai.com/v1/chat/completions",
                APIKey: Agora.String("<your_llm_key>"),
                SystemMessages: []map[string]interface{}{
                    {"role": "system", "content": "You are a helpful chatbot."},
                },
                Params: map[string]interface{}{
                    "model": "gpt-4o-mini",
                },
                MaxHistory:       Agora.Int(32),
                GreetingMessage:  Agora.String("Hello, how can I assist you today?"),
                FailureMessage:   Agora.String("Please hold on a second."),
            },
        },
    }

    ctx := context.Background()
    resp, err := c.Agents.Start(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    _ = resp
}
```

## Using Pointer Helpers

The API uses pointer types for optional fields. Use `Agora.Int()`, `Agora.String()`, `Agora.Bool()` from the root package:

```go
IdleTimeout:   Agora.Int(120),
Language:      Agora.String("en-US"),
GreetingMessage: Agora.String("Hello!"),
```

## MLLM (Raw API)

For MLLM flow without the builder pattern, pass `mllm` and `turn_detection` directly. See the [MLLM Overview](https://docs.agora.io/en/conversational-ai/models/mllm/overview) for details.

```go
req := &Agora.StartAgentsRequest{
    Appid: "your_app_id",
    Name:  "mllm_agent",
    Properties: &Agora.StartAgentsRequestProperties{
        Channel:       "channel_name",
        Token:         "your_token",
        AgentRtcUID:   "1001",
        RemoteRtcUIDs: []string{"1002"},
        IdleTimeout:   Agora.Int(120),
        AdvancedFeatures: &Agora.StartAgentsRequestPropertiesAdvancedFeatures{
            EnableMllm: Agora.Bool(true),
        },
        Mllm: &Agora.StartAgentsRequestPropertiesMllm{
            URL:    "wss://api.openai.com/v1/realtime",
            APIKey: Agora.String("<your_openai_api_key>"),
            Vendor: Agora.StartAgentsRequestPropertiesMllmVendorOpenai.Ptr(),
            Params: map[string]interface{}{
                "model": "gpt-4o-realtime-preview",
                "voice": "alloy",
            },
            InputModalities:  []string{"audio"},
            OutputModalities: []string{"text", "audio"},
            GreetingMessage:  Agora.String("Hello! I'm ready to chat in real-time."),
        },
        TurnDetection: &Agora.StartAgentsRequestPropertiesTurnDetection{
            Type:               Agora.StartAgentsRequestPropertiesTurnDetectionTypeServerVad.Ptr(),
            Threshold:          Agora.Float64(0.5),
            SilenceDurationMs:  Agora.Int(500),
        },
        Tts: &Agora.Tts{
            Elevenlabs: &Agora.ElevenLabsTts{
                Params: &Agora.ElevenLabsTtsParams{
                    Key:     "your-elevenlabs-key",
                    ModelID: "eleven_flash_v2_5",
                    VoiceID: "your-voice-id",
                },
            },
        },
        Llm: &Agora.StartAgentsRequestPropertiesLlm{
            URL: "https://api.openai.com/v1/chat/completions",
        },
    },
}
```

For more on the agentkit-based MLLM flow, see [MLLM Flow](./mllm-flow.md).
