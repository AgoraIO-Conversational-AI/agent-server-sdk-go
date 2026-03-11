---
sidebar_position: 8
title: Advanced
description: Headers, retries, timeouts, raw response, and custom HTTP client.
---

# Advanced

## Additional Headers

Send additional headers with `option.WithHTTPHeader`:

```go
package main

import (
    "net/http"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    c := client.NewClient(
        option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    )

    headers := http.Header{}
    headers.Set("X-Custom-Header", "custom value")

    resp, err := c.Agents.Start(ctx, req, option.WithHTTPHeader(headers))
    _ = resp
    _ = err
}
```

## Additional Query Parameters

Use `option.WithQueryParameters`:

```go
import (
    "net/url"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

params := url.Values{}
params.Set("customQueryParamKey", "custom query param value")

resp, err := c.Agents.Start(ctx, req, option.WithQueryParameters(params))
```

## Retries

The SDK retries automatically with exponential backoff when a request returns:

- **408** (Request Timeout)
- **429** (Too Many Requests)
- **5XX** (Internal Server Errors)

Default retry limit: 2. Override with `option.WithMaxAttempts`:

```go
c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    option.WithMaxAttempts(0), // disable retries
)

// Or per-request:
resp, err := c.Agents.Start(ctx, req, option.WithMaxAttempts(1))
```

## Timeouts

Go uses `context` for request-scoped timeouts. Use `context.WithTimeout`:

```go
package main

import (
    "context"
    "time"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
)

func main() {
    c := client.NewClient(option.WithBasicAuth("<customer_id>", "<customer_secret>"))

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    resp, err := c.Agents.Start(ctx, req)
    _ = resp
    _ = err
}
```

For a global default timeout, pass a custom `*http.Client` with `Timeout` set:

```go
import (
    "net/http"
    "time"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    option.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

## Cancelling Requests

Use `context.WithCancel` or `context.WithTimeout` and cancel when needed:

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    time.Sleep(2 * time.Second)
    cancel() // abort the request
}()

resp, err := c.Agents.Start(ctx, req)
```

## Access Raw Response Data

Use `client.Agents.WithRawResponse` to get headers and status code:

```go
import (
    "fmt"
    "log"
)

resp, err := c.Agents.WithRawResponse.Start(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.StatusCode)
fmt.Println(resp.Header.Get("X-My-Header"))
fmt.Println(resp.Body) // typed response body
```

## Custom HTTP Client

Pass a custom `*http.Client` for proxies, custom transports, or other options:

```go
import (
    "net/http"
    "net/url"

    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

proxyURL, _ := url.Parse("http://my.proxy.example.com:8080")

c := client.NewClient(
    option.WithBasicAuth("<customer_id>", "<customer_secret>"),
    option.WithHTTPClient(&http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
        Timeout: 60 * time.Second,
    }),
)
```
