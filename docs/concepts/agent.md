---
sidebar_position: 2
title: Agent
description: The Agent builder — functional options pattern, vendor chaining, and ToProperties conversion.
---

# Agent

The `agentkit.Agent` is the central configuration object. It defines what LLM, TTS, STT, MLLM, and avatar vendors your agent uses, along with instructions, greeting messages, and advanced settings.

## Functional Options Pattern

`agentkit.NewAgent` uses Go's [functional options pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis). Instead of a large config struct, you pass option functions that each configure one aspect of the agent:

```go
agent := agentkit.NewAgent(
    agentkit.WithName("my-assistant"),
    agentkit.WithInstructions("You are a helpful voice assistant."),
    agentkit.WithGreeting("Hello! How can I help?"),
    agentkit.WithFailureMessage("Sorry, something went wrong."),
    agentkit.WithMaxHistory(10),
)
```

Each `With*` function has the signature `func(...) AgentOption`, where `AgentOption` is `func(*Agent)`. This pattern lets you:

- Omit any option you don't need (sensible defaults)
- Add new options without breaking existing code
- Compose options from helper functions

## AgentOption Functions

These are passed to `agentkit.NewAgent(opts ...AgentOption)`:

| Function | Parameter | Description |
|---|---|---|
| `WithName(name string)` | Agent name | Identifier for the agent |
| `WithInstructions(instructions string)` | System prompt | Injected as the LLM system message |
| `WithGreeting(greeting string)` | Greeting text | First message the agent speaks |
| `WithFailureMessage(msg string)` | Fallback message | Spoken when the LLM fails |
| `WithMaxHistory(n int)` | History depth | Max conversation turns to retain |
| `WithTurnDetectionConfig(td *TurnDetectionConfig)` | VAD config | Server-side voice activity detection |
| `WithSalConfig(sal *SalConfig)` | SAL config | Speech analytics configuration |
| `WithAdvancedFeatures(af *AdvancedFeatures)` | Feature flags | `EnableMllm`, `EnableAivad`, etc. |
| `WithParameters(params *SessionParams)` | Session params | Additional session parameters |

## Vendor Chaining Methods

After creating an agent with `NewAgent`, attach vendors using method chaining. Each method returns a **new** `*Agent` (the original is not modified — immutable cloning):

```go
agent := agentkit.NewAgent(
    agentkit.WithName("assistant"),
    agentkit.WithInstructions("You are helpful."),
).WithLlm(
    vendors.NewOpenAI(vendors.OpenAIOptions{
        APIKey: "<key>",
        Model:  "gpt-4o-mini",
    }),
).WithTts(
    vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
        Key:     "<key>",
        ModelID: "eleven_turbo_v2_5",
        VoiceID: "<voice_id>",
    }),
).WithStt(
    vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
        APIKey: "<key>",
    }),
)
```

| Method | Parameter Type | Description |
|---|---|---|
| `WithLlm(vendor vendors.LLM)` | `vendors.LLM` interface | Set the LLM vendor |
| `WithTts(vendor vendors.TTS)` | `vendors.TTS` interface | Set the TTS vendor (also captures sample rate) |
| `WithStt(vendor vendors.STT)` | `vendors.STT` interface | Set the STT vendor |
| `WithMllm(vendor vendors.MLLM)` | `vendors.MLLM` interface | Set the MLLM vendor (for multimodal flow) |
| `WithAvatar(vendor vendors.Avatar)` | `vendors.Avatar` interface | Set the avatar vendor (validates sample rate) |
| `WithTurnDetection(td *TurnDetectionConfig)` | Pointer to config | Override turn detection on a cloned agent |
| `WithInstructions(instructions string)` | String | Override instructions on a cloned agent |
| `WithGreeting(greeting string)` | String | Override greeting on a cloned agent |
| `WithName(name string)` | String | Override name on a cloned agent |

Note: `WithInstructions`, `WithGreeting`, and `WithName` exist as both `AgentOption` functions (for `NewAgent`) and as methods on `*Agent` (for cloning with overrides).

## Agent Getters

```go
agent.Name() string
agent.Instructions() string
agent.Greeting() string
agent.LlmConfig() map[string]interface{}
agent.TtsConfig() map[string]interface{}
agent.SttConfig() map[string]interface{}
agent.MllmConfig() map[string]interface{}
agent.TtsSampleRate() *vendors.SampleRate
agent.AvatarRequiredSampleRate() *vendors.SampleRate
```

## ToProperties

`ToProperties` converts the agent configuration into a `*Agora.StartAgentsRequestProperties` suitable for the API:

```go
props, err := agent.ToProperties(agentkit.ToPropertiesOptions{
    Channel:        "my-channel",
    AgentUID:       "1001",
    RemoteUIDs:     []string{"1002"},
    AppID:          "<app_id>",
    AppCertificate: "<app_cert>",
})
if err != nil {
    log.Fatalf("Failed to build properties: %v", err)
}
```

### ToPropertiesOptions Fields

| Field | Type | Description |
|---|---|---|
| `Channel` | `string` | Agora channel name |
| `AgentUID` | `string` | UID for the agent in the channel |
| `RemoteUIDs` | `[]string` | UIDs of remote participants |
| `Token` | `string` | Pre-generated RTC token (if provided, skips generation) |
| `AppID` | `string` | Agora App ID (for token generation) |
| `AppCertificate` | `string` | Agora App Certificate (for token generation) |
| `UID` | `uint32` | Numeric UID for token generation |
| `TokenExpirySeconds` | `int` | Token expiry (default: 3600) |
| `IdleTimeout` | `*int` | Session idle timeout in seconds |
| `EnableStringUID` | `*bool` | Enable string UIDs |

If `Token` is empty, `ToProperties` generates one from `AppID` + `AppCertificate`. If both are empty, it returns an error.

In cascading mode, `ToProperties` requires both LLM and TTS to be configured — it returns an error if either is missing. In MLLM mode (when `AdvancedFeatures.EnableMllm` is `true`), LLM and TTS are not required.

## Type Aliases

The agentkit package defines convenient type aliases for common Agora types:

```go
type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters
```
