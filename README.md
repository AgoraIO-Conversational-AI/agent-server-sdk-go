# Agora Conversational AI Go SDK

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=https%3A%2F%2Fgithub.com%2FAgoraIO-Conversational-AI%2Fagora-agent-go-sdk)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk)](https://pkg.go.dev/github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk)

The Agora Conversational AI SDK for Go provides convenient access to the Agora Conversational AI APIs, enabling you to build voice-powered AI agents with support for both **cascading flows** (ASR → LLM → TTS) and **multimodal flows** (MLLM) for real-time audio processing.

## Installation

```sh
go get github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk
```

## Quick Start

Use the **builder pattern** with `Agent` and `AgentSession`:

```go
package main

import (
	"context"
	"log"

	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit/vendors"
)

func main() {
	c := client.NewClient(
		option.WithArea(option.AreaUS),
	)

	agent := agentkit.NewAgent(
		agentkit.WithName("support-assistant"),
		agentkit.WithInstructions("You are a helpful voice assistant."),
		agentkit.WithGreeting("Hello! How can I help you today?"),
		agentkit.WithMaxHistory(10),
	).
		WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
			APIKey: "your-openai-key",
			Model:  "gpt-4o-mini",
		})).
		WithTts(vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
			Key:     "your-elevenlabs-key",
			ModelID: "eleven_flash_v2_5",
			VoiceID: "your-voice-id",
		})).
		WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
			APIKey: "your-deepgram-key",
		}))

	session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
		Client:                   c.Agents,
		Agent:                    agent,
		AppID:                    "your-app-id",
		AppCertificate:           "your-app-certificate",
		Name:                     "support-assistant",
		Channel:                  "support-room-123",
		AgentUID:                 "1",
		RemoteUIDs:               []string{"100"},
		UseAppCredentialsForREST: true,
	})

	ctx := context.Background()
	agentID, err := session.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_ = agentID

	_ = session.Say(ctx, "Hello! How can I help you today?", nil, nil)
	_ = session.Stop(ctx)
}
```

## Documentation

| Topic              | Link                                                                               |
| ------------------ | ---------------------------------------------------------------------------------- |
| **API docs**       | [docs.agora.io](https://docs.agora.io/en/conversational-ai/overview)               |
| **Installation**   | [docs/getting-started/installation.md](./docs/getting-started/installation.md)     |
| **Authentication** | [docs/getting-started/authentication.md](./docs/getting-started/authentication.md) |
| **Quick Start**    | [docs/getting-started/quick-start.md](./docs/getting-started/quick-start.md)       |
| **Cascading flow** | [docs/guides/cascading-flow.md](./docs/guides/cascading-flow.md)                   |
| **MLLM flow**      | [docs/guides/mllm-flow.md](./docs/guides/mllm-flow.md)                             |
| **Avatars**        | [docs/guides/avatars.md](./docs/guides/avatars.md)                                 |
| **Regional routing** | [docs/guides/regional-routing.md](./docs/guides/regional-routing.md)               |
| **Error handling** | [docs/guides/error-handling.md](./docs/guides/error-handling.md)                     |
| **Pagination**     | [docs/guides/pagination.md](./docs/guides/pagination.md)                           |
| **Advanced**       | [docs/guides/advanced.md](./docs/guides/advanced.md)                               |
| **Low-level API**  | [docs/guides/low-level-api.md](./docs/guides/low-level-api.md)                     |
| **API reference**  | [reference.md](./reference.md)                                                     |

## Contributing

This library is generated programmatically. Contributions to the README and docs are welcome. For code changes, open an issue first to discuss.
