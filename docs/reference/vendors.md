---
sidebar_position: 4
title: Vendors Reference
description: Complete API reference for all vendor constructors and configuration structs.
---

# Vendors Reference

Package: `github.com/AgoraIO/agora-agents-go/v2/agentkit/vendors`

## SampleRate

<!-- snippet: fragment -->
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

---

## LLM Vendors

### NewOpenAI

<!-- snippet: fragment -->
```go
func NewOpenAI(opts OpenAIOptions) *OpenAI
```

Panics if `APIKey` is empty unless `Model` is one of the supported Agora-managed OpenAI models (`gpt-4o-mini`, `gpt-4.1-mini`, `gpt-5-nano`, `gpt-5-mini`) and `BaseURL` / `Vendor` are not set.

#### OpenAIOptions

| Field             | Type                       | Required | Default                                        | Description             |
| ----------------- | -------------------------- | -------- | ---------------------------------------------- | ----------------------- |
| `APIKey`          | `string`                   | No       | —                                              | OpenAI API key. Optional for supported Agora-managed OpenAI models. |
| `Model`           | `string`                   | No       | `"gpt-4o-mini"`                                | Model identifier        |
| `BaseURL`         | `string`                   | No       | `"https://api.openai.com/v1/chat/completions"` | API endpoint            |
| `Temperature`     | `*float64`                 | No       | —                                              | Sampling temperature    |
| `TopP`            | `*float64`                 | No       | —                                              | Nucleus sampling        |
| `MaxTokens`       | `*int`                     | No       | —                                              | Max tokens in response  |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                                              | System messages         |
| `GreetingMessage` | `string`                   | No       | —                                              | Initial greeting        |
| `FailureMessage`  | `string`                   | No       | —                                              | Fallback on error       |
| `InputModalities` | `[]string`                 | No       | `["text"]`                                     | Input modality types    |
| `OutputModalities` | `[]string`                | No       | —                                              | Output modality types   |
| `Params`          | `map[string]interface{}`   | No       | —                                              | Additional model params |
| `Headers`         | `map[string]string`        | No       | —                                              | Custom HTTP headers forwarded to the LLM provider |
| `GreetingConfigs` | `map[string]interface{}`   | No       | —                                              | Greeting playback configuration |
| `TemplateVariables` | `map[string]string`      | No       | —                                              | Template variables for messages |

### NewAzureOpenAI

<!-- snippet: fragment -->
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
| `Model`           | `string`                   | No       | —                      | Deployment's base model name (e.g., `"gpt-4o"`). Emitted as `params.model` for parity with the TypeScript SDK. |
| `Temperature`     | `*float64`                 | No       | —                      | Sampling temperature |
| `TopP`            | `*float64`                 | No       | —                      | Nucleus sampling     |
| `MaxTokens`       | `*int`                     | No       | —                      | Max tokens           |
| `SystemMessages`  | `[]map[string]interface{}` | No       | —                      | System messages      |
| `GreetingMessage` | `string`                   | No       | —                      | Initial greeting     |
| `FailureMessage`  | `string`                   | No       | —                      | Fallback on error    |
| `InputModalities` | `[]string`                 | No       | `["text"]`             | Input modality types |
| `OutputModalities` | `[]string`                | No       | —                      | Output modality types |
| `Params`          | `map[string]interface{}`   | No       | —                      | Additional model params |
| `Headers`         | `map[string]string`        | No       | —                      | Custom HTTP headers forwarded to the LLM provider |
| `GreetingConfigs` | `map[string]interface{}`   | No       | —                      | Greeting playback configuration |
| `TemplateVariables` | `map[string]string`      | No       | —                      | Template variables for messages |

### NewAnthropic

<!-- snippet: fragment -->
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
| `OutputModalities` | `[]string`                | No       | —                              | Output modality types |
| `Params`          | `map[string]interface{}`   | No       | —                              | Additional model params |
| `Headers`         | `map[string]string`        | No       | —                              | Custom HTTP headers forwarded to the LLM provider |
| `GreetingConfigs` | `map[string]interface{}`   | No       | —                              | Greeting playback configuration |
| `TemplateVariables` | `map[string]string`      | No       | —                              | Template variables for messages |

### NewGemini

<!-- snippet: fragment -->
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
| `OutputModalities` | `[]string`                | No       | —                        | Output modality types |
| `Params`          | `map[string]interface{}`   | No       | —                        | Additional model params |
| `Headers`         | `map[string]string`        | No       | —                        | Custom HTTP headers forwarded to the LLM provider |
| `GreetingConfigs` | `map[string]interface{}`   | No       | —                        | Greeting playback configuration |
| `TemplateVariables` | `map[string]string`      | No       | —                        | Template variables for messages |

### Other LLM vendors

The SDK also includes named helpers for the remaining Agora-supported LLM providers. These helpers choose the correct request format internally.

| Constructor | Options Struct | Required Fields |
|---|---|---|
| `NewGroq` | `GroqOptions` | `APIKey` |
| `NewVertexAILLM` | `VertexAILLMOptions` | `APIKey`, `ProjectID`, `Location` |
| `NewAmazonBedrock` | `AmazonBedrockOptions` | `APIKey`, `URL`, `Model` |
| `NewDify` | `DifyOptions` | `APIKey`, `URL` |
| `NewCustomLLM` | `CustomLLMOptions` | `APIKey`, `BaseURL`, `Model` |

---

## TTS Vendors

### NewElevenLabsTTS

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
func NewOpenAITTS(opts OpenAITTSOptions) *OpenAITTS
```

Panics if `Voice` is empty. `APIKey` is optional for the Agora-managed `tts-1` path. Always returns `SampleRate24kHz` from `GetSampleRate()`.

#### OpenAITTSOptions

| Field            | Type       | Required | Description                        |
| ---------------- | ---------- | -------- | ---------------------------------- |
| `APIKey`         | `string`   | No       | OpenAI API key. Optional for the Agora-managed `tts-1` path. |
| `Voice`          | `string`   | Yes      | Voice name                         |
| `Model`          | `string`   | No       | Model identifier                   |
| `ResponseFormat` | `string`   | No       | Audio format (e.g., `"pcm"`)       |
| `Speed`          | `*float64` | No       | Speech speed multiplier            |
| `SkipPatterns`   | `[]int`    | No       | Patterns to skip                   |

### NewCartesiaTTS

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

### NewDeepgramTTS

<!-- snippet: fragment -->
```go
func NewDeepgramTTS(opts DeepgramTTSOptions) *DeepgramTTS
```

Panics if `APIKey` or `Model` is empty.

#### DeepgramTTSOptions

| Field          | Type                     | Required | Description |
| -------------- | ------------------------ | -------- | ----------- |
| `APIKey`       | `string`                 | Yes      | Deepgram API key |
| `Model`        | `string`                 | Yes      | Deepgram TTS model (e.g., `"aura-2-thalia-en"`) |
| `BaseURL`      | `string`                 | No       | WebSocket endpoint; defaults server-side to `wss://api.deepgram.com/v1/speak` |
| `SampleRate`   | `*SampleRate`            | No       | Output sample rate |
| `Params`       | `map[string]interface{}` | No       | Additional Deepgram TTS parameters |
| `SkipPatterns` | `[]int`                  | No       | Patterns to skip |

### NewHumeAITTS

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
func NewMiniMaxTTS(opts MiniMaxTTSOptions) *MiniMaxTTS
```

Panics if `Model` is empty. `Key` is optional for supported Agora-managed MiniMax models (`speech-2.6-turbo`, `speech_2_6_turbo`, `speech-2.8-turbo`, `speech_2_8_turbo`). BYOK still requires `Key` and `GroupID`, and Agora-managed mode must not set `GroupID`, `VoiceID`, or `URL`.

#### MiniMaxTTSOptions

| Field          | Type     | Required | Description                               |
| -------------- | -------- | -------- | ----------------------------------------- |
| `Key`          | `string` | No       | MiniMax API key. Optional for supported Agora-managed MiniMax models. |
| `GroupID`      | `string` | No       | MiniMax group ID. Required for BYOK.      |
| `Model`        | `string` | Yes      | Model name (e.g., `speech-02-turbo`)      |
| `VoiceID`      | `string` | No       | Voice style identifier. BYOK only.        |
| `URL`          | `string` | No       | WebSocket endpoint. BYOK only.            |
| `SkipPatterns` | `[]int`  | No       | Patterns to skip                          |

### NewMurfTTS

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
func NewSpeechmaticsSTT(opts SpeechmaticsSTTOptions) *SpeechmaticsSTT
```

Panics if `APIKey` or `Language` is empty.

#### SpeechmaticsSTTOptions

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | Speechmatics API key |
| `Language` | `string` | Yes | Speechmatics language code |
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |
| `Model` | `string` | No | Model identifier |
| `URI` | `string` | No | Speechmatics streaming WebSocket URL |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

### NewDeepgramSTT

<!-- snippet: fragment -->
```go
func NewDeepgramSTT(opts DeepgramSTTOptions) *DeepgramSTT
```

Does not panic. `APIKey` is optional for Agora-managed Deepgram presets.

#### DeepgramSTTOptions

| Field              | Type                     | Required | Description              |
| ------------------ | ------------------------ | -------- | ------------------------ |
| `APIKey`           | `string`                 | No       | Deepgram API key         |
| `Model`            | `string`                 | No       | Model (e.g., `"nova-2"`) |
| `Language`         | `string`                 | No       | Language code            |
| `InteractionLanguage` | `string`              | No       | Agora `asr.language` override |
| `SmartFormat`      | `*bool`                  | No       | Enable smart formatting  |
| `Punctuation`      | `*bool`                  | No       | Enable punctuation       |
| `AdditionalParams` | `map[string]interface{}` | No       | Additional vendor params |

### NewMicrosoftSTT

<!-- snippet: fragment -->
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
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |

### NewOpenAISTT

<!-- snippet: fragment -->
```go
func NewOpenAISTT(opts OpenAISTTOptions) *OpenAISTT
```

Panics if `APIKey` is empty.

#### OpenAISTTOptions

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | OpenAI API key |
| `Model` | `string` | No | Transcription model |
| `Language` | `string` | No | OpenAI transcription language |
| `Prompt` | `string` | No | Prompt for OpenAI transcription |
| `InputAudioTranscription` | `map[string]interface{}` | No | OpenAI transcription settings |
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

### NewGoogleSTT

<!-- snippet: fragment -->
```go
func NewGoogleSTT(opts GoogleSTTOptions) *GoogleSTT
```

Panics if `ProjectID`, `Location`, or `ADCCredentialsString` is empty.

#### GoogleSTTOptions

| Field | Type | Required | Description |
|---|---|---|---|
| `ProjectID` | `string` | Yes | Google Cloud project ID |
| `Location` | `string` | Yes | Google Cloud region |
| `ADCCredentialsString` | `string` | Yes | Google service account credentials JSON string |
| `Language` | `string` | No | Google recognition language |
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |
| `Model` | `string` | No | Model identifier |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

### NewAmazonSTT

<!-- snippet: fragment -->
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
| `Language`  | `string` | No       | Amazon `language_code` |
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

### NewAssemblyAISTT

<!-- snippet: fragment -->
```go
func NewAssemblyAISTT(opts AssemblyAISTTOptions) *AssemblyAISTT
```

Panics if `APIKey` is empty.

#### AssemblyAISTTOptions

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | AssemblyAI API key |
| `Language` | `string` | No | AssemblyAI language code |
| `InteractionLanguage` | `string` | No | Agora `asr.language` override |
| `URI` | `string` | No | AssemblyAI streaming WebSocket URL |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

### NewAresSTT

<!-- snippet: fragment -->
```go
func NewAresSTT(opts AresSTTOptions) *AresSTT
```

Ares is an Agora built-in STT service — no external API key required.

#### AresSTTOptions

| Field              | Type                     | Required | Description              |
| ------------------ | ------------------------ | -------- | ------------------------ |
| `Language`         | `string`                 | No       | Language code            |
| `AdditionalParams` | `map[string]interface{}` | No       | Additional vendor params |

### NewSarvamSTT

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
func NewOpenAIRealtime(opts OpenAIRealtimeOptions) *OpenAIRealtime
```

Panics if `APIKey` is empty.

#### OpenAIRealtimeOptions

| Field             | Type       | Required | Default                     | Description                                        |
| ----------------- | ---------- | -------- | --------------------------- | -------------------------------------------------- |
| `APIKey`          | `string`   | Yes      | —                           | OpenAI API key                                     |
| `Model`           | `string`   | No       | `"gpt-4o-realtime-preview"` | Model identifier                                   |
| `URL`             | `string`   | No       | —                           | Custom realtime WebSocket URL                      |
| `GreetingMessage` | `string`                   | No       | —                           | Initial greeting                                   |
| `FailureMessage`  | `string`                   | No       | —                           | Fallback message                                   |
| `InputModalities` | `[]string`                 | No       | —                           | Input modalities                                   |
| `OutputModalities` | `[]string`                | No       | —                           | Output modalities                                  |
| `Messages`        | `[]map[string]interface{}` | No       | —                           | Conversation messages for short-term memory        |
| `Params`          | `map[string]interface{}`   | No       | —                           | Additional realtime params such as `voice`         |
| `TurnDetection`   | `*Agora.StartAgentsRequestPropertiesMllmTurnDetection` | No | — | MLLM turn detection configuration; overrides top-level turn detection |

### NewGeminiLive

<!-- snippet: fragment -->
```go
func NewGeminiLive(opts GeminiLiveOptions) *GeminiLive
```

Panics if `APIKey` or `Model` is empty.

#### GeminiLiveOptions

| Field              | Type                       | Required | Default | Description |
| ------------------ | -------------------------- | -------- | ------- | ----------- |
| `APIKey`           | `string`                   | Yes      | —       | Google AI API key |
| `Model`            | `string`                   | Yes      | —       | Gemini Live model identifier |
| `URL`              | `string`                   | No       | —       | Custom realtime WebSocket URL |
| `Instructions`     | `string`                   | No       | —       | System instruction |
| `Voice`            | `string`                   | No       | —       | Voice name |
| `GreetingMessage`  | `string`                   | No       | —       | Initial greeting |
| `FailureMessage`   | `string`                   | No       | —       | Fallback message |
| `InputModalities`  | `[]string`                 | No       | —       | Input modalities |
| `OutputModalities` | `[]string`                 | No       | —       | Output modalities |
| `Messages`         | `[]map[string]interface{}` | No       | —       | Conversation messages |
| `AdditionalParams` | `map[string]interface{}`   | No       | —       | Additional Gemini params |
| `TurnDetection`    | `*Agora.StartAgentsRequestPropertiesMllmTurnDetection` | No | — | MLLM turn detection configuration; overrides top-level turn detection |

### NewXaiGrok

<!-- snippet: fragment -->
```go
func NewXaiGrok(opts XaiGrokOptions) *XaiGrok
```

xAI Grok MLLM vendor (`mllm.vendor`: `"xai"`). Panics if `APIKey` is empty. Defaults `URL` to `wss://api.x.ai/v1/realtime`.

> `NewXAIGrok` / `XAIGrokOptions` are deprecated aliases.

#### XaiGrokOptions

Same fields as `XAIGrokOptions` below.

### NewXAIGrok (deprecated)

<!-- snippet: fragment -->
```go
func NewXAIGrok(opts XAIGrokOptions) *XAIGrok
```

Deprecated. Use `NewXaiGrok` instead.

#### XAIGrokOptions

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `APIKey` | `string` | Yes | — | xAI API key |
| `URL` | `string` | No | `"wss://api.x.ai/v1/realtime"` | Realtime WebSocket URL |
| `Voice` | `string` | No | — | Voice identifier |
| `Language` | `string` | No | — | Language code |
| `SampleRate` | `*int` | No | — | Audio sample rate in Hz |
| `GreetingMessage` | `string` | No | — | Initial greeting |
| `FailureMessage` | `string` | No | — | Fallback message |
| `InputModalities` | `[]string` | No | — | Input modalities |
| `OutputModalities` | `[]string` | No | — | Output modalities |
| `Messages` | `[]map[string]interface{}` | No | — | Conversation messages |
| `Params` | `map[string]interface{}` | No | — | Additional xAI params |
| `TurnDetection` | `*Agora.StartAgentsRequestPropertiesMllmTurnDetection` | No | — | `agora_vad` / `server_vad` turn detection |

### NewVertexAI

<!-- snippet: fragment -->
```go
func NewVertexAI(opts VertexAIOptions) *VertexAI
```

 Panics if `ProjectID` or `ADCredentialsString` is empty.

#### VertexAIOptions

| Field             | Type       | Required | Default                  | Description                                     |
| ----------------- | ---------- | -------- | ------------------------ | ----------------------------------------------- |
| `ProjectID`       | `string`   | Yes      | —                        | GCP project ID                                  |
| `Location`        | `string`   | No       | `"us-central1"`          | GCP region                                      |
| `Model`           | `string`   | No       | `"gemini-2.0-flash-exp"` | Model identifier                                |
| `URL`             | `string`   | No       | —                        | Custom realtime WebSocket URL                   |
| `ADCredentialsString` | `string` | Yes     | —                        | ADC JSON credentials string                     |
| `Voice`           | `string`   | No       | —                        | Voice name                                      |
| `Instructions`    | `string`                   | No       | —                        | System instruction                              |
| `Messages`        | `[]map[string]interface{}` | No       | —                        | Conversation messages                           |
| `GreetingMessage` | `string`                   | No       | —                        | Initial greeting                                |
| `FailureMessage`  | `string`                   | No       | —                        | Fallback message                                |
| `InputModalities` | `[]string`                 | No       | —                        | Input modalities                                |
| `OutputModalities` | `[]string`                | No       | —                        | Output modalities                               |
| `AdditionalParams` | `map[string]interface{}`  | No       | —                        | Additional Vertex/Gemini params                 |
| `TurnDetection`    | `*Agora.StartAgentsRequestPropertiesMllmTurnDetection` | No | — | MLLM turn detection configuration; overrides top-level turn detection |

---

## Avatar Vendors

### NewLiveAvatarAvatar

<!-- snippet: fragment -->
```go
func NewLiveAvatarAvatar(opts LiveAvatarAvatarOptions) *LiveAvatarAvatar
```

Panics if `APIKey` or `AgoraUID` is empty, or if `Quality` is not one of `"low"`, `"medium"`, `"high"`.

Required TTS sample rate: **24kHz** (`SampleRate24kHz`)

#### LiveAvatarAvatarOptions

| Field                 | Type     | Required | Description                                      |
| --------------------- | -------- | -------- | ------------------------------------------------ |
| `APIKey`              | `string` | Yes      | LiveAvatar API key                               |
| `Quality`             | `string` | Yes      | `"low"`, `"medium"`, or `"high"`                 |
| `AgoraUID`            | `string` | Yes      | UID for avatar's video stream                    |
| `AgoraToken`          | `string` | No       | Avatar Agora token. Auto-generated when omitted. |
| `AvatarID`            | `string` | No       | LiveAvatar avatar ID                             |
| `Enable`              | `*bool`  | No       | Enable or disable the avatar (default: `true`)   |
| `DisableIdleTimeout`  | `*bool`  | No       | Disable the idle timeout                         |
| `ActivityIdleTimeout` | `*int`   | No       | Idle timeout in seconds (default: 120)           |

### NewAkoolAvatar

<!-- snippet: fragment -->
```go
func NewAkoolAvatar(opts AkoolAvatarOptions) *AkoolAvatar
```

Panics if `APIKey` is empty.

Required TTS sample rate: **16kHz** (`SampleRate16kHz`)

#### AkoolAvatarOptions

| Field              | Type                     | Required | Description |
| ------------------ | ------------------------ | -------- | ----------- |
| `APIKey`           | `string`                 | Yes      | Akool API key |
| `AvatarID`         | `string`                 | No       | Avatar ID |
| `Enable`           | `*bool`                  | No       | Enable or disable the avatar |
| `AdditionalParams` | `map[string]interface{}` | No       | Additional vendor params |

### NewHeyGenAvatar (deprecated)

<!-- snippet: fragment -->
```go
func NewHeyGenAvatar(opts HeyGenAvatarOptions) *HeyGenAvatar
```

Deprecated alias for the LiveAvatar-compatible avatar shape. Use `NewLiveAvatarAvatar` for new integrations.

### NewAnamAvatar

<!-- snippet: fragment -->
```go
func NewAnamAvatar(opts AnamAvatarOptions) *AnamAvatar
```

Panics if `APIKey` is empty.

### NewGenericAvatar

<!-- snippet: fragment -->
```go
func NewGenericAvatar(opts GenericAvatarOptions) *GenericAvatar
```

Panics if `APIKey`, `APIBaseURL`, `AvatarID`, or `AgoraUID` is empty. `AgoraAppID`, `AgoraChannel`, and `AgoraToken` are optional; AgentKit fills them from the session on `Start()` when omitted.

Generic avatars do not enforce a fixed TTS sample rate. Use the sample rate required by your avatar provider.

#### GenericAvatarOptions

| Field | Type | Required | Description |
|---|---|---|---|
| `APIKey` | `string` | Yes | Generic avatar vendor API key |
| `APIBaseURL` | `string` | Yes | Generic avatar API endpoint |
| `AvatarID` | `string` | Yes | Avatar identifier |
| `AgoraUID` | `string` | Yes | UID for avatar video stream; use a different UID from `AgentUID` |
| `AgoraToken` | `string` | No | Avatar token; auto-generated with the same token format as agent tokens when omitted |
| `AgoraAppID` | `string` | No | Overrides session App ID |
| `AgoraChannel` | `string` | No | Overrides session channel |
| `Enable` | `*bool` | No | Enable or disable the avatar |
| `AdditionalParams` | `map[string]interface{}` | No | Additional vendor params |

---

## Sample Rate Constants

<!-- snippet: fragment -->
```go
const (
    LiveAvatarRequiredSampleRate = SampleRate24kHz
    AkoolRequiredSampleRate  = SampleRate16kHz  // 16000 Hz
)
```
