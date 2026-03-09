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

	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit/vendors"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
	c := client.NewClient(option.WithArea(option.AreaUS))

	agent := agentkit.NewAgent(
		agentkit.WithName("support-assistant"),
		agentkit.WithInstructions("You are a helpful voice assistant."),
		agentkit.WithGreeting("Hello! How can I help you today?"),
		agentkit.WithMaxHistory(10),
	).
		// Create Agent: STT → LLM → TTS → (optional) Avatar
		WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{APIKey: "your-deepgram-key"})).
		WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{APIKey: "your-openai-key", Model: "gpt-4o-mini"})).
		WithTts(vendors.NewElevenLabsTTS(vendors.ElevenLabsTTSOptions{
			Key: "your-elevenlabs-key", ModelID: "eleven_flash_v2_5", VoiceID: "your-voice-id",
		}))
	// .WithAvatar(vendors.NewHeyGenAvatar(...)) // optional

	expiresIn, _ := agentkit.ExpiresInHours(12) // optional — default is 86400 s (24 h)

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
		ExpiresIn:                expiresIn,
	})

	ctx := context.Background()
	// Start() returns a session ID — a unique identifier for this agent session
	agentSessionId, err := session.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_ = agentSessionId

	// In production, Stop is typically called when your client signals the session has ended.
	// Your server receives that request and calls session.Stop().
	_ = session.Stop(ctx)
}
```

### Session lifecycle

`Start()` joins the agent to the channel and returns a **session ID** — a unique identifier for this agent session. The session stays active until `Stop()` is called.

There are two ways to stop a session depending on how your server is structured:

**Option 1 — Hold the session in memory:**
```go
// start-session handler
agentSessionId, err := session.Start(ctx) // unique ID for this session
// stop-session handler (same process, session still in scope)
_ = session.Stop(ctx)
```

**Option 2 — Store the session ID and stop by ID (stateless servers):**
```go
// start-session handler: return session ID to your client app
agentSessionId, err := session.Start(ctx)
// ... return agentSessionId to client ...

// stop-session handler: client sends back agentSessionId
c := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
    Area: option.AreaUS, AppID: "your-app-id", AppCertificate: "your-app-certificate",
})
if err := c.StopAgent(ctx, agentSessionId); err != nil {
    log.Printf("stop failed: %v", err)
}
```

### Manual tokens (for debugging)

Generate tokens yourself and pass them in — useful when inspecting or reusing tokens:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
	const (
		appID    = "your-app-id"
		appCert  = "your-app-certificate"
		channel  = "support-room-123"
		agentUID = "1"
	)

	expiry, _ := agentkit.ExpiresInHours(12)

	// Auth header token — used by the SDK to authenticate REST API calls
	authToken, err := agentkit.GenerateConvoAIToken(agentkit.GenerateConvoAITokenOptions{
		AppID: appID, AppCertificate: appCert,
		ChannelName: channel, Account: agentUID,
		TokenExpire: expiry,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Channel join token — embedded in the start request so the agent can join the channel
	joinToken, err := agentkit.GenerateConvoAIToken(agentkit.GenerateConvoAITokenOptions{
		AppID: appID, AppCertificate: appCert,
		ChannelName: channel, Account: agentUID,
		TokenExpire: expiry,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Pass the auth token to the REST client using option.WithToken
	c := client.NewClient(
		option.WithArea(option.AreaUS),
		option.WithToken(authToken), // sets Authorization: agora token=<authToken>
	)

	session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
		Client:     c.Agents,
		Agent:      agent,
		AppID:      appID,
		Channel:    channel,
		AgentUID:   agentUID,
		RemoteUIDs: []string{"100"},
		Token:      joinToken, // channel join token
	})

	agentSessionId, err := session.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_ = agentSessionId
}
```

## Documentation

| Topic                | Link                                                                               |
| -------------------- | ---------------------------------------------------------------------------------- |
| **API docs**         | [docs.agora.io](https://docs.agora.io/en/conversational-ai/overview)               |
| **Installation**     | [docs/getting-started/installation.md](./docs/getting-started/installation.md)     |
| **Authentication**   | [docs/getting-started/authentication.md](./docs/getting-started/authentication.md) |
| **Quick Start**      | [docs/getting-started/quick-start.md](./docs/getting-started/quick-start.md)       |
| **Cascading flow**   | [docs/guides/cascading-flow.md](./docs/guides/cascading-flow.md)                   |
| **MLLM flow**        | [docs/guides/mllm-flow.md](./docs/guides/mllm-flow.md)                             |
| **Avatars**          | [docs/guides/avatars.md](./docs/guides/avatars.md)                                 |
| **Regional routing** | [docs/guides/regional-routing.md](./docs/guides/regional-routing.md)               |
| **Error handling**   | [docs/guides/error-handling.md](./docs/guides/error-handling.md)                   |
| **Pagination**       | [docs/guides/pagination.md](./docs/guides/pagination.md)                           |
| **Advanced**         | [docs/guides/advanced.md](./docs/guides/advanced.md)                               |
| **Low-level API**    | [docs/guides/low-level-api.md](./docs/guides/low-level-api.md)                     |
| **API reference**    | [reference.md](./reference.md)                                                     |

## Contributing

This library is generated programmatically. Contributions to the README and docs are welcome. For code changes, open an issue first to discuss.
