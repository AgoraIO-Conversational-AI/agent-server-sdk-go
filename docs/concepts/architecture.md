---
sidebar_position: 1
title: Architecture
description: Two-layer architecture of the Agora Go SDK — Fern-generated client vs. hand-written agent layer.
---

# Architecture

The SDK is organized into two layers with distinct responsibilities.

## Layer Diagram

```
┌─────────────────────────────────────────────────────┐
│                  Your Application                    │
├─────────────────────────────────────────────────────┤
│  agentkit/         │  agentkit/vendors/              │
│  ─────────         │  ────────────────               │
│  NewAgent()        │  NewOpenAI()                    │
│  CreateSession()   │  NewElevenLabsTTS()             │
│  AgentOption funcs │  NewDeepgramSTT()               │
│  Session lifecycle │  NewOpenAIRealtime()            │
│  Token generation  │  NewHeyGenAvatar()              │
│                    │  ... (30+ vendor constructors)  │
├─────────────────────────────────────────────────────┤
│  client/           │  option/                        │
│  ─────────         │  ────────                       │
│  NewClient()       │  WithBasicAuth()                │
│  client.Agents     │  WithToken()                    │
│  client.Telephony  │  WithBaseURL()                  │
│  client.PhoneNums  │  WithArea()                     │
│                    │  WithMaxAttempts()               │
├─────────────────────────────────────────────────────┤
│  Agora (root)      │  core/                          │
│  ──────────        │  ──────                         │
│  Request/Response  │  HTTP caller, retries,          │
│  types, pointer    │  Pool, Area, DNS resolver       │
│  helpers, enums    │                                 │
└─────────────────────────────────────────────────────┘
```

## Fern-Generated Layer (`client/`, `option/`, root types)

Auto-generated from the Agora OpenAPI specification using [Fern](https://buildwithfern.com). This layer provides:

- **Typed request/response structs** — `Agora.StartAgentsRequest`, `Agora.StartAgentsRequestProperties`, etc.
- **Sub-clients** — `client.Agents`, `client.Telephony`, `client.PhoneNumbers`
- **Request options** — authentication, retries, timeouts, custom headers
- **Pointer helpers** — `Agora.String()`, `Agora.Bool()`, `Agora.Int()`, `Agora.Float64()` for optional fields
- **Automatic retries** with exponential backoff

All API methods on sub-clients take `context.Context` as the first argument:

<!-- snippet: fragment -->
```go
resp, err := c.Agents.Start(ctx, &Agora.StartAgentsRequest{...})
```

## Hand-Written Agentkit Layer (`agentkit/`, `agentkit/vendors/`)

Built on top of the Fern-generated client. This layer provides:

- **`agentkit.NewAgent`** — functional options pattern for building agent configurations
- **Vendor constructors** — `vendors.NewOpenAI()`, `vendors.NewElevenLabsTTS()`, etc. with validation
- **`agent.CreateSession`** — session lifecycle management (start, stop, say, interrupt, update)
- **Automatic token generation** — generates RTC tokens from app credentials
- **Sample rate validation** — prevents avatar/TTS sample rate mismatches at build time

## When to Use Which Layer

| Use Case | Layer | Example |
|---|---|---|
| Build a conversational agent | `agentkit` | `NewAgent` -> `WithLlm` -> `WithTts` -> `CreateSession` -> `Start` |
| Make a telephony call | `client` | `c.Telephony.Call(ctx, req)` |
| Manage phone numbers | `client` | `c.PhoneNumbers.List(ctx, req)` |
| Custom request construction | `client` | `c.Agents.Start(ctx, req)` with manually built properties |
| Regional routing with DNS selection | `option` | `option.WithArea(option.AreaUS)` |

## Context Propagation

All API calls in the Go SDK take `context.Context` as their first argument. This is idiomatic Go and enables:

- **Timeouts:** `context.WithTimeout(ctx, 5*time.Second)`
- **Cancellation:** `context.WithCancel(ctx)` for graceful shutdown
- **Deadline propagation:** Pass request-scoped deadlines through the call chain

<!-- snippet: fragment -->
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

agentID, err := session.Start(ctx)
if err != nil {
    log.Fatalf("Start failed: %v", err)
}
```
