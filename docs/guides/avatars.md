---
sidebar_position: 3
title: Avatars
description: Add visual avatars to your agent — HeyGen and Akool with sample rate requirements.
---

# Avatars

Avatars provide a visual representation of your AI agent. The SDK supports two avatar vendors, each with a strict TTS sample rate requirement.

## Avatar Vendors

| Vendor | Constructor | Required Sample Rate | Required Fields |
|---|---|---|---|
| HeyGen | `vendors.NewHeyGenAvatar` | 24kHz (`SampleRate24kHz`) | `APIKey`, `Quality` (low/medium/high), `AgoraUID` |
| Akool | `vendors.NewAkoolAvatar` | 16kHz (`SampleRate16kHz`) | `APIKey` |

## HeyGen Avatar Example

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
        AppID:          "<app_id>",
        AppCertificate: "<app_cert>",
    })

    sr := vendors.SampleRate24kHz

    agent := agentkit.NewAgent(
        agentkit.WithName("avatar-agent"),
        agentkit.WithInstructions("You are a friendly virtual assistant with a visual avatar."),
        agentkit.WithGreeting("Hello! I can see you and you can see me!"),
    ).WithLlm(
        vendors.NewOpenAI(vendors.OpenAIOptions{
            APIKey: "<openai_key>",
        }),
    ).WithTts(
        // TTS sample rate MUST match the avatar's required rate (24kHz for HeyGen)
        vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
            Key:        "<elevenlabs_key>",
            ModelID:    "eleven_turbo_v2_5",
            VoiceID:    "<voice_id>",
            SampleRate: &sr,
        }),
    ).WithStt(
        vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
            APIKey: "<deepgram_key>",
        }),
    ).WithAvatar(
        vendors.NewHeyGenAvatar(vendors.HeyGenAvatarOptions{
            APIKey:   "<heygen_key>",
            Quality:  "high",
            AgoraUID: "2001",
        }),
    )

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:    "avatar-channel",
        AgentUID:   "1001",
        RemoteUIDs: []string{"1002"},
    })

    ctx := context.Background()

    agentID, err := session.Start(ctx)
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    fmt.Println("Avatar agent running:", agentID)

    err = session.Stop(ctx)
    if err != nil {
        log.Fatalf("Failed to stop: %v", err)
    }
}
```

## Akool Avatar Example

```go
sr := vendors.SampleRate16kHz

agent := agentkit.NewAgent(
    agentkit.WithName("akool-avatar"),
    agentkit.WithInstructions("You are a virtual presenter."),
).WithLlm(
    vendors.NewOpenAI(vendors.OpenAIOptions{
        APIKey: "<openai_key>",
    }),
).WithTts(
    // Akool requires 16kHz
    vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
        Key:        "<elevenlabs_key>",
        ModelID:    "eleven_turbo_v2_5",
        VoiceID:    "<voice_id>",
        SampleRate: &sr,
    }),
).WithStt(
    vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
        APIKey: "<deepgram_key>",
    }),
).WithAvatar(
    vendors.NewAkoolAvatar(vendors.AkoolAvatarOptions{
        APIKey: "<akool_key>",
    }),
)
```

## Sample Rate Validation and Panic Behavior

The Go SDK enforces the TTS/avatar sample rate constraint at runtime using `panic()`. This differs from TypeScript (compile-time phantom types) and Python (`ValueError`).

### When does the panic occur?

`WithAvatar()` checks if the agent already has a TTS configured with a mismatched sample rate. If so, it panics immediately:

```
Avatar requires TTS sample rate of 24000 Hz, but TTS is configured with 16000 Hz. Please update your TTS sample_rate to 24000.
```

### Why panic instead of returning an error?

In Go, `panic` is idiomatic for programmer errors — configuration mistakes that should be caught during development, not handled at runtime. A sample rate mismatch is a static configuration error, not a transient failure.

### How to prevent the panic

Always configure TTS with the correct sample rate **before** calling `WithAvatar()`:

```go
// Correct: TTS sample rate matches avatar requirement
sr := vendors.SampleRate24kHz
agent := agentkit.NewAgent(...).
    WithTts(vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
        Key:        "<key>",
        ModelID:    "<model>",
        VoiceID:    "<voice>",
        SampleRate: &sr,           // 24kHz for HeyGen
    })).
    WithAvatar(vendors.NewHeyGenAvatar(vendors.HeyGenAvatarOptions{
        APIKey:   "<key>",
        Quality:  "high",
        AgoraUID: "2001",
    }))
```

```go
// Wrong: This panics — TTS is 16kHz but HeyGen requires 24kHz
sr := vendors.SampleRate16kHz
agent := agentkit.NewAgent(...).
    WithTts(vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
        Key:        "<key>",
        ModelID:    "<model>",
        VoiceID:    "<voice>",
        SampleRate: &sr,           // 16kHz — mismatch!
    })).
    WithAvatar(vendors.NewHeyGenAvatar(vendors.HeyGenAvatarOptions{
        APIKey:   "<key>",
        Quality:  "high",
        AgoraUID: "2001",
    }))
// panic: Avatar requires TTS sample rate of 24000 Hz, but TTS is configured with 16000 Hz. Please update your TTS sample_rate to 24000.
```

### Additional validation at Start

`AgentSession.Start()` also validates the sample rate match before making the API call. If the mismatch was introduced after `WithAvatar()` (e.g., by cloning the agent with a different TTS), `Start()` returns an error instead of panicking.

## HeyGenAvatarOptions Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | HeyGen API key |
| `Quality` | `string` | Yes | `"low"`, `"medium"`, or `"high"` |
| `AgoraUID` | `string` | Yes | UID for the avatar's video stream |
| `AvatarName` | `string` | No | Specific avatar model name |
| `VoiceID` | `string` | No | Override voice for the avatar |
| `Language` | `string` | No | Language code |
| `Version` | `string` | No | API version |

## AkoolAvatarOptions Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | Akool API key |
