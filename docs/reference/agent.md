---
sidebar_position: 2
title: Agent Reference
description: Complete API reference for agentkit.Agent — functional options, methods, and ToProperties.
---

# Agent Reference

Package: `github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit`

## NewAgent

<!-- snippet: fragment -->
```go
func NewAgent(opts ...AgentOption) *Agent
```

Creates a new `Agent` with the given functional options.

## AgentOption Type

<!-- snippet: fragment -->
```go
type AgentOption func(*Agent)
```

## AgentOption Functions

### WithName

<!-- snippet: fragment -->
```go
func WithName(name string) AgentOption
```

Sets the agent name identifier.

### WithInstructions

<!-- snippet: fragment -->
```go
func WithInstructions(instructions string) AgentOption
```

Sets the system prompt injected into the LLM configuration.

### WithGreeting

<!-- snippet: fragment -->
```go
func WithGreeting(greeting string) AgentOption
```

Sets the greeting message the agent speaks first.

### WithFailureMessage

<!-- snippet: fragment -->
```go
func WithFailureMessage(msg string) AgentOption
```

Sets the fallback message spoken when the LLM fails.

### WithMaxHistory

<!-- snippet: fragment -->
```go
func WithMaxHistory(n int) AgentOption
```

Sets the maximum number of conversation turns to retain.

### WithTurnDetectionConfig

<!-- snippet: fragment -->
```go
func WithTurnDetectionConfig(td *TurnDetectionConfig) AgentOption
```

Sets the turn detection configuration. Use `Config.StartOfSpeech` and `Config.EndOfSpeech` for the preferred SOS/EOS model.

### WithSalConfig

<!-- snippet: fragment -->
```go
func WithSalConfig(sal *SalConfig) AgentOption
```

Sets the speech analytics configuration.

### WithAdvancedFeatures

<!-- snippet: fragment -->
```go
func WithAdvancedFeatures(af *AdvancedFeatures) AgentOption
```

Sets advanced feature flags (e.g., `EnableMllm`, `EnableAivad`).

### WithParameters

<!-- snippet: fragment -->
```go
func WithParameters(params *SessionParams) AgentOption
```

Sets additional session parameters.

### WithGeofence

<!-- snippet: fragment -->
```go
func WithGeofence(gf *GeofenceConfig) AgentOption
```

Sets geofence configuration (restricts backend server regions).

### WithLabels

<!-- snippet: fragment -->
```go
func WithLabels(labels map[string]string) AgentOption
```

Sets custom labels (key-value pairs returned in notification callbacks).

### WithRtc

<!-- snippet: fragment -->
```go
func WithRtc(rtc *RtcConfig) AgentOption
```

Sets RTC configuration.

### WithFillerWords

<!-- snippet: fragment -->
```go
func WithFillerWords(fw *FillerWordsConfig) AgentOption
```

Sets filler words configuration (played while waiting for LLM response).

## Agent Methods

All vendor-chaining methods return a **new** `*Agent` (immutable clone). The original agent is not modified.

### WithLlm

<!-- snippet: fragment -->
```go
func (a *Agent) WithLlm(vendor vendors.LLM) *Agent
```

### WithTts

<!-- snippet: fragment -->
```go
func (a *Agent) WithTts(vendor vendors.TTS) *Agent
```

Also captures the vendor's sample rate for avatar validation.

### WithStt

<!-- snippet: fragment -->
```go
func (a *Agent) WithStt(vendor vendors.STT) *Agent
```

### WithMllm

<!-- snippet: fragment -->
```go
func (a *Agent) WithMllm(vendor vendors.MLLM) *Agent
```

### WithAvatar

<!-- snippet: fragment -->
```go
func (a *Agent) WithAvatar(vendor vendors.Avatar) *Agent
```

**Panics** if TTS is already configured with a sample rate that doesn't match the avatar's requirement. See [Avatars Guide](../guides/avatars.md).

### WithTurnDetection (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithTurnDetection(td *TurnDetectionConfig) *Agent
```

Sets the turn detection configuration. Use `Config.StartOfSpeech` and `Config.EndOfSpeech` for the preferred SOS/EOS model.

### WithInstructions (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithInstructions(instructions string) *Agent
```

### WithGreeting (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithGreeting(greeting string) *Agent
```

### WithName (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithName(name string) *Agent
```

### WithSal (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithSal(sal *SalConfig) *Agent
```

### WithAdvancedFeatures (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithAdvancedFeatures(af *AdvancedFeatures) *Agent
```

### WithParameters (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithParameters(params *SessionParams) *Agent
```

### WithFailureMessage (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithFailureMessage(msg string) *Agent
```

### WithMaxHistory (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithMaxHistory(n int) *Agent
```

### WithGeofence (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithGeofence(gf *GeofenceConfig) *Agent
```

### WithLabels (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithLabels(labels map[string]string) *Agent
```

### WithRtc (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithRtc(rtc *RtcConfig) *Agent
```

### WithFillerWords (method)

<!-- snippet: fragment -->
```go
func (a *Agent) WithFillerWords(fw *FillerWordsConfig) *Agent
```

## Getters

<!-- snippet: fragment -->
```go
func (a *Agent) Name() string
func (a *Agent) Instructions() string
func (a *Agent) Greeting() string
func (a *Agent) FailureMessage() string
func (a *Agent) MaxHistory() *int
func (a *Agent) LlmConfig() map[string]interface{}
func (a *Agent) TtsConfig() map[string]interface{}
func (a *Agent) SttConfig() map[string]interface{}
func (a *Agent) MllmConfig() map[string]interface{}
func (a *Agent) TtsSampleRate() *vendors.SampleRate
func (a *Agent) AvatarRequiredSampleRate() *vendors.SampleRate
func (a *Agent) Avatar() map[string]interface{}
func (a *Agent) TurnDetection() *TurnDetectionConfig
func (a *Agent) Sal() *SalConfig
func (a *Agent) AdvancedFeatures() *AdvancedFeatures
func (a *Agent) Parameters() *SessionParams
func (a *Agent) Geofence() *GeofenceConfig
func (a *Agent) Labels() map[string]string
func (a *Agent) Rtc() *RtcConfig
func (a *Agent) FillerWords() *FillerWordsConfig
```

## ToProperties

<!-- snippet: fragment -->
```go
func (a *Agent) ToProperties(opts ToPropertiesOptions) (*Agora.StartAgentsRequestProperties, error)
```

Converts the agent configuration into API request properties. Handles token generation, LLM/TTS config merging, and validation.

Returns an error if:
- Neither `Token` nor `AppID`+`AppCertificate` is provided
- In cascading mode: LLM or TTS is not configured
- Config marshaling fails

### ToPropertiesOptions

<!-- snippet: fragment -->
```go
type ToPropertiesOptions struct {
    Channel         string
    AgentUID        string
    RemoteUIDs      []string
    Token           string
    AppID           string
    AppCertificate  string
    ExpiresIn       int
    IdleTimeout     *int
    EnableStringUID *bool
}
```

| Field | Type | Description |
|---|---|---|
| `Channel` | `string` | Agora channel name |
| `AgentUID` | `string` | Agent's UID in the channel |
| `RemoteUIDs` | `[]string` | Remote participant UIDs |
| `Token` | `string` | Pre-generated RTC+RTM token (skips generation if set) |
| `AppID` | `string` | Agora App ID (for token generation) |
| `AppCertificate` | `string` | Agora App Certificate (for token generation) |
| `ExpiresIn` | `int` | Token lifetime in seconds (default: `86400` = 24 h, Agora max). Use `ExpiresInHours()` / `ExpiresInMinutes()` for clarity. Valid range: 1–86400. |
| `IdleTimeout` | `*int` | Session idle timeout |
| `EnableStringUID` | `*bool` | Enable string UID mode |

## Type Aliases

<!-- snippet: fragment -->
```go
type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters
type GeofenceConfig = Agora.StartAgentsRequestPropertiesGeofence
type RtcConfig = Agora.StartAgentsRequestPropertiesRtc
type FillerWordsConfig = Agora.StartAgentsRequestPropertiesFillerWords
type LlmConfig = Agora.StartAgentsRequestPropertiesLlm
type MllmConfig = Agora.StartAgentsRequestPropertiesMllm
type AsrConfig = Agora.StartAgentsRequestPropertiesAsr
type TtsConfig = Agora.Tts
type AvatarConfig = Agora.StartAgentsRequestPropertiesAvatar
```

Additional SOS/EOS turn detection aliases: `TurnDetectionNestedConfig`, `StartOfSpeechConfig`, `EndOfSpeechConfig`, and related sub-types. See the agentkit package for the full list.

## Token Generation

<!-- snippet: fragment -->
```go
func GenerateRtcToken(opts GenerateTokenOptions) (string, error)
```

### GenerateTokenOptions

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
const (
    RolePublisher      = 1
    RoleSubscriber     = 2
    DefaultExpirySeconds = 3600
)
```
