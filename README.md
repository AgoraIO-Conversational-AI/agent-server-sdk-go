# Agora Agent Server SDK for Go

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=https%3A%2F%2Fgithub.com%2FAgoraIO-Conversational-AI%2Fagent-server-sdk-go)

The Agora Agent Server SDK for Go lets you build real-time voice agents on Agora Conversational AI with a high-level `Agent` / `AgentSession` API and a generated low-level REST client.

## Requirements

- Go 1.21+

## Installation

```sh
go mod init example.com/voice-agent
go get github.com/AgoraIO-Conversational-AI/agent-server-sdk-go
```

## Quick Start

The recommended onboarding path is a server-side builder flow: define the agent once, configure preset-backed providers in the builder, and let AgentKit infer the reseller `preset` values when the session starts.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

const (
    agentPrompt = "You are a concise, technically credible voice assistant. Keep replies short unless the user asks for detail."
    greeting    = "Hi there! I am your Agora voice assistant. How can I help?"
)

func stringPtr(v string) *string { return &v }
func intPtr(v int) *int { return &v }
func float64Ptr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool { return &v }

func requireEnv(name string) (string, error) {
    value := os.Getenv(name)
    if value == "" {
        return "", fmt.Errorf("missing required environment variable: %s", name)
    }
    return value, nil
}

func startConversation(ctx context.Context) (string, error) {
    appID, err := requireEnv("AGORA_APP_ID")
    if err != nil {
        return "", err
    }
    appCertificate, err := requireEnv("AGORA_APP_CERTIFICATE")
    if err != nil {
        return "", err
    }
    expiresIn, err := agentkit.ExpiresInHours(1)
    if err != nil {
        return "", err
    }

    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          appID,
        AppCertificate: appCertificate,
    })

    agent := agentkit.NewAgent(
        agentkit.WithName(fmt.Sprintf("conversation-%d", time.Now().UnixMilli())),
        agentkit.WithInstructions(agentPrompt),
        agentkit.WithGreeting(greeting),
        agentkit.WithFailureMessage("Please wait a moment."),
        agentkit.WithMaxHistory(50),
        agentkit.WithTurnDetectionConfig(&agentkit.TurnDetectionConfig{
            Config: &agentkit.TurnDetectionNestedConfig{
                SpeechThreshold: float64Ptr(0.5),
                StartOfSpeech: &agentkit.StartOfSpeechConfig{
                    Mode: agentkit.StartOfSpeechMode("vad"),
                    VadConfig: &agentkit.StartOfSpeechVadConfig{
                        InterruptDurationMs: intPtr(160),
                        PrefixPaddingMs:     intPtr(300),
                    },
                },
                EndOfSpeech: &agentkit.EndOfSpeechConfig{
                    Mode: agentkit.EndOfSpeechMode("vad"),
                    VadConfig: &agentkit.EndOfSpeechVadConfig{
                        SilenceDurationMs: intPtr(480),
                    },
                },
            },
        }),
        agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
            EnableRtm:   boolPtr(true),
            EnableTools: boolPtr(true),
        }),
        agentkit.WithParameters(&agentkit.SessionParams{
            DataChannel:        &agentkit.DataChannelRtm,
            EnableErrorMessage: boolPtr(true),
        }),
    ).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
        Model: "nova-3",
        Language: "en",
    })).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
        Model:           "gpt-4o-mini",
        GreetingMessage: greeting,
        FailureMessage:  "Please wait a moment.",
        MaxHistory:      intPtr(15),
        Params: map[string]interface{}{
            "max_tokens": 1024,
            "temperature": 0.7,
            "top_p": 0.95,
        },
    })).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
        Model:   "speech_2_6_turbo",
        VoiceID: "English_captivating_female1",
    }))

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:     fmt.Sprintf("demo-channel-%d", time.Now().UnixMilli()),
        AgentUID:    "123456",
        RemoteUIDs:  []string{"*"},
        IdleTimeout: intPtr(30),
        ExpiresIn:   expiresIn,
        Debug:       false,
    })

    return session.Start(ctx)
}

func main() {
    agentID, err := startConversation(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Println(agentID)
}
```

### Why no token or vendor key in the example?

`AgoraClient` generates the required ConvoAI REST auth and RTC join tokens automatically when you provide `AppID` and `AppCertificate`. AgentKit then inspects the builder-provided vendor configs and infers the matching supported `preset` values for reseller-backed models, so you do not pass vendor API keys in this flow.

### BYOK version of the same builder flow

Use the same `Agent` builder shape, but provide credentials explicitly when you want vendor-managed billing and routing instead of Agora-managed presets.

```go
agent := agentkit.NewAgent(
    agentkit.WithInstructions(agentPrompt),
    agentkit.WithGreeting(greeting),
).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
    APIKey:   os.Getenv("DEEPGRAM_API_KEY"),
    Model:    "nova-3",
    Language: "en",
})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
    APIKey:      os.Getenv("OPENAI_API_KEY"),
    Model:       "gpt-4o-mini",
    MaxTokens:   intPtr(1024),
    Temperature: float64Ptr(0.7),
    TopP:        float64Ptr(0.95),
})).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
    Key:     os.Getenv("MINIMAX_API_KEY"),
    GroupID: os.Getenv("MINIMAX_GROUP_ID"),
    Model:   "speech_2_6_turbo",
    VoiceID: "English_captivating_female1",
    URL:     "wss://api-uw.minimax.io/ws/v1/t2a_v2",
}))
```

## BYOK

If you want to bring your own vendor credentials instead of using Agora-managed presets, use the BYOK guide:

- [BYOK Guide](./docs/guides/byok.md)

## MLLM (Realtime / Multimodal)

Use `WithMllm()` for OpenAI Realtime or Gemini Live. No STT, LLM, or TTS vendor is needed when MLLM mode is enabled.

```go
agent := agentkit.NewAgent(
    agentkit.WithName("realtime-assistant"),
).WithMllm(vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
    APIKey:          os.Getenv("OPENAI_API_KEY"),
    Model:           "gpt-4o-realtime-preview",
    GreetingMessage: "Hello! Ready to chat.",
}))
```

See the [MLLM Flow guide](./docs/guides/mllm-flow.md) for full examples with Gemini Live and Vertex AI.

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
