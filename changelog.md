# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [v1.4.0] — 2026-05-13

### Added

- **`NewDeepgramTTS`** — New TTS vendor wrapper for Deepgram (Beta). Accepts `APIKey`, `Model`, `BaseURL`, `SampleRate`, `Params`, and `SkipPatterns`.
- **`Agent.WithTools(enabled bool)` / `WithTools` option** — Dedicated builder method and constructor option to enable MCP tool invocation (`advanced_features.enable_tools`). Replaces the raw `WithAdvancedFeatures(&AdvancedFeatures{EnableTools: Agora.Bool(true)})` call.
- **LLM vendors: `Headers` field** — All four LLM vendors (`OpenAI`, `AzureOpenAI`, `Anthropic`, `Gemini`) now accept an optional `Headers map[string]string` field. Use this to pass custom HTTP headers to the LLM provider (e.g., tenant identifiers, routing headers).
- **`AgentSession.Think()` / `ThinkWithOptions()`** — Send a custom instruction to a running agent through the agent management API.
- **`Agent.WithInterruption()` / `WithInterruptionConfig()`** — Configure the new top-level `interruption` object for unified interruption control.
- **MLLM turn detection** — `NewOpenAIRealtime`, `NewGeminiLive`, and `NewVertexAI` now accept `TurnDetection`, which maps to `mllm.turn_detection` and overrides top-level turn detection for MLLM sessions.
- **`AudioScenario` AgentKit support** — Session params and AgentKit request construction now expose the top-level `parameters.audio_scenario` field.

### Fixed

- **MiniMax TTS preset stripping** — When a MiniMax reseller preset is inferred (`minimax_speech_2_6_turbo` or `minimax_speech_2_8_turbo`), the `group_id` and `url` fields are now correctly stripped from `tts.params` alongside `key` and `model`. Previously they were forwarded to the API, causing request failures.
- **MLLM enable flag** — `Agent.WithMllm()` now sets `mllm.enable = true` and removes the deprecated `advanced_features.enable_mllm` flag from generated requests.
- **MLLM wrapper shape** — MLLM vendors no longer emit removed fields such as `style`; docs and tests now reflect the v2.6 MLLM contract.
- **Preset-backed OpenAI TTS** — `NewOpenAITTS` no longer requires `APIKey` when a reseller preset supplies credentials server-side.
- **AgentKit parity coverage** — Added regression coverage for interruption, MLLM turn detection, Deepgram TTS, LLM headers, and deprecated MLLM flag cleanup.

## [v1.3.0] — 2026-04-02

### Added

- **`AgentSession.GetTurns()`** — Retrieve conversation turn analytics for a running session.
- **Session-level `Preset` and `PipelineID`** — `CreateSessionOptions` now accepts `Preset` and `PipelineID`. Automatic reseller preset inference for supported Deepgram (nova-2, nova-3), OpenAI (gpt-4o-mini, gpt-4.1-mini, gpt-5-nano, gpt-5-mini), and MiniMax (speech-2.6-turbo, speech-2.8-turbo) models.
- **`AgentPresets` constants** — Type-safe preset constant map for discoverable session preset composition.
- **`LiveAvatarAvatar` and `AnamAvatar`** — New avatar vendor wrappers. `LiveAvatarAvatar` is the successor to `HeyGenAvatar` and emits `vendor: "liveavatar"`.

### Changed

- **`GeminiLive` MLLM wrapper** — Output now matches the Agora low-level MLLM contract; `messages` stays at top level.
- **Avatar sample-rate validation** — `CreateSession` validates TTS sample rate against avatar vendor requirements and emits a warning (via `WarnHook`) when using `HeyGenAvatar` or `LiveAvatarAvatar` without an explicit 24 kHz TTS sample rate.

### Fixed

- **MLLM wrapper fields** — Removed unsupported wrapper-only fields that caused the Go surface to diverge from the generated Agora API contract.

## [v1.2.0] — 2026-03-27

### Fixed

- **`AresSTT`** — Removed redundant `language` key from the `Params` map. Language is now emitted only at the top level.
- **`OpenAIRealtime` / `GeminiLive` (MLLM)** — Agent-level `greeting`, `FailureMessage`, and `MaxHistory` overrides are now correctly applied in MLLM mode.
- **`GeminiLive` (MLLM)** — `messages` is now correctly placed inside `params` as required by the Gemini Live API.

### Changed

- **`OpenAITTS`** — Constructor field renamed `Key` → `APIKey`. ⚠️ **Breaking change.**
- **`CartesiaTTS`** — Constructor field renamed `Key` → `APIKey`. Voice is now serialized as `{"mode": "id", "id": "<voice_id>"}`. ⚠️ **Breaking change.**
- **`HeyGenAvatar`** — Removed legacy fields `AvatarName`, `VoiceID`, `Language`, `Version`. Added `AgoraToken`, `AvatarID`, `Enable`, `DisableIdleTimeout`, `ActivityIdleTimeout`. Config now includes a top-level `enable` field (defaults `true`). ⚠️ **Breaking change.**

### Added

- **`OpenAITTS`** — New optional fields: `ResponseFormat` (string) and `Speed` (float64).
- **`CartesiaTTS`** — `VoiceID` user-facing field is preserved; serialized to the required nested object format automatically.
- **`RimeTTS`** — New optional fields: `Lang` (string), `SamplingRate` (int), `SpeedAlpha` (float64).
- **`OpenAIRealtime`** — New optional fields: `PredefinedTools` ([]string), `FailureMessage` (string), `MaxHistory` (int).
- **`GeminiLive` (MLLM)** — New optional fields: `PredefinedTools` ([]string), `FailureMessage` (string), `MaxHistory` (int).

## [v1.1.0] — 2026-03-17

### Added

- **`NewMurfTTS`** — New TTS vendor wrapper for Murf.
- **`NewSarvamTTS`** — New TTS vendor wrapper for Sarvam.
- **`NewSarvamSTT`** — New STT vendor wrapper for Sarvam.
- **`AgentSession.Think()` / `ThinkWithOptions()`** — Send a custom instruction to the agent mid-session. Routes to the `agentManagement` client.
- **`ThinkOptions`** — Idiomatic Go struct for `Think` call options (`OnListeningAction`, `OnThinkingAction`, `OnSpeakingAction`, `Interruptable`, `Metadata`).
- All LLM vendors: added `MaxHistory` field for conversation history caching.
- **`AzureOpenAI`** — Added `Params` escape hatch for arbitrary API parameters.
- **`Anthropic`** — Added `URL` for custom endpoints and `Params` escape hatch.
- **`Gemini`** — Added `URL` for custom endpoints and `Params` escape hatch.
- **`SpeechmaticsSTT`**, **`SarvamSTT`** — Added optional `Model` field.

### Fixed

- **`MiniMaxTTS`** — Added required `GroupID`, `URL`, and correctly nested `VoiceSetting.VoiceID` fields.
- **`SarvamTTS`** — Corrected schema to `Key` + `Speaker` + `TargetLanguageCode`.

## [v1.0.0] — 2026-03-11

Initial stable release of Agora Agents Go.

### Added

- `Agent` builder with immutable fluent API (`WithLlm()`, `WithTts()`, `WithStt()`, `WithMllm()`, `WithAvatar()`) and `NewAgent(opts...)` constructor.
- `AgentSession` for session lifecycle management (`Start()`, `Stop()`, `Say()`, `Interrupt()`, `Update()`, `GetHistory()`, `GetInfo()`).
- App-credentials auth — pass `AppID` + `AppCertificate` to `AgoraClientOptions` and ConvoAI tokens are generated per request.
- Token utilities: `GenerateRtcToken`, `GenerateConvoAIToken`.
- Turn detection configuration via `TurnDetectionConfig` with full `StartOfSpeechConfig` / `EndOfSpeechConfig` nested types and constants.
- SAL (Selective Attention Locking) via `SalConfig` with `SalModeLocking` / `SalModeRecognition` constants.
- Filler words support via `FillerWordsConfig`.
- Session parameters: `SessionParams`, `SilenceConfig`, `FarewellConfig`.
- Geofencing via `GeofenceConfig` with `GeofenceArea` constants.
- Advanced features (MLLM mode) via `AdvancedFeatures`.
- `WarnHook` callback for runtime warnings (avatar sample rate mismatches, handler panics).
- Vendor integrations:
  - **LLM**: `NewOpenAI`, `NewAzureOpenAI`, `NewAnthropic`, `NewGemini`
  - **MLLM**: `NewOpenAIRealtime`, `NewGeminiLive`
  - **TTS**: `NewElevenLabsTTS`, `NewMicrosoftTTS`, `NewOpenAITTS`, `NewCartesiaTTS`, `NewGoogleTTS`, `NewAmazonTTS`, `NewHumeAITTS`, `NewRimeTTS`, `NewFishAudioTTS`, `NewMiniMaxTTS`
  - **STT**: `NewDeepgramSTT`, `NewMicrosoftSTT`, `NewOpenAISTT`, `NewGoogleSTT`, `NewAmazonSTT`, `NewAssemblyAISTT`, `NewAresSTT`, `NewSpeechmaticsSTT`
  - **Avatar**: `NewHeyGenAvatar`, `NewAkoolAvatar`
