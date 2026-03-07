---
sidebar_position: 6
title: Error Handling
description: Handle API errors with core.APIError and errors.As.
---

# Error Handling

When the API returns a non-success status code (4xx or 5xx response), the SDK returns a `*core.APIError`. Use `errors.As` to inspect it:

```go
package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"

    core "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/core"
)

func main() {
    resp, err := client.Agents.Start(ctx, req)
    if err != nil {
        var apiErr *core.APIError
        if errors.As(err, &apiErr) {
            fmt.Println("Status:", apiErr.StatusCode)
            fmt.Println("Message:", apiErr.Error())
            fmt.Println("Headers:", apiErr.Header)
        } else {
            log.Fatal(err)
        }
        return
    }
    _ = resp
}
```

## APIError Fields

| Field | Type | Description |
|-------|------|--------------|
| `StatusCode` | `int` | HTTP status code (e.g., 400, 404, 500) |
| `Header` | `http.Header` | Response headers |
| `Error()` | `string` | Combined status and message |

## Graceful 404 Handling

For operations like `Stop` where a 404 may indicate the agent was already stopped, check the status code:

```go
err := session.Stop(ctx)
if err != nil {
    var apiErr *core.APIError
    if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
        // Agent already stopped or never existed — treat as success
        return nil
    }
    return err
}
```

Use this for both the agentkit layer (e.g. `session.Start()`, `session.Stop()`) and direct client calls (`client.Agents.Start(...)`).
