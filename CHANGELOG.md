# Changelog

## Unreleased

### Added
- `Agent.CreateSession(...)` for ergonomic session creation from an `Agent`.
- Session preset support, including preset normalization and preset / `pipeline_id` forwarding at start time.
- `AgentSession.GetTurns()` for turn analytics without dropping to the raw client.
- `AgentSession.Off(...)` for unregistering event handlers.
- `AgentSession` `Debug` and `Warn` options for session-level logging and warning hooks.
- `vendors.NewGeminiLive(...)`, `vendors.NewLiveAvatarAvatar(...)`, and `vendors.NewAnamAvatar(...)`.
- `agentkit/presets.go` helper utilities for preset inference and normalization.

### Changed
- Aligned Go AgentKit session, avatar, and MLLM wrappers with the conservative low-level SDK contract.
- Removed unsupported MLLM wrapper-only fields and kept wrapper behavior limited to generated API-backed fields.
- Corrected Gemini and Vertex MLLM `messages` placement to match the generated contract.
- Updated Go AgentKit docs and examples to reflect the supported session and MLLM APIs.

### Fixed
- Preset-backed sessions now work without requiring explicit LLM/TTS configuration when the low-level contract allows it.
- Avatar validation and sample-rate warnings now match the current wrapper behavior.
- Added focused tests for preset flow, `GetTurns`, event unregistration, avatar support, and MLLM shape enforcement.
