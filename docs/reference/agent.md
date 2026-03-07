---
sidebar_position: 2
title: Agent Reference
description: Complete API reference for agentkit.Agent — functional options, methods, and ToProperties.
---

# Agent Reference

Package: `github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit`

## NewAgent

```go
func NewAgent(opts ...AgentOption) *Agent
```

Creates a new `Agent` with the given functional options.

## AgentOption Type

```go
type AgentOption func(*Agent)
```

## AgentOption Functions

### WithName

```go
func WithName(name string) AgentOption
```

Sets the agent name identifier.

### WithInstructions

```go
func WithInstructions(instructions string) AgentOption
```

Sets the system prompt injected into the LLM configuration.

### WithGreeting

```go
func WithGreeting(greeting string) AgentOption
```

Sets the greeting message the agent speaks first.

### WithFailureMessage

```go
func WithFailureMessage(msg string) AgentOption
```

Sets the fallback message spoken when the LLM fails.

### WithMaxHistory

```go
func WithMaxHistory(n int) AgentOption
```

Sets the maximum number of conversation turns to retain.

### WithTurnDetectionConfig

```go
func WithTurnDetectionConfig(td *TurnDetectionConfig) AgentOption
```

Sets the voice activity detection configuration.

### WithSalConfig

```go
func WithSalConfig(sal *SalConfig) AgentOption
```

Sets the speech analytics configuration.

### WithAdvancedFeatures

```go
func WithAdvancedFeatures(af *AdvancedFeatures) AgentOption
```

Sets advanced feature flags (e.g., `EnableMllm`, `EnableAivad`).

### WithParameters

```go
func WithParameters(params *SessionParams) AgentOption
```

Sets additional session parameters.

## Agent Methods

All vendor-chaining methods return a **new** `*Agent` (immutable clone). The original agent is not modified.

### WithLlm

```go
func (a *Agent) WithLlm(vendor vendors.LLM) *Agent
```

### WithTts

```go
func (a *Agent) WithTts(vendor vendors.TTS) *Agent
```

Also captures the vendor's sample rate for avatar validation.

### WithStt

```go
func (a *Agent) WithStt(vendor vendors.STT) *Agent
```

### WithMllm

```go
func (a *Agent) WithMllm(vendor vendors.MLLM) *Agent
```

### WithAvatar

```go
func (a *Agent) WithAvatar(vendor vendors.Avatar) *Agent
```

**Panics** if TTS is already configured with a sample rate that doesn't match the avatar's requirement. See [Avatars Guide](../guides/avatars.md).

### WithTurnDetection

```go
func (a *Agent) WithTurnDetection(td *TurnDetectionConfig) *Agent
```

### WithInstructions (method)

```go
func (a *Agent) WithInstructions(instructions string) *Agent
```

### WithGreeting (method)

```go
func (a *Agent) WithGreeting(greeting string) *Agent
```

### WithName (method)

```go
func (a *Agent) WithName(name string) *Agent
```

## Getters

```go
func (a *Agent) Name() string
func (a *Agent) Instructions() string
func (a *Agent) Greeting() string
func (a *Agent) LlmConfig() map[string]interface{}
func (a *Agent) TtsConfig() map[string]interface{}
func (a *Agent) SttConfig() map[string]interface{}
func (a *Agent) MllmConfig() map[string]interface{}
func (a *Agent) TtsSampleRate() *vendors.SampleRate
func (a *Agent) AvatarRequiredSampleRate() *vendors.SampleRate
```

## ToProperties

```go
func (a *Agent) ToProperties(opts ToPropertiesOptions) (*Agora.StartAgentsRequestProperties, error)
```

Converts the agent configuration into API request properties. Handles token generation, LLM/TTS config merging, and validation.

Returns an error if:
- Neither `Token` nor `AppID`+`AppCertificate` is provided
- In cascading mode: LLM or TTS is not configured
- Config marshaling fails

### ToPropertiesOptions

```go
type ToPropertiesOptions struct {
    Channel            string
    AgentUID           string
    RemoteUIDs         []string
    Token              string
    AppID              string
    AppCertificate     string
    UID                uint32
    TokenExpirySeconds int
    IdleTimeout        *int
    EnableStringUID    *bool
}
```

| Field | Type | Description |
|---|---|---|
| `Channel` | `string` | Agora channel name |
| `AgentUID` | `string` | Agent's UID in the channel |
| `RemoteUIDs` | `[]string` | Remote participant UIDs |
| `Token` | `string` | Pre-generated RTC token (skips generation if set) |
| `AppID` | `string` | Agora App ID (for token generation) |
| `AppCertificate` | `string` | Agora App Certificate (for token generation) |
| `UID` | `uint32` | Numeric UID for token generation |
| `TokenExpirySeconds` | `int` | Token TTL in seconds (default: 3600) |
| `IdleTimeout` | `*int` | Session idle timeout |
| `EnableStringUID` | `*bool` | Enable string UID mode |

## Type Aliases

```go
type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters
```

## Token Generation

```go
func GenerateRtcToken(opts GenerateTokenOptions) (string, error)
```

### GenerateTokenOptions

```go
type GenerateTokenOptions struct {
    AppID          string
    AppCertificate string
    Channel        string
    UID            uint32
    Role           int
    ExpirySeconds  int
}
```

### Constants

```go
const (
    RolePublisher      = 1
    RoleSubscriber     = 2
    DefaultExpirySeconds = 3600
)
```
