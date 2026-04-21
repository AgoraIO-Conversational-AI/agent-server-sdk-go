---
sidebar_position: 7
title: Pagination
description: Iterate over paginated list endpoints.
---

# Pagination

List endpoints such as `client.Agents.List` and `client.Telephony.List` return a `*core.Page` that you can iterate over.

## Iterating Over All Items

Use the `Iterator()` method to loop over all items across pages:

```go
package main

import (
    "context"
    "fmt"
    "log"

    Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
    "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

func main() {
    c := client.NewClient(
        option.WithToken("<your_rest_auth_token>"),
    )

    ctx := context.Background()
    page, err := c.Agents.List(ctx, &Agora.ListAgentsRequest{
        Appid: "your_app_id",
    })
    if err != nil {
        log.Fatal(err)
    }

    iter := page.Iterator()
    for iter.Next(ctx) {
        item := iter.Current()
        fmt.Println(item)
    }
    if iter.Err() != nil {
        log.Fatal(iter.Err())
    }
}
```

## Manual Page-by-Page Iteration

Fetch pages explicitly with `GetNextPage`:

```go
page, err := c.Agents.List(ctx, &Agora.ListAgentsRequest{Appid: "your_app_id"})
if err != nil {
    log.Fatal(err)
}

for {
    for _, item := range page.Results {
        fmt.Println(item)
    }
    var nextErr error
    page, nextErr = page.GetNextPage(ctx)
    if nextErr != nil || page == nil {
        break
    }
}
```

When no more pages exist, `GetNextPage` returns `core.ErrNoPages`. The iterator treats this as a normal end-of-stream (not an error).
