# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [v1.2.0] — 2026-03-27

### Fixed

- **`OpenAIRealtime` / `VertexAI` (MLLM)** — Agent-level `greeting`, `failure_message`, and `max_history` overrides are now correctly applied when the agent is in MLLM mode. Previously these values were silently dropped.
- **`VertexAI` (MLLM)** — Corrected vendor string from `"vertex_ai"` to `"vertexai"` to match the Agora server API.

### Changed

- **`OpenAITTS`** — Renamed struct field `Key` → `APIKey` to match the Agora server API expectation. ⚠️ **Breaking change.**
- **`CartesiaTTS`** — Renamed struct field `Key` → `APIKey`. Voice is now serialized as `{"mode": "id", "id": "<VoiceID>"}` instead of a flat `voice_id` string. ⚠️ **Breaking change.**
- **`HeyGenAvatar`** — Removed legacy fields `AvatarName`, `VoiceID`, `Language`, `Version`. Added `AgoraToken`, `AvatarID`, `Enable`, `DisableIdleTimeout`, `ActivityIdleTimeout`. The config now includes a top-level `enable` field (defaults `true`). ⚠️ **Breaking change.**

### Added

- **`OpenAITTS`** — New optional fields: `ResponseFormat` (string, e.g. `"pcm"`) and `Speed` (*float64).
- **`CartesiaTTS`** — `VoiceID` user-facing field is preserved; voice is serialized to the required nested object format automatically.
- **`RimeTTS`** — New optional fields: `Lang` (string), `SamplingRate` (*int, serialized as `samplingRate`), `SpeedAlpha` (*float64, serialized as `speedAlpha`).
- **`OpenAIRealtime`** — New optional fields: `PredefinedTools` ([]string), `FailureMessage` (string), `MaxHistory` (*int).
- **`VertexAI` (MLLM)** — New optional fields: `PredefinedTools` ([]string), `FailureMessage` (string), `MaxHistory` (*int).
- **`HeyGenAvatar`** — New fields: `AgoraToken` (string), `AvatarID` (string), `Enable` (*bool, default `true`), `DisableIdleTimeout` (*bool), `ActivityIdleTimeout` (*int).

## [v1.1.0] — 2026-03-17

### Added

- `MurfTTS` vendor
- `AdditionalParams map[string]interface{}` on all STT vendors for passing unlisted API parameters


### Fixed

- `ElevenLabsTTS`: added missing voice tuning fields — `OptimizeStreamingLatency`, `Stability`, `SimilarityBoost`, `Style`, `UseSpeakerBoost`
- All LLM vendors: added `MaxHistory *int` for conversation history caching
- `AzureOpenAI` LLM: added `Params map[string]interface{}` escape hatch (was missing, unlike other vendors)
- `Anthropic` LLM: added `URL string` for custom endpoints and `Params map[string]interface{}` escape hatch
- `Gemini` LLM: added `URL string` for custom endpoints and `Params map[string]interface{}` escape hatch
- `MiniMaxTTS`: added required `GroupID`, `URL`, and correctly nested `voice_setting.voice_id`
- `SarvamTTS`: corrected schema to `Key` + `Speaker` + `TargetLanguageCode` (was incorrectly using `APIKey`, `VoiceID`, `Model`)
- All STT vendors: added top-level `language` field to `ToConfig()` output, required by the Agora platform for routing
- `GoogleSTT`: renamed `Key` → `APIKey`; corrected payload key from `"key"` to `"api_key"`
- `AresSTT`: removed erroneous `APIKey` requirement — Ares is an Agora built-in service with no external key
- `AssemblyAISTT`: added missing `Language` field
- `DeepgramSTT`: added `SmartFormat` and `Punctuation` fields; made `APIKey` optional to match other SDKs
- `SarvamSTT`: added `Language` validation

## [v1.0.0] — 2026-03-11

Initial stable release of the Agora Agent Server SDK for Go.

### Added

- `Agent` builder with functional options (`WithLlm()`, `WithTts()`, `WithStt()`, `WithMllm()`, `WithAvatar()`)
- `AgentSession` for session lifecycle management (`Start()`, `Stop()`)
- Automatic token generation — pass `AppID` + `AppCertificate` and tokens are handled internally
- Token utilities: `ExpiresInHours()`, `ExpiresInMinutes()`
- Turn detection configuration via `TurnDetectionConfig` with nested `StartOfSpeechConfig` and `EndOfSpeechConfig`
- SAL (Selective Attention Locking) via `SalConfig` with `SalMode`
- Filler words support: `FillerWordsConfig`, `FillerWordsTrigger`, `FillerWordsContent`
- Session parameters: `SessionParams`, `SilenceConfig`, `FarewellConfig`, `ParametersDataChannel`
- Geofencing via `GeofenceConfig`
- Advanced features (MLLM mode) via `AdvancedFeatures`
- Vendor integrations:
  - **LLM**: `OpenAI`, `AzureOpenAI`, `Anthropic`, `Gemini`, `VertexAI`
  - **MLLM**: `OpenAIRealtime`
  - **TTS**: `ElevenLabsTTS`, `MicrosoftTTS`, `OpenAITTS`, `CartesiaTTS`, `GoogleTTS`, `AmazonTTS`, `HumeAITTS`, `RimeTTS`, `FishAudioTTS`, `MiniMaxTTS`, `SarvamTTS`
  - **STT**: `DeepgramSTT`, `MicrosoftSTT`, `OpenAISTT`, `GoogleSTT`, `AmazonSTT`, `AssemblyAISTT`, `AresSTT`, `SarvamSTT`, `SpeechmaticsSTT`
  - **Avatar**: `HeyGenAvatar`, `AkoolAvatar`
