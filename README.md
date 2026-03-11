# Agoraio Go Library

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=https%3A%2F%2Fgithub.com%2FAgoraIO-Conversational-AI%2Fagent-server-sdk-go)

The Agora Conversational AI SDK provides convenient access to the Agora Conversational AI APIs, 
enabling you to build voice-powered AI agents with support for both cascading flows (ASR -> LLM -> TTS) 
and multimodal flows (MLLM) for real-time audio processing.


## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Reference](#reference)
- [Mllm Flow Multimodal](#mllm-flow-multimodal)
- [Mllm Flow Multimodal](#mllm-flow-multimodal)
- [Usage](#usage)
- [Environments](#environments)
- [Errors](#errors)
- [Request Options](#request-options)
- [Advanced](#advanced)
  - [Response Headers](#response-headers)
  - [Retries](#retries)
  - [Timeouts](#timeouts)
  - [Explicit Null](#explicit-null)
- [Contributing](#contributing)

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

API reference documentation is available [here](https://docs.agora.io/en/conversational-ai/overview).

## Reference

A full reference for this library is available [here](https://github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/blob/HEAD/./reference.md).

## MLLM Flow (Multimodal)

For real-time audio processing using OpenAI's Realtime API or Google Gemini Live, use the MLLM (Multimodal Large Language Model) flow instead of the cascading ASR -> LLM -> TTS flow. See the [MLLM Overview](https://docs.agora.io/en/conversational-ai/models/mllm/overview) for more details.

```go
package main

import (
    "context"
    client "github.com/{{ owner }}/{{ repo }}/client"
    option "github.com/{{ owner }}/{{ repo }}/option"
    Agora "github.com/{{ owner }}/{{ repo }}"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customerId>", "<customerSecret>"),
    )

    c.Agents.Start(
        context.TODO(),
        &Agora.StartAgentsRequest{
            Appid: "your_app_id",
            Name:  "mllm_agent",
            Properties: &Agora.StartAgentsRequestProperties{
                Channel:       "channel_name",
                Token:         "your_token",
                AgentRtcUID:   "1001",
                RemoteRtcUIDs: []string{"1002"},
                IdleTimeout:   Agora.Int(120),
                AdvancedFeatures: &Agora.StartAgentsRequestPropertiesAdvancedFeatures{
                    EnableMllm: Agora.Bool(true),
                },
                Mllm: &Agora.StartAgentsRequestPropertiesMllm{
                    URL:    Agora.String("wss://api.openai.com/v1/realtime"),
                    APIKey: Agora.String("<your_openai_api_key>"),
                    Vendor: Agora.StartAgentsRequestPropertiesMllmVendorOpenai,
                    Params: map[string]any{
                        "model": "gpt-4o-realtime-preview",
                        "voice": "alloy",
                    },
                    InputModalities:  []string{"audio"},
                    OutputModalities: []string{"text", "audio"},
                    GreetingMessage:  Agora.String("Hello! I'm ready to chat in real-time."),
                },
                TurnDetection: &Agora.StartAgentsRequestPropertiesTurnDetection{
                    Type:              Agora.StartAgentsRequestPropertiesTurnDetectionTypeServerVad,
                    Threshold:         Agora.Float64(0.5),
                    SilenceDurationMs: Agora.Int(500),
                },
                // TTS and LLM are still required but not used when MLLM is enabled
                Tts: &Agora.StartAgentsRequestPropertiesTts{
                    Vendor: Agora.StartAgentsRequestPropertiesTtsVendorMicrosoft,
                    Params: map[string]any{},
                },
                Llm: &Agora.StartAgentsRequestPropertiesLlm{
                    URL: "https://api.openai.com/v1/chat/completions",
                },
            },
        },
    )
}
```

## MLLM Flow (Multimodal)

For real-time audio processing using OpenAI's Realtime API or Google Gemini Live, use the MLLM (Multimodal Large Language Model) flow instead of the cascading ASR -> LLM -> TTS flow. See the [MLLM Overview](https://docs.agora.io/en/conversational-ai/models/mllm/overview) for more details.

```go
package main

import (
    "context"
    client "github.com/{{ owner }}/{{ repo }}/client"
    option "github.com/{{ owner }}/{{ repo }}/option"
    Agora "github.com/{{ owner }}/{{ repo }}"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customerId>", "<customerSecret>"),
    )

    c.Agents.Start(
        context.TODO(),
        &Agora.StartAgentsRequest{
            Appid: "your_app_id",
            Name:  "mllm_agent",
            Properties: &Agora.StartAgentsRequestProperties{
                Channel:       "channel_name",
                Token:         "your_token",
                AgentRtcUID:   "1001",
                RemoteRtcUIDs: []string{"1002"},
                IdleTimeout:   Agora.Int(120),
                AdvancedFeatures: &Agora.StartAgentsRequestPropertiesAdvancedFeatures{
                    EnableMllm: Agora.Bool(true),
                },
                Mllm: &Agora.StartAgentsRequestPropertiesMllm{
                    URL:    Agora.String("wss://api.openai.com/v1/realtime"),
                    APIKey: Agora.String("<your_openai_api_key>"),
                    Vendor: Agora.StartAgentsRequestPropertiesMllmVendorOpenai,
                    Params: map[string]any{
                        "model": "gpt-4o-realtime-preview",
                        "voice": "alloy",
                    },
                    InputModalities:  []string{"audio"},
                    OutputModalities: []string{"text", "audio"},
                    GreetingMessage:  Agora.String("Hello! I'm ready to chat in real-time."),
                },
                TurnDetection: &Agora.StartAgentsRequestPropertiesTurnDetection{
                    Type:              Agora.StartAgentsRequestPropertiesTurnDetectionTypeServerVad,
                    Threshold:         Agora.Float64(0.5),
                    SilenceDurationMs: Agora.Int(500),
                },
                // TTS and LLM are still required but not used when MLLM is enabled
                Tts: &Agora.StartAgentsRequestPropertiesTts{
                    Vendor: Agora.StartAgentsRequestPropertiesTtsVendorMicrosoft,
                    Params: map[string]any{},
                },
                Llm: &Agora.StartAgentsRequestPropertiesLlm{
                    URL: "https://api.openai.com/v1/chat/completions",
                },
            },
        },
    )
}
```


## Usage

Instantiate and use the client with the following:

```go
package example

import (
    client "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    option "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
    Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
    context "context"
)

func do() {
    client := client.NewClient(
        option.WithBasicAuth(
            "<username>",
            "<password>",
        ),
    )
    request := &Agora.StartAgentsRequest{
        Appid: "appid",
        Name: "unique_name",
        Properties: &Agora.StartAgentsRequestProperties{
            Channel: "channel_name",
            Token: "token",
            AgentRtcUID: "1001",
            RemoteRtcUIDs: []string{
                "1002",
            },
            IdleTimeout: Agora.Int(
                120,
            ),
            Asr: &Agora.StartAgentsRequestPropertiesAsr{
                Language: Agora.String(
                    "en-US",
                ),
            },
            Tts: &Agora.Tts{
                Microsoft: &Agora.MicrosoftTts{
                    Params: &Agora.MicrosoftTtsParams{
                        Key: "key",
                        Region: "region",
                        VoiceName: "voice_name",
                    },
                },
            },
            Llm: &Agora.StartAgentsRequestPropertiesLlm{
                URL: "https://api.openai.com/v1/chat/completions",
                APIKey: Agora.String(
                    "<your_llm_key>",
                ),
                SystemMessages: []map[string]any{
                    map[string]any{
                        "role": "system",
                        "content": "You are a helpful chatbot.",
                    },
                },
                Params: map[string]any{
                    "model": "gpt-4o-mini",
                },
                MaxHistory: Agora.Int(
                    32,
                ),
                GreetingMessage: Agora.String(
                    "Hello, how can I assist you today?",
                ),
                FailureMessage: Agora.String(
                    "Please hold on a second.",
                ),
            },
        },
    }
    client.Agents.Start(
        context.TODO(),
        request,
    )
}
```

## Environments

You can choose between different environments by using the `option.WithBaseURL` option. You can configure any arbitrary base
URL, which is particularly useful in test environments.

```go
client := client.NewClient(
    option.WithBaseURL(Agora.Environments.Default),
)
```

## Errors

Structured error types are returned from API calls that return non-success status codes. These errors are compatible
with the `errors.Is` and `errors.As` APIs, so you can access the error like so:

```go
response, err := client.Agents.Start(...)
if err != nil {
    var apiError *core.APIError
    if errors.As(err, apiError) {
        // Do something with the API error ...
    }
    return err
}
```

## Request Options

A variety of request options are included to adapt the behavior of the library, which includes configuring
authorization tokens, or providing your own instrumented `*http.Client`.

These request options can either be
specified on the client so that they're applied on every request, or for an individual request, like so:

> Providing your own `*http.Client` is recommended. Otherwise, the `http.DefaultClient` will be used,
> and your client will wait indefinitely for a response (unless the per-request, context-based timeout
> is used).

```go
// Specify default options applied on every request.
client := client.NewClient(
    option.WithToken("<YOUR_API_KEY>"),
    option.WithHTTPClient(
        &http.Client{
            Timeout: 5 * time.Second,
        },
    ),
)

// Specify options for an individual request.
response, err := client.Agents.Start(
    ...,
    option.WithToken("<YOUR_API_KEY>"),
)
```

## Advanced

### Response Headers

You can access the raw HTTP response data by using the `WithRawResponse` field on the client. This is useful
when you need to examine the response headers received from the API call. (When the endpoint is paginated,
the raw HTTP response data will be included automatically in the Page response object.)

```go
response, err := client.Agents.WithRawResponse.Start(...)
if err != nil {
    return err
}
fmt.Printf("Got response headers: %v", response.Header)
fmt.Printf("Got status code: %d", response.StatusCode)
```

### Retries

The SDK is instrumented with automatic retries with exponential backoff. A request will be retried as long
as the request is deemed retryable and the number of retry attempts has not grown larger than the configured
retry limit (default: 2).

A request is deemed retryable when any of the following HTTP status codes is returned:

- [408](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408) (Timeout)
- [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) (Too Many Requests)
- [5XX](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500) (Internal Server Errors)

If the `Retry-After` header is present in the response, the SDK will prioritize respecting its value exactly
over the default exponential backoff.

Use the `option.WithMaxAttempts` option to configure this behavior for the entire client or an individual request:

```go
client := client.NewClient(
    option.WithMaxAttempts(1),
)

response, err := client.Agents.Start(
    ...,
    option.WithMaxAttempts(1),
)
```

### Timeouts

Setting a timeout for each individual request is as simple as using the standard context library. Setting a one second timeout for an individual API call looks like the following:

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()

response, err := client.Agents.Start(ctx, ...)
```

### Explicit Null

If you want to send the explicit `null` JSON value through an optional parameter, you can use the setters\
that come with every object. Calling a setter method for a property will flip a bit in the `explicitFields`
bitfield for that setter's object; during serialization, any property with a flipped bit will have its
omittable status stripped, so zero or `nil` values will be sent explicitly rather than omitted altogether:

```go
type ExampleRequest struct {
    // An optional string parameter.
    Name *string `json:"name,omitempty" url:"-"`

    // Private bitmask of fields set to an explicit value and therefore not to be omitted
    explicitFields *big.Int `json:"-" url:"-"`
}

request := &ExampleRequest{}
request.SetName(nil)

response, err := client.Agents.Start(ctx, request, ...)
```

## Contributing

While we value open-source contributions to this SDK, this library is generated programmatically.
Additions made directly to this library would have to be moved over to our generation code,
otherwise they would be overwritten upon the next generated release. Feel free to open a PR as
a proof of concept, but know that we will not be able to merge it as-is. We suggest opening
an issue first to discuss with us!

On the other hand, contributions to the README are always very welcome!