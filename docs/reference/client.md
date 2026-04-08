---
sidebar_position: 1
title: Client Reference
description: API reference for the Fern-generated client, sub-clients, and request options.
---

# Client Reference

## client.NewClient

<!-- snippet: fragment -->
```go
func NewClient(opts ...option.RequestOption) *Client
```

Creates a new API client. All sub-clients share the same configuration.

<!-- snippet: fragment -->
```go
import (
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

c := client.NewClient(
    option.WithToken("<your_rest_auth_token>"),
)
```

## Sub-Clients

| Field | Type | Description |
|---|---|---|
| `c.Agents` | `*agents.Client` | Agent lifecycle (start, stop, speak, interrupt, update, get, getHistory) |
| `c.Telephony` | `*telephony.Client` | Telephony operations (call, hangup) |
| `c.PhoneNumbers` | `*phonenumbers.Client` | Phone number management |

All sub-client methods take `context.Context` as their first argument. See the [generated reference](https://github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/blob/HEAD/./reference.md) for full method signatures.

## Request Options

Request options configure authentication, transport, and retry behavior. They can be set at client creation (applied to all requests) or per-request.

### option.WithToken

<!-- snippet: fragment -->
```go
func WithToken(token string) *core.TokenOption
```

Sets a bearer token for authentication. This is the recommended option for new integrations and matches the onboarding guides.

<!-- snippet: fragment -->
```go
c := client.NewClient(
    option.WithToken("<your_rest_auth_token>"),
)
```

### option.WithBasicAuth

<!-- snippet: fragment -->
```go
func WithBasicAuth(username, password string) *core.BasicAuthOption
```

Sets the `Authorization: Basic <base64>` header using your Agora Customer ID and Customer Secret.

Use this only when you explicitly want Basic Auth. It is still supported by the client, but token auth is the preferred path in the narrative docs and quick starts.

<!-- snippet: fragment -->
```go
c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
)
```

### option.WithBaseURL

<!-- snippet: fragment -->
```go
func WithBaseURL(baseURL string) *core.BaseURLOption
```

Overrides the default API endpoint.

<!-- snippet: fragment -->
```go
import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

c := client.NewClient(
    option.WithBaseURL(Agora.Environments.Default),
)
```

### option.WithArea

<!-- snippet: fragment -->
```go
func WithArea(area core.Area) *core.AreaRequestOption
```

Enables regional routing with automatic DNS-based domain selection. See [Regional Routing](../guides/regional-routing.md).

<!-- snippet: fragment -->
```go
c := client.NewClient(
    option.WithArea(option.AreaUS),
)
```

### option.WithPool

<!-- snippet: fragment -->
```go
func WithPool(pool *core.Pool) *core.AreaRequestOption
```

Uses a pre-configured `Pool` for regional routing. See [Regional Routing](../guides/regional-routing.md).

### option.WithHTTPClient

<!-- snippet: fragment -->
```go
func WithHTTPClient(httpClient core.HTTPClient) *core.HTTPClientOption
```

Provides a custom HTTP client. Recommended for production to set timeouts.

<!-- snippet: fragment -->
```go
import "net/http"

c := client.NewClient(
    option.WithHTTPClient(&http.Client{
        Timeout: 10 * time.Second,
    }),
)
```

### option.WithHTTPHeader

<!-- snippet: fragment -->
```go
func WithHTTPHeader(httpHeader http.Header) *core.HTTPHeaderOption
```

Adds custom HTTP headers to every request.

### option.WithBodyProperties

<!-- snippet: fragment -->
```go
func WithBodyProperties(bodyProperties map[string]interface{}) *core.BodyPropertiesOption
```

Adds extra properties to the JSON request body.

### option.WithQueryParameters

<!-- snippet: fragment -->
```go
func WithQueryParameters(queryParameters url.Values) *core.QueryParametersOption
```

Adds query parameters to the request URL.

### option.WithMaxAttempts

<!-- snippet: fragment -->
```go
func WithMaxAttempts(attempts uint) *core.MaxAttemptsOption
```

Configures the maximum number of retry attempts (default: 2). Retries use exponential backoff for status codes 408, 429, and 5xx.

<!-- snippet: fragment -->
```go
c := client.NewClient(
    option.WithMaxAttempts(3),
)
```

## Area Constants

<!-- snippet: fragment -->
```go
option.AreaUS      // United States (west + east)
option.AreaEU      // Europe (west + central)
option.AreaAP      // Asia-Pacific (southeast + northeast)
option.AreaCN      // Chinese Mainland (east + north)
option.AreaUnknown // Default
```

## Environments

<!-- snippet: fragment -->
```go
import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

Agora.Environments.Default  // "https://api.agora.io/api/conversational-ai-agent"
```

## Pointer Helpers

The root `Agora` package provides helper functions for creating pointers to literal values, required for optional fields in request structs:

| Function | Signature | Example |
|---|---|---|
| `Agora.Bool` | `func(bool) *bool` | `EnableMllm: Agora.Bool(true)` |
| `Agora.Int` | `func(int) *int` | `IdleTimeout: Agora.Int(120)` |
| `Agora.String` | `func(string) *string` | `APIKey: Agora.String("<key>")` |
| `Agora.Float64` | `func(float64) *float64` | `Threshold: Agora.Float64(0.5)` |
| `Agora.Float32` | `func(float32) *float32` | — |
| `Agora.Int8/16/32/64` | `func(intN) *intN` | — |
| `Agora.Uint/8/16/32/64` | `func(uintN) *uintN` | — |
| `Agora.UUID` | `func(uuid.UUID) *uuid.UUID` | — |
| `Agora.Time` | `func(time.Time) *time.Time` | — |
