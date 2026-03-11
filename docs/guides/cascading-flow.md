---
sidebar_position: 1
title: Cascading Flow
description: Build agents using the ASR -> LLM -> TTS cascading flow with different vendor combinations.
---

# Cascading Flow (ASR -> LLM -> TTS)

The cascading flow chains three vendor services: speech-to-text (STT/ASR) converts audio to text, an LLM generates a response, and text-to-speech (TTS) converts it back to audio.

## Example 1: OpenAI + ElevenLabs + Deepgram (Basic Auth)

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
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    agent := agentkit.NewAgent(
        agentkit.WithName("openai-assistant"),
        agentkit.WithInstructions("You are a helpful voice assistant. Keep responses under 3 sentences."),
        agentkit.WithGreeting("Hi there! What can I help you with?"),
        agentkit.WithMaxHistory(20),
    ).WithLlm(
        vendors.NewOpenAI(vendors.OpenAIOptions{
            APIKey: "<openai_key>",
            Model:  "gpt-4o-mini",
        }),
    ).WithTts(
        vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
            Key:     "<elevenlabs_key>",
            ModelID: "eleven_turbo_v2_5",
            VoiceID: "<voice_id>",
        }),
    ).WithStt(
        vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
            APIKey:   "<deepgram_key>",
            Model:    "nova-2",
            Language: "en-US",
        }),
    )

    session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
        Client:         c.Agents,
        Agent:          agent,
        AppID:          "<app_id>",
        AppCertificate: "<app_cert>",
        Channel:        "demo-channel",
        AgentUID:       "1001",
        RemoteUIDs:     []string{"1002"},
    })

    ctx := context.Background()

    agentID, err := session.Start(ctx)
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    fmt.Println("Running:", agentID)

    // ... interact with the agent ...

    err = session.Stop(ctx)
    if err != nil {
        log.Fatalf("Failed to stop: %v", err)
    }
}
```

## Example 2: Anthropic + Microsoft TTS + Microsoft STT (App Credentials)

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
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    agent := agentkit.NewAgent(
        agentkit.WithName("claude-assistant"),
        agentkit.WithInstructions("You are a customer service agent for Acme Corp."),
        agentkit.WithGreeting("Welcome to Acme Corp! How may I assist you?"),
        agentkit.WithFailureMessage("I apologize for the inconvenience. Please try again shortly."),
    ).WithLlm(
        vendors.NewAnthropic(vendors.AnthropicOptions{
            APIKey: "<anthropic_key>",
            Model:  "claude-3-5-sonnet-20241022",
        }),
    ).WithTts(
        vendors.NewMicrosoftTTS(vendors.MicrosoftTTSOptions{
            Key:       "<microsoft_key>",
            Region:    "eastus",
            VoiceName: "en-US-JennyNeural",
        }),
    ).WithStt(
        vendors.NewMicrosoftSTT(vendors.MicrosoftSTTOptions{
            Key:    "<microsoft_key>",
            Region: "eastus",
        }),
    )

    session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
        Client:         c.Agents,
        Agent:          agent,
        AppID:          "<app_id>",
        AppCertificate: "<app_cert>",
        Channel:        "support-channel",
        AgentUID:       "1001",
        RemoteUIDs:     []string{"1002"},
    })

    ctx := context.Background()

    agentID, err := session.Start(ctx)
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    fmt.Println("Running:", agentID)

    err = session.Stop(ctx)
    if err != nil {
        log.Fatalf("Failed to stop: %v", err)
    }
}
```

## Required Vendors

In cascading mode, **LLM and TTS are required**. STT is optional — if omitted, the platform uses a default ASR provider. `ToProperties` returns an error if LLM or TTS is missing:

```go
agent := agentkit.NewAgent(agentkit.WithName("no-tts"))
// No TTS or LLM configured

_, err := agent.ToProperties(agentkit.ToPropertiesOptions{...})
// err: "TTS configuration is required; use WithTts() to set it"
```

## Turn Detection

Add server-side voice activity detection to control when the agent starts processing:

```go
import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

agent := agentkit.NewAgent(
    agentkit.WithName("vad-agent"),
    agentkit.WithTurnDetectionConfig(&agentkit.TurnDetectionConfig{
        Type:              agentkit.TurnDetectionTypeServerVad.Ptr(), // deprecated; use Config.EndOfSpeech instead
        Threshold:         Agora.Float64(0.5),
        SilenceDurationMs: Agora.Int(500),
    }),
    // ... other options
)
```
