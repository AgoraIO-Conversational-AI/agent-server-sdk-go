---
sidebar_position: 1
title: Overview
description: Introduction to the Agora Conversational AI Go SDK — architecture, installation, and navigation guide.
---

# Agora Conversational AI Go SDK

The Agora Conversational AI Go SDK lets you build voice-powered AI agents on the [Agora Conversational AI](https://docs.agora.io/en/conversational-ai/overview) platform.

Source: [github.com/AgoraIO/agora-agents-go](https://github.com/AgoraIO/agora-agents-go) — module `github.com/AgoraIO/agora-agents-go/v2`.

## Conversation flows

**Cascading flow** uses ASR -> LLM -> TTS and supports the broadest set of vendor combinations.

**MLLM flow** uses a multimodal model such as OpenAI Realtime or Gemini Live for end-to-end audio.

## Start here

- Start with [Quick Start](./getting-started/quick-start.md). It shows the baseline app-credentials setup and starts a cascading ASR -> LLM -> TTS agent.
- Use [MLLM Flow](./guides/mllm-flow.md) when your agent uses one realtime multimodal model, such as OpenAI Realtime or Gemini Live.
- Use [Cascading Flow](./guides/cascading-flow.md) for more examples of the default ASR -> LLM -> TTS flow, including provider-specific configuration.

## How the SDK is organized

| Layer | Package | Description |
|---|---|---|
| **Generated REST clients** | `client`, `option`, `agents`, `telephony`, `phonenumbers` | Typed access to REST APIs not covered by AgentKit |
| **AgentKit** | `agentkit`, `agentkit/vendors` | Agent builder, session lifecycle, and typed vendors |

## Installation

```sh
go get github.com/AgoraIO/agora-agents-go/v2@v2.0.0
```

Requires Go 1.21 or later.

## Documentation

| Section | What you'll find |
|---|---|
| [Installation](./getting-started/installation.md) | Prerequisites, package install, import paths |
| [Authentication](./getting-started/authentication.md) | App credentials for REST auth and RTC joins |
| [Quick Start](./getting-started/quick-start.md) | App credentials and AgentKit |
| [BYOK](./guides/byok.md) | Bring your own vendor credentials and config |
| [Architecture](./concepts/architecture.md) | SDK structure and generated REST clients |
| [Agent](./concepts/agent.md) | Builder pattern, immutable reuse, vendor configuration |
| [AgentSession](./concepts/session.md) | State machine, lifecycle methods, events |
| [Vendors](./concepts/vendors.md) | LLM, TTS, STT, MLLM, and Avatar provider catalog |
| [Cascading Flow](./guides/cascading-flow.md) | Step-by-step ASR -> LLM -> TTS |
| [MLLM Flow](./guides/mllm-flow.md) | OpenAI Realtime and Vertex AI Gemini Live |
| [Avatars](./guides/avatars.md) | LiveAvatar, Generic Avatar, Anam, Akool, and deprecated HeyGen |
| [Regional Routing](./guides/regional-routing.md) | Area enum, domain pool, failover |
| [Error Handling](./guides/error-handling.md) | API errors and Go error handling patterns |
| [Pagination](./guides/pagination.md) | Iterate over paginated list endpoints |
| [Advanced](./guides/advanced.md) | Headers, retries, timeouts, raw response, custom HTTP client |
| [Low-Level API](./guides/low-level-api.md) | Generated REST APIs |
| [Client Reference](./reference/client.md) | Constructor options, public methods |
| [Agent Reference](./reference/agent.md) | Full builder API |
| [Session Reference](./reference/session.md) | All methods and payload types |
| [Vendor Reference](./reference/vendors.md) | Constructor options for every vendor class |

For generated REST API types, see the [API Reference](../../reference.md).
