---
sidebar_position: 1
title: Overview
description: Introduction to the Agora Conversational AI Go SDK — architecture, installation, and navigation guide.
---

# Agora Conversational AI Go SDK

The Agora Conversational AI Go SDK enables you to build voice-powered AI agents that interact with users in real time. It supports two conversation flows:

- **Cascading Flow (ASR → LLM → TTS):** Speech-to-text, language model processing, then text-to-speech — three separate vendor services chained together.
- **MLLM Flow (Multimodal):** A single multimodal model (e.g., OpenAI Realtime, Google Gemini Live) handles audio input and output directly.

## Two-Layer Architecture

The SDK has two distinct layers:

| Layer | Package | Description |
|-------|---------|-------------|
| **Fern-generated client** | `client`, `option`, `agents`, `telephony`, `phonenumbers` | Auto-generated from the Agora OpenAPI spec. Provides typed request/response structs and low-level API access. |
| **Hand-written agentkit layer** | `agentkit`, `agentkit/vendors` | Ergonomic builder layer with functional options, vendor constructors, session lifecycle management, and automatic token generation. |

Most applications should use the `agentkit` layer. Use the raw `client` layer when you need direct control over request construction or access to endpoints the agentkit package doesn't cover (e.g., Telephony, PhoneNumbers).

## Installation

```sh
go get github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk
```

Requires Go 1.21 or later.

## Import Paths

```go
import (
    Agora "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk"  // Root types, pointer helpers, environments
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client" // Fern-generated API client
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option" // Request options (auth, base URL, retries)
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit" // Agent builder and session management
    "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit/vendors" // Vendor constructors (LLM, TTS, STT, MLLM, Avatar)
)
```

## Documentation

| Section | What You'll Find |
|---------|------------------|
| [Installation](./getting-started/installation.md) | Prerequisites, package install, import paths |
| [Authentication](./getting-started/authentication.md) | App credentials, pre-built tokens, Basic Auth |
| [Quick Start](./getting-started/quick-start.md) | End-to-end cascading flow example |
| [Architecture](./concepts/architecture.md) | Two-layer design, when to use agentkit vs. raw client |
| [Agent](./concepts/agent.md) | Builder pattern, immutable reuse, vendor configuration |
| [AgentSession](./concepts/session.md) | State machine, lifecycle methods, events |
| [Vendors](./concepts/vendors.md) | LLM, TTS, STT, MLLM, and Avatar provider catalog |
| [Cascading Flow](./guides/cascading-flow.md) | Step-by-step ASR → LLM → TTS |
| [MLLM Flow](./guides/mllm-flow.md) | OpenAI Realtime and Vertex AI Gemini Live |
| [Avatars](./guides/avatars.md) | HeyGen and Akool with sample-rate constraints |
| [Regional Routing](./guides/regional-routing.md) | Area enum, domain pool, failover |
| [Error Handling](./guides/error-handling.md) | APIError and API error handling |
| [Pagination](./guides/pagination.md) | Iterate over paginated list endpoints |
| [Advanced](./guides/advanced.md) | Headers, retries, timeouts, raw response, custom HTTP client |
| [Low-Level API](./guides/low-level-api.md) | Direct client.Agents.Start() without builder |
| [Client Reference](./reference/client.md) | Constructor options, public methods |
| [Agent Reference](./reference/agent.md) | Full builder API |
| [Session Reference](./reference/session.md) | All methods, events, and payload types |
| [Vendor Reference](./reference/vendors.md) | Constructor options for every vendor class |

For Fern-generated raw API types, see the [API Reference](../../reference.md).
