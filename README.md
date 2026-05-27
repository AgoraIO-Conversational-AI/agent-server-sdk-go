# Agora Agents Go

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=https%3A%2F%2Fgithub.com%2FAgoraIO%2Fagora-agents-go)

The Agora Conversational AI SDK provides convenient access to the Agora Conversational AI APIs, 
enabling you to build voice-powered AI agents with support for both cascading flows (ASR -> LLM -> TTS) 
and multimodal flows (MLLM) for real-time audio processing.


## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Byok](#byok)
- [Mllm Realtime Multimodal](#mllm-realtime-multimodal)
- [Documentation](#documentation)
- [Reference](#reference)
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

## Requirements

- Go 1.21+

## Installation

```sh
go mod init example.com/voice-agent
go get github.com/AgoraIO/agora-agents-go
```

## Quick Start

The recommended onboarding path is a server-side builder flow: define the agent once, configure preset-backed providers in the builder, and let AgentKit infer the reseller `preset` values when the session starts.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/AgoraIO/agora-agents-go/agentkit"
    "github.com/AgoraIO/agora-agents-go/agentkit/vendors"
    "github.com/AgoraIO/agora-agents-go/option"
)

const (
    agentPrompt = "You are a concise, technically credible voice assistant. Keep replies short unless the user asks for detail."
    greeting    = "Hi there! I am your Agora voice assistant. How can I help?"
)

func stringPtr(v string) *string { return &v }
func intPtr(v int) *int { return &v }
func float64Ptr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool { return &v }

func requireEnv(name string) (string, error) {
    value := os.Getenv(name)
    if value == "" {
        return "", fmt.Errorf("missing required environment variable: %s", name)
    }
    return value, nil
}

func startConversation(ctx context.Context) (string, error) {
    appID, err := requireEnv("AGORA_APP_ID")
    if err != nil {
        return "", err
    }
    appCertificate, err := requireEnv("AGORA_APP_CERTIFICATE")
    if err != nil {
        return "", err
    }
    expiresIn, err := agentkit.ExpiresInHours(1)
    if err != nil {
        return "", err
    }

    client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
        Area:           option.AreaUS,
        AppID:          appID,
        AppCertificate: appCertificate,
    })

    agent := agentkit.NewAgent(
        agentkit.WithName(fmt.Sprintf("conversation-%d", time.Now().UnixMilli())),
        agentkit.WithInstructions(agentPrompt),
        agentkit.WithGreeting(greeting),
        agentkit.WithFailureMessage("Please wait a moment."),
        agentkit.WithMaxHistory(50),
        agentkit.WithTurnDetectionConfig(&agentkit.TurnDetectionConfig{
            Config: &agentkit.TurnDetectionNestedConfig{
                SpeechThreshold: float64Ptr(0.5),
                StartOfSpeech: &agentkit.StartOfSpeechConfig{
                    Mode: agentkit.StartOfSpeechMode("vad"),
                    VadConfig: &agentkit.StartOfSpeechVadConfig{
                        InterruptDurationMs: intPtr(160),
                        PrefixPaddingMs:     intPtr(300),
                    },
                },
                EndOfSpeech: &agentkit.EndOfSpeechConfig{
                    Mode: agentkit.EndOfSpeechMode("vad"),
                    VadConfig: &agentkit.EndOfSpeechVadConfig{
                        SilenceDurationMs: intPtr(480),
                    },
                },
            },
        }),
        agentkit.WithAdvancedFeatures(&agentkit.AdvancedFeatures{
            EnableRtm:   boolPtr(true),
            EnableTools: boolPtr(true),
        }),
        agentkit.WithParameters(&agentkit.SessionParams{
            DataChannel:        &agentkit.DataChannelRtm,
            EnableErrorMessage: boolPtr(true),
        }),
    ).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
        Model: "nova-3",
        Language: "en",
    })).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
        Model:           "gpt-4o-mini",
        GreetingMessage: greeting,
        FailureMessage:  "Please wait a moment.",
        MaxHistory:      intPtr(15),
        Params: map[string]interface{}{
            "max_tokens": 1024,
            "temperature": 0.7,
            "top_p": 0.95,
        },
    })).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
        Model:   "speech_2_6_turbo",
        VoiceID: "English_captivating_female1",
    }))

    session := agent.CreateSession(client, agentkit.CreateSessionOptions{
        Channel:     fmt.Sprintf("demo-channel-%d", time.Now().UnixMilli()),
        AgentUID:    "123456",
        RemoteUIDs:  []string{"*"},
        IdleTimeout: intPtr(30),
        ExpiresIn:   expiresIn,
        Debug:       false,
    })

    return session.Start(ctx)
}

func main() {
    agentID, err := startConversation(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Println(agentID)
}
```

### Why no token or vendor key in the example?

`AgoraClient` generates the required ConvoAI REST auth and RTC join tokens automatically when you provide `AppID` and `AppCertificate`. AgentKit then inspects the builder-provided vendor configs and infers the matching supported `preset` values for reseller-backed models, so you do not pass vendor API keys in this flow.

### BYOK version of the same builder flow

Use the same `Agent` builder shape, but provide credentials explicitly when you want vendor-managed billing and routing instead of Agora-managed presets.

```go
agent := agentkit.NewAgent(
    agentkit.WithInstructions(agentPrompt),
    agentkit.WithGreeting(greeting),
).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
    APIKey:   os.Getenv("DEEPGRAM_API_KEY"),
    Model:    "nova-3",
    Language: "en",
})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
    APIKey:      os.Getenv("OPENAI_API_KEY"),
    Model:       "gpt-4o-mini",
    MaxTokens:   intPtr(1024),
    Temperature: float64Ptr(0.7),
    TopP:        float64Ptr(0.95),
})).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
    Key:     os.Getenv("MINIMAX_API_KEY"),
    GroupID: os.Getenv("MINIMAX_GROUP_ID"),
    Model:   "speech_2_6_turbo",
    VoiceID: "English_captivating_female1",
    URL:     "wss://api-uw.minimax.io/ws/v1/t2a_v2",
}))
```

## BYOK

If you want to bring your own vendor credentials instead of using Agora-managed presets, use the BYOK guide:

- [BYOK Guide](./docs/guides/byok.md)

## MLLM (Realtime / Multimodal)

Use `WithMllm()` for OpenAI Realtime or Gemini Live. No STT, LLM, or TTS vendor is needed when MLLM mode is enabled.

```go
agent := agentkit.NewAgent(
    agentkit.WithName("realtime-assistant"),
).WithMllm(vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
    APIKey:          os.Getenv("OPENAI_API_KEY"),
    Model:           "gpt-4o-realtime-preview",
    GreetingMessage: "Hello! Ready to chat.",
}))
```

See the [MLLM Flow guide](./docs/guides/mllm-flow.md) for full examples with Gemini Live and Vertex AI.

## Documentation

API reference documentation is available [here](https://docs.agora.io/en/conversational-ai/overview).

## Reference

A full reference for this library is available [here](https://github.com/AgoraIO/agora-agents-go/blob/HEAD/./reference.md).

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
    client "github.com/AgoraIO/agora-agents-go/client"
    option "github.com/AgoraIO/agora-agents-go/option"
    Agora "github.com/AgoraIO/agora-agents-go"
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
            TurnDetection: &Agora.StartAgentsRequestPropertiesTurnDetection{
                Config: &Agora.StartAgentsRequestPropertiesTurnDetectionConfig{
                    EndOfSpeech: &Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeech{
                        Mode: Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeechModeSemantic.Ptr(),
                    },
                },
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
