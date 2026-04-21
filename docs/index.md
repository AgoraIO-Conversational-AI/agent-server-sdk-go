---
sidebar_position: 1
title: Overview
description: Introduction to the Agora Conversational AI Go SDK — architecture, installation, and navigation guide.
---

# Agora Conversational AI Go SDK

The Agora Conversational AI Go SDK lets you build voice-powered AI agents on the [Agora Conversational AI](https://docs.agora.io/en/conversational-ai/overview) platform.

## Conversation flows

**Cascading flow** uses ASR -> LLM -> TTS and supports the broadest set of vendor combinations.

**MLLM flow** uses a multimodal model such as OpenAI Realtime or Gemini Live for end-to-end audio.

## Choose a starting point

- Use [Quick Start](./getting-started/quick-start.md) if you want the recommended preset-based path with no vendor keys.
- Use [MLLM Flow](./guides/mllm-flow.md) if you want realtime end-to-end audio with OpenAI Realtime or Gemini Live.
- Use [Cascading Flow](./guides/cascading-flow.md) if you want separate ASR, LLM, and TTS vendors.

## SDK layers

| Layer | Package | Description |
|---|---|---|
| **Fern-generated client** | `client`, `option`, `agents`, `telephony`, `phonenumbers` | Low-level typed REST API access |
| **Hand-written agentkit layer** | `agentkit`, `agentkit/vendors` | High-level builder pattern, presets, lifecycle, typed vendors |

## Installation

```sh
go get github.com/AgoraIO-Conversational-AI/agent-server-sdk-go
```

Requires Go 1.21 or later.

## Documentation

| Section | What you'll find |
|---|---|
| [Installation](./getting-started/installation.md) | Prerequisites, package install, import paths |
| [Authentication](./getting-started/authentication.md) | Token auth for REST and RTC joins |
| [Quick Start](./getting-started/quick-start.md) | Recommended preset-based onboarding flow |
| [BYOK](./guides/byok.md) | Bring your own vendor credentials and config |
| [Architecture](./concepts/architecture.md) | Two-layer design, when to use agentkit vs. raw client |
| [Agent](./concepts/agent.md) | Builder pattern, immutable reuse, vendor configuration |
| [AgentSession](./concepts/session.md) | State machine, lifecycle methods, events |
| [Vendors](./concepts/vendors.md) | LLM, TTS, STT, MLLM, and Avatar provider catalog |
| [Cascading Flow](./guides/cascading-flow.md) | Step-by-step ASR -> LLM -> TTS |
| [MLLM Flow](./guides/mllm-flow.md) | OpenAI Realtime and Vertex AI Gemini Live |
| [Avatars](./guides/avatars.md) | HeyGen and Akool with sample-rate constraints |
| [Regional Routing](./guides/regional-routing.md) | Area enum, domain pool, failover |
| [Error Handling](./guides/error-handling.md) | API errors and Go error handling patterns |
| [Pagination](./guides/pagination.md) | Iterate over paginated list endpoints |
| [Advanced](./guides/advanced.md) | Headers, retries, timeouts, raw response, custom HTTP client |
| [Low-Level API](./guides/low-level-api.md) | Direct `client.Agents.Start()` usage |
| [Client Reference](./reference/client.md) | Constructor options, public methods |
| [Agent Reference](./reference/agent.md) | Full builder API |
| [Session Reference](./reference/session.md) | All methods and payload types |
| [Vendor Reference](./reference/vendors.md) | Constructor options for every vendor class |

For Fern-generated raw API types, see the [API Reference](../../reference.md).
