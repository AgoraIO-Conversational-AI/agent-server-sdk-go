---
sidebar_position: 4
title: Vendors
description: Vendor catalog — LLM, TTS, STT, MLLM, and Avatar constructors with configuration structs.
---

# Vendors

The `agentkit/vendors` package provides constructor functions for all supported third-party vendors. Each vendor implements one of five interfaces: `LLM`, `TTS`, `STT`, `MLLM`, or `Avatar`.

## Vendor Interfaces

<!-- snippet: fragment -->
```go
type LLM interface {
    ToConfig() map[string]interface{}
}

type TTS interface {
    ToConfig() map[string]interface{}
    GetSampleRate() *SampleRate
}

type STT interface {
    ToConfig() map[string]interface{}
}

type MLLM interface {
    ToConfig() map[string]interface{}
}

type Avatar interface {
    ToConfig() map[string]interface{}
    RequiredSampleRate() SampleRate
}
```

## LLM Vendors

| Constructor | Options Struct | Required Fields | Default Model |
|---|---|---|---|
| `NewOpenAI` | `OpenAIOptions` | `APIKey` for BYOK; none for supported preset-backed OpenAI models | `gpt-4o-mini` |
| `NewAzureOpenAI` | `AzureOpenAIOptions` | `APIKey`, `Endpoint`, `DeploymentName` | — |
| `NewAnthropic` | `AnthropicOptions` | `APIKey` | `claude-3-5-sonnet-20241022` |
| `NewGemini` | `GeminiOptions` | `APIKey` | `gemini-2.0-flash-exp` |

<!-- snippet: fragment -->
```go
llm := vendors.NewOpenAI(vendors.OpenAIOptions{
    Model: "gpt-5-mini",
})

agent := agentkit.NewAgent(...).WithLlm(llm)
```

## TTS Vendors

| Constructor | Options Struct | Required Fields |
|---|---|---|
| `NewElevenLabsTTS` | `ElevenLabsTTSOptions` | `Key`, `ModelID`, `VoiceID` |
| `NewMicrosoftTTS` | `MicrosoftTTSOptions` | `Key`, `Region`, `VoiceName` |
| `NewOpenAITTS` | `OpenAITTSOptions` | `Voice` for preset-backed `tts-1`; `APIKey`, `Voice` for BYOK |
| `NewCartesiaTTS` | `CartesiaTTSOptions` | `Key`, `VoiceID` |
| `NewGoogleTTS` | `GoogleTTSOptions` | `Key`, `VoiceName` |
| `NewAmazonTTS` | `AmazonTTSOptions` | `AccessKey`, `SecretKey`, `Region`, `VoiceID` |
| `NewHumeAITTS` | `HumeAITTSOptions` | `Key` |
| `NewRimeTTS` | `RimeTTSOptions` | `Key`, `Speaker` |
| `NewFishAudioTTS` | `FishAudioTTSOptions` | `Key`, `ReferenceID` |
| `NewGroqTTS` | `GroqTTSOptions` | `Key` |
| `NewMiniMaxTTS` | `MiniMaxTTSOptions` | `Model` for supported preset-backed MiniMax models; `Key`, `GroupID`, `Model` for BYOK |
| `NewSarvamTTS` | `SarvamTTSOptions` | `APIKey` |

<!-- snippet: fragment -->
```go
tts := vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
    Key:        "<key>",
    ModelID:    "eleven_turbo_v2_5",
    VoiceID:    "<voice_id>",
    SampleRate: &vendors.SampleRate24kHz,
})

agent = agent.WithTts(tts)
```

### SampleRate Constants

<!-- snippet: fragment -->
```go
vendors.SampleRate8kHz   // 8000
vendors.SampleRate16kHz  // 16000
vendors.SampleRate22kHz  // 22050
vendors.SampleRate24kHz  // 24000
vendors.SampleRate44kHz  // 44100
vendors.SampleRate48kHz  // 48000
```

Note: `OpenAITTS` always returns `SampleRate24kHz`. Other TTS vendors return their configured sample rate or `nil`.

## STT Vendors

| Constructor | Options Struct | Required Fields |
|---|---|---|
| `NewSpeechmaticsSTT` | `SpeechmaticsSTTOptions` | `APIKey` |
| `NewDeepgramSTT` | `DeepgramSTTOptions` | `APIKey` for BYOK; none for supported preset-backed Deepgram models |
| `NewMicrosoftSTT` | `MicrosoftSTTOptions` | `Key`, `Region` |
| `NewOpenAISTT` | `OpenAISTTOptions` | `APIKey` |
| `NewGoogleSTT` | `GoogleSTTOptions` | `Key` |
| `NewAmazonSTT` | `AmazonSTTOptions` | `AccessKey`, `SecretKey`, `Region` |
| `NewAssemblyAISTT` | `AssemblyAISTTOptions` | `APIKey` |
| `NewAresSTT` | `AresSTTOptions` | `APIKey` |
| `NewSonioxSTT` | `SonioxSTTOptions` | `APIKey` |
| `NewSarvamSTT` | `SarvamSTTOptions` | `APIKey` |

<!-- snippet: fragment -->
```go
stt := vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
    APIKey:   "<key>",
    Model:    "nova-2",
    Language: "en-US",
})

agent = agent.WithStt(stt)
```

## MLLM Vendors

| Constructor | Options Struct | Required Fields | Default Model |
|---|---|---|---|
| `NewOpenAIRealtime` | `OpenAIRealtimeOptions` | `APIKey` | `gpt-4o-realtime-preview` |
| `NewVertexAI` | `VertexAIOptions` | `ProjectID` | `gemini-2.0-flash-exp` |

<!-- snippet: fragment -->
```go
mllm := vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
    APIKey: "<key>",
    Model:  "gpt-4o-realtime-preview",
    Voice:  "alloy",
})

agent = agent.WithMllm(mllm)
```

## Avatar Vendors

| Constructor | Options Struct | Required Fields | Required TTS Sample Rate |
|---|---|---|---|
| `NewHeyGenAvatar` | `HeyGenAvatarOptions` | `APIKey`, `Quality`, `AgoraUID` | 24kHz |
| `NewAkoolAvatar` | `AkoolAvatarOptions` | `APIKey` | 16kHz |

<!-- snippet: fragment -->
```go
avatar := vendors.NewHeyGenAvatar(vendors.HeyGenAvatarOptions{
    APIKey:   "<key>",
    Quality:  "high",
    AgoraUID: "2001",
})

// TTS must be configured with matching sample rate BEFORE WithAvatar
agent = agent.WithTts(tts).WithAvatar(avatar)
```

See [Avatars Guide](../guides/avatars.md) for sample rate requirements and the panic behavior when rates mismatch.

## Validation

All vendor constructors validate required fields and `panic` if they are missing. For example:

<!-- snippet: fragment -->
```go
// This panics because gpt-4o is not a supported preset-backed model.
vendors.NewOpenAI(vendors.OpenAIOptions{Model: "gpt-4o"})
```

This is Go-idiomatic for configuration errors that indicate programmer mistakes rather than runtime conditions. Handle these by ensuring BYOK fields are populated when you are not using a supported preset-backed model.
