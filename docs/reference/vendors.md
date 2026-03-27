---
sidebar_position: 4
title: Vendors Reference
description: Complete API reference for all vendor constructors and configuration structs.
---

# Vendors Reference

Package: `github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors`

## SampleRate

```go
type SampleRate int

const (
    SampleRate8kHz  SampleRate = 8000
    SampleRate16kHz SampleRate = 16000
    SampleRate22kHz SampleRate = 22050
    SampleRate24kHz SampleRate = 24000
    SampleRate44kHz SampleRate = 44100
    SampleRate48kHz SampleRate = 48000
)
```

## Interfaces

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

---

## LLM Vendors

### NewOpenAI

```go
func NewOpenAI(opts OpenAIOptions) *OpenAI
```

Panics if `APIKey` is empty.

#### OpenAIOptions

| Field             | Type                       | Required | Default                                        | Description             |
| ----------------- | -------------------------- | -------- | ---------------------------------------------- | ----------------------- |
| `APIKey`          | `string`                   | Yes      | —                                              | OpenAI API key          |
| `Model`           | `string`                   | No       | `"gpt-4o-mini"`                                | Model identifier        |
| `BaseURL`         | `string`                   | No       | `"https://api.openai.com/v1/chat/completions"` | API endpoint            |
| `Temperature`     | `*float64`                 | No       | —                                              | Sampling temperature    |
| `TopP`            | `*float64`                 | No       | —                                              | Nucleus sampling        |
| `MaxTokens`       | `*int`                     | No       | —                                              | Max tokens in response  |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                                              | System messages         |
| `GreetingMessage` | `string`                   | No       | —                                              | Initial greeting        |
| `FailureMessage`  | `string`                   | No       | —                                              | Fallback on error       |
| `InputModalities` | `[]string`                 | No       | `["text"]`                                     | Input modality types    |
| `Params`          | `map[string]interface{}`   | No       | —                                              | Additional model params |

### NewAzureOpenAI

```go
func NewAzureOpenAI(opts AzureOpenAIOptions) *AzureOpenAI
```

Panics if `APIKey`, `Endpoint`, or `DeploymentName` is empty.

#### AzureOpenAIOptions

| Field             | Type                       | Required | Default                | Description          |
| ----------------- | -------------------------- | -------- | ---------------------- | -------------------- |
| `APIKey`          | `string`                   | Yes      | —                      | Azure OpenAI API key |
| `Endpoint`        | `string`                   | Yes      | —                      | Azure endpoint URL   |
| `DeploymentName`  | `string`                   | Yes      | —                      | Deployment name      |
| `APIVersion`      | `string`                   | No       | `"2024-08-01-preview"` | API version          |
| `Temperature`     | `*float64`                 | No       | —                      | Sampling temperature |
| `TopP`            | `*float64`                 | No       | —                      | Nucleus sampling     |
| `MaxTokens`       | `*int`                     | No       | —                      | Max tokens           |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                      | System messages      |
| `GreetingMessage` | `string`                   | No       | —                      | Initial greeting     |
| `FailureMessage`  | `string`                   | No       | —                      | Fallback on error    |
| `InputModalities` | `[]string`                 | No       | `["text"]`             | Input modality types |

### NewAnthropic

```go
func NewAnthropic(opts AnthropicOptions) *Anthropic
```

Panics if `APIKey` is empty.

#### AnthropicOptions

| Field             | Type                       | Required | Default                        | Description          |
| ----------------- | -------------------------- | -------- | ------------------------------ | -------------------- |
| `APIKey`          | `string`                   | Yes      | —                              | Anthropic API key    |
| `Model`           | `string`                   | No       | `"claude-3-5-sonnet-20241022"` | Model identifier     |
| `MaxTokens`       | `*int`                     | No       | —                              | Max tokens           |
| `Temperature`     | `*float64`                 | No       | —                              | Sampling temperature |
| `TopP`            | `*float64`                 | No       | —                              | Nucleus sampling     |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                              | System messages      |
| `GreetingMessage` | `string`                   | No       | —                              | Initial greeting     |
| `FailureMessage`  | `string`                   | No       | —                              | Fallback on error    |
| `InputModalities` | `[]string`                 | No       | `["text"]`                     | Input modality types |

### NewGemini

```go
func NewGemini(opts GeminiOptions) *Gemini
```

Panics if `APIKey` is empty.

#### GeminiOptions

| Field             | Type                       | Required | Default                  | Description          |
| ----------------- | -------------------------- | -------- | ------------------------ | -------------------- |
| `APIKey`          | `string`                   | Yes      | —                        | Google AI API key    |
| `Model`           | `string`                   | No       | `"gemini-2.0-flash-exp"` | Model identifier     |
| `Temperature`     | `*float64`                 | No       | —                        | Sampling temperature |
| `TopP`            | `*float64`                 | No       | —                        | Nucleus sampling     |
| `TopK`            | `*int`                     | No       | —                        | Top-K sampling       |
| `MaxOutputTokens` | `*int`                     | No       | —                        | Max output tokens    |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                        | System messages      |
| `GreetingMessage` | `string`                   | No       | —                        | Initial greeting     |
| `FailureMessage`  | `string`                   | No       | —                        | Fallback on error    |
| `InputModalities` | `[]string`                 | No       | `["text"]`               | Input modality types |

---

## TTS Vendors

### NewElevenLabsTTS

```go
func NewElevenLabsTTS(opts ElevenLabsTTSOptions) *ElevenLabsTTS
```

Panics if `Key`, `ModelID`, or `VoiceID` is empty.

#### ElevenLabsTTSOptions

| Field          | Type          | Required | Description                                    |
| -------------- | ------------- | -------- | ---------------------------------------------- |
| `Key`          | `string`      | Yes      | ElevenLabs API key                             |
| `ModelID`      | `string`      | Yes      | Model identifier (e.g., `"eleven_turbo_v2_5"`) |
| `VoiceID`      | `string`      | Yes      | Voice identifier                               |
| `BaseURL`      | `string`      | No       | Custom API endpoint                            |
| `SampleRate`   | `*SampleRate` | No       | Output sample rate                             |
| `SkipPatterns` | `[]int`       | No       | Patterns to skip in TTS output                 |

### NewMicrosoftTTS

```go
func NewMicrosoftTTS(opts MicrosoftTTSOptions) *MicrosoftTTS
```

Panics if `Key`, `Region`, or `VoiceName` is empty.

#### MicrosoftTTSOptions

| Field          | Type          | Required | Description                              |
| -------------- | ------------- | -------- | ---------------------------------------- |
| `Key`          | `string`      | Yes      | Azure Speech Services key                |
| `Region`       | `string`      | Yes      | Azure region (e.g., `"eastus"`)          |
| `VoiceName`    | `string`      | Yes      | Voice name (e.g., `"en-US-JennyNeural"`) |
| `SampleRate`   | `*SampleRate` | No       | Output sample rate                       |
| `SkipPatterns` | `[]int`       | No       | Patterns to skip                         |

### NewOpenAITTS

```go
func NewOpenAITTS(opts OpenAITTSOptions) *OpenAITTS
```

Panics if `APIKey` or `Voice` is empty. Always returns `SampleRate24kHz` from `GetSampleRate()`.

#### OpenAITTSOptions

| Field            | Type       | Required | Description                        |
| ---------------- | ---------- | -------- | ---------------------------------- |
| `APIKey`         | `string`   | Yes      | OpenAI API key                     |
| `Voice`          | `string`   | Yes      | Voice name                         |
| `Model`          | `string`   | No       | Model identifier                   |
| `ResponseFormat` | `string`   | No       | Audio format (e.g., `"pcm"`)       |
| `Speed`          | `*float64` | No       | Speech speed multiplier            |
| `SkipPatterns`   | `[]int`    | No       | Patterns to skip                   |

### NewCartesiaTTS

```go
func NewCartesiaTTS(opts CartesiaTTSOptions) *CartesiaTTS
```

Panics if `APIKey` or `VoiceID` is empty.

#### CartesiaTTSOptions

| Field          | Type          | Required | Description                                          |
| -------------- | ------------- | -------- | ---------------------------------------------------- |
| `APIKey`       | `string`      | Yes      | Cartesia API key                                     |
| `VoiceID`      | `string`      | Yes      | Voice identifier (serialized as `{"mode":"id","id":"..."}`) |
| `ModelID`      | `string`      | No       | Model identifier                                     |
| `SampleRate`   | `*SampleRate` | No       | Output sample rate                                   |
| `SkipPatterns` | `[]int`       | No       | Patterns to skip                                     |

### NewGoogleTTS

```go
func NewGoogleTTS(opts GoogleTTSOptions) *GoogleTTS
```

Panics if `Key` or `VoiceName` is empty.

#### GoogleTTSOptions

| Field          | Type     | Required | Description          |
| -------------- | -------- | -------- | -------------------- |
| `Key`          | `string` | Yes      | Google Cloud API key |
| `VoiceName`    | `string` | Yes      | Voice name           |
| `LanguageCode` | `string` | No       | Language code        |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip     |

### NewAmazonTTS

```go
func NewAmazonTTS(opts AmazonTTSOptions) *AmazonTTS
```

Panics if `AccessKey`, `SecretKey`, `Region`, or `VoiceID` is empty.

#### AmazonTTSOptions

| Field          | Type     | Required | Description      |
| -------------- | -------- | -------- | ---------------- |
| `AccessKey`    | `string` | Yes      | AWS access key   |
| `SecretKey`    | `string` | Yes      | AWS secret key   |
| `Region`       | `string` | Yes      | AWS region       |
| `VoiceID`      | `string` | Yes      | Polly voice ID   |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip |

### NewHumeAITTS

```go
func NewHumeAITTS(opts HumeAITTSOptions) *HumeAITTS
```

Panics if `Key` is empty.

#### HumeAITTSOptions

| Field          | Type     | Required | Description      |
| -------------- | -------- | -------- | ---------------- |
| `Key`          | `string` | Yes      | Hume AI API key  |
| `ConfigID`     | `string` | No       | Configuration ID |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip |

### NewRimeTTS

```go
func NewRimeTTS(opts RimeTTSOptions) *RimeTTS
```

Panics if `Key` or `Speaker` is empty.

#### RimeTTSOptions

| Field          | Type       | Required | Description                                   |
| -------------- | ---------- | -------- | --------------------------------------------- |
| `Key`          | `string`   | Yes      | Rime API key                                  |
| `Speaker`      | `string`   | Yes      | Speaker identifier                            |
| `ModelID`      | `string`   | No       | Model identifier                              |
| `Lang`         | `string`   | No       | Language code                                 |
| `SamplingRate` | `*int`     | No       | Sampling rate in Hz (serialized as `samplingRate`) |
| `SpeedAlpha`   | `*float64` | No       | Speed multiplier (serialized as `speedAlpha`) |
| `SkipPatterns` | `[]int`    | No       | Patterns to skip                              |

### NewFishAudioTTS

```go
func NewFishAudioTTS(opts FishAudioTTSOptions) *FishAudioTTS
```

Panics if `Key` or `ReferenceID` is empty.

#### FishAudioTTSOptions

| Field          | Type     | Required | Description        |
| -------------- | -------- | -------- | ------------------ |
| `Key`          | `string` | Yes      | FishAudio API key  |
| `ReferenceID`  | `string` | Yes      | Reference audio ID |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip   |

### NewMiniMaxTTS

```go
func NewMiniMaxTTS(opts MiniMaxTTSOptions) *MiniMaxTTS
```

Panics if `Key`, `GroupID`, `Model`, `VoiceID`, or `URL` is empty.

#### MiniMaxTTSOptions

| Field          | Type     | Required | Description                               |
| -------------- | -------- | -------- | ----------------------------------------- |
| `Key`          | `string` | Yes      | MiniMax API key                           |
| `GroupID`      | `string` | Yes      | MiniMax group ID                          |
| `Model`        | `string` | Yes      | Model name (e.g., `speech-02-turbo`)      |
| `VoiceID`      | `string` | Yes      | Voice style identifier                    |
| `URL`          | `string` | Yes      | WebSocket endpoint                        |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip                          |

### NewMurfTTS

```go
func NewMurfTTS(opts MurfTTSOptions) *MurfTTS
```

Panics if `Key` or `VoiceID` is empty.

#### MurfTTSOptions

| Field          | Type     | Required | Description                              |
| -------------- | -------- | -------- | ---------------------------------------- |
| `Key`          | `string` | Yes      | Murf API key                             |
| `VoiceID`      | `string` | Yes      | Voice ID (e.g., `Ariana`, `Natalie`)     |
| `Style`        | `string` | No       | Voice style (e.g., `Conversational`)     |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip                         |

### NewSarvamTTS

```go
func NewSarvamTTS(opts SarvamTTSOptions) *SarvamTTS
```

Panics if `Key`, `Speaker`, or `TargetLanguageCode` is empty.

#### SarvamTTSOptions

| Field                | Type     | Required | Description          |
| -------------------- | -------- | -------- | -------------------- |
| `Key`                | `string` | Yes      | Sarvam API key       |
| `Speaker`            | `string` | Yes      | Speaker name         |
| `TargetLanguageCode` | `string` | Yes      | Target language code |
| `SkipPatterns`       | `[]int`  | No       | Patterns to skip     |

---

## STT Vendors

### NewSpeechmaticsSTT

```go
func NewSpeechmaticsSTT(opts SpeechmaticsSTTOptions) *SpeechmaticsSTT
```

Panics if `APIKey` is empty.

#### SpeechmaticsSTTOptions

| Field      | Type     | Required | Description          |
| ---------- | -------- | -------- | -------------------- |
| `APIKey`   | `string` | Yes      | Speechmatics API key |
| `Language` | `string` | No       | Language code        |
| `Model`    | `string` | No       | Model identifier     |

### NewDeepgramSTT

```go
func NewDeepgramSTT(opts DeepgramSTTOptions) *DeepgramSTT
```

Panics if `APIKey` is empty.

#### DeepgramSTTOptions

| Field      | Type     | Required | Description              |
| ---------- | -------- | -------- | ------------------------ |
| `APIKey`   | `string` | Yes      | Deepgram API key         |
| `Model`    | `string` | No       | Model (e.g., `"nova-2"`) |
| `Language` | `string` | No       | Language code            |

### NewMicrosoftSTT

```go
func NewMicrosoftSTT(opts MicrosoftSTTOptions) *MicrosoftSTT
```

Panics if `Key` or `Region` is empty.

#### MicrosoftSTTOptions

| Field      | Type     | Required | Description               |
| ---------- | -------- | -------- | ------------------------- |
| `Key`      | `string` | Yes      | Azure Speech Services key |
| `Region`   | `string` | Yes      | Azure region              |
| `Language` | `string` | No       | Language code             |

### NewOpenAISTT

```go
func NewOpenAISTT(opts OpenAISTTOptions) *OpenAISTT
```

Panics if `APIKey` is empty.

#### OpenAISTTOptions

| Field      | Type     | Required | Description      |
| ---------- | -------- | -------- | ---------------- |
| `APIKey`   | `string` | Yes      | OpenAI API key   |
| `Model`    | `string` | No       | Model identifier |
| `Language` | `string` | No       | Language code    |

### NewGoogleSTT

```go
func NewGoogleSTT(opts GoogleSTTOptions) *GoogleSTT
```

Panics if `Key` is empty.

#### GoogleSTTOptions

| Field      | Type     | Required | Description          |
| ---------- | -------- | -------- | -------------------- |
| `Key`      | `string` | Yes      | Google Cloud API key |
| `Language` | `string` | No       | Language code        |
| `Model`    | `string` | No       | Model identifier     |

### NewAmazonSTT

```go
func NewAmazonSTT(opts AmazonSTTOptions) *AmazonSTT
```

Panics if `AccessKey`, `SecretKey`, or `Region` is empty.

#### AmazonSTTOptions

| Field       | Type     | Required | Description    |
| ----------- | -------- | -------- | -------------- |
| `AccessKey` | `string` | Yes      | AWS access key |
| `SecretKey` | `string` | Yes      | AWS secret key |
| `Region`    | `string` | Yes      | AWS region     |
| `Language`  | `string` | No       | Language code  |

### NewAssemblyAISTT

```go
func NewAssemblyAISTT(opts AssemblyAISTTOptions) *AssemblyAISTT
```

Panics if `APIKey` is empty.

#### AssemblyAISTTOptions

| Field    | Type     | Required | Description        |
| -------- | -------- | -------- | ------------------ |
| `APIKey` | `string` | Yes      | AssemblyAI API key |

### NewAresSTT

```go
func NewAresSTT(opts AresSTTOptions) *AresSTT
```

Panics if `APIKey` is empty.

#### AresSTTOptions

| Field    | Type     | Required | Description  |
| -------- | -------- | -------- | ------------ |
| `APIKey` | `string` | Yes      | Ares API key |

### NewSonioxSTT

```go
func NewSonioxSTT(opts SonioxSTTOptions) *SonioxSTT
```

Panics if `APIKey` is empty.

#### SonioxSTTOptions

| Field    | Type     | Required | Description    |
| -------- | -------- | -------- | -------------- |
| `APIKey` | `string` | Yes      | Soniox API key |

### NewSarvamSTT

```go
func NewSarvamSTT(opts SarvamSTTOptions) *SarvamSTT
```

Panics if `APIKey` is empty.

#### SarvamSTTOptions

| Field      | Type     | Required | Description      |
| ---------- | -------- | -------- | ---------------- |
| `APIKey`   | `string` | Yes      | Sarvam API key   |
| `Language` | `string` | No       | Language code    |
| `Model`    | `string` | No       | Model identifier |

---

## MLLM Vendors

### NewOpenAIRealtime

```go
func NewOpenAIRealtime(opts OpenAIRealtimeOptions) *OpenAIRealtime
```

Panics if `APIKey` is empty.

#### OpenAIRealtimeOptions

| Field             | Type       | Required | Default                     | Description                                        |
| ----------------- | ---------- | -------- | --------------------------- | -------------------------------------------------- |
| `APIKey`          | `string`   | Yes      | —                           | OpenAI API key                                     |
| `Model`           | `string`   | No       | `"gpt-4o-realtime-preview"` | Model identifier                                   |
| `Voice`           | `string`   | No       | —                           | Voice name (e.g., `"alloy"`)                       |
| `Temperature`     | `*float64` | No       | —                           | Sampling temperature                               |
| `MaxOutputTokens` | `*int`     | No       | —                           | Max output tokens                                  |
| `SystemMessage`   | `string`   | No       | —                           | System instruction                                 |
| `PredefinedTools` | `[]string` | No       | —                           | Predefined tools (e.g., `["_publish_message"]`)    |
| `FailureMessage`  | `string`   | No       | —                           | Message played when the model call fails           |
| `MaxHistory`      | `*int`     | No       | —                           | Maximum conversation history length                |

### NewVertexAI

```go
func NewVertexAI(opts VertexAIOptions) *VertexAI
```

Panics if `ProjectID` is empty.

#### VertexAIOptions

| Field             | Type       | Required | Default                  | Description                                     |
| ----------------- | ---------- | -------- | ------------------------ | ----------------------------------------------- |
| `ProjectID`       | `string`   | Yes      | —                        | GCP project ID                                  |
| `Location`        | `string`   | No       | `"us-central1"`          | GCP region                                      |
| `Model`           | `string`   | No       | `"gemini-2.0-flash-exp"` | Model identifier                                |
| `Voice`           | `string`   | No       | —                        | Voice name                                      |
| `Language`        | `string`   | No       | —                        | Language code                                   |
| `SystemMessage`   | `string`   | No       | —                        | System instruction                              |
| `PredefinedTools` | `[]string` | No       | —                        | Predefined tools (e.g., `["_publish_message"]`) |
| `FailureMessage`  | `string`   | No       | —                        | Message played when the model call fails        |
| `MaxHistory`      | `*int`     | No       | —                        | Maximum conversation history length             |

---

## Avatar Vendors

### NewHeyGenAvatar

```go
func NewHeyGenAvatar(opts HeyGenAvatarOptions) *HeyGenAvatar
```

Panics if `APIKey` or `AgoraUID` is empty, or if `Quality` is not one of `"low"`, `"medium"`, `"high"`.

Required TTS sample rate: **24kHz** (`SampleRate24kHz`)

#### HeyGenAvatarOptions

| Field                 | Type     | Required | Description                                      |
| --------------------- | -------- | -------- | ------------------------------------------------ |
| `APIKey`              | `string` | Yes      | HeyGen API key                                   |
| `Quality`             | `string` | Yes      | `"low"`, `"medium"`, or `"high"`                 |
| `AgoraUID`            | `string` | Yes      | UID for avatar's video stream                    |
| `AgoraToken`          | `string` | No       | RTC token for avatar authentication              |
| `AvatarID`            | `string` | No       | HeyGen avatar ID                                 |
| `Enable`              | `*bool`  | No       | Enable or disable the avatar (default: `true`)   |
| `DisableIdleTimeout`  | `*bool`  | No       | Disable the idle timeout                         |
| `ActivityIdleTimeout` | `*int`   | No       | Idle timeout in seconds (default: 120)           |

### NewAkoolAvatar

```go
func NewAkoolAvatar(opts AkoolAvatarOptions) *AkoolAvatar
```

Panics if `APIKey` is empty.

Required TTS sample rate: **16kHz** (`SampleRate16kHz`)

#### AkoolAvatarOptions

| Field    | Type     | Required | Description   |
| -------- | -------- | -------- | ------------- |
| `APIKey` | `string` | Yes      | Akool API key |

---

## Sample Rate Constants

```go
const (
    HeyGenRequiredSampleRate = SampleRate24kHz  // 24000 Hz
    AkoolRequiredSampleRate  = SampleRate16kHz  // 16000 Hz
)
```
