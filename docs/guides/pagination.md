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

    Agora "github.com/AgoraIO/agora-agents-go"
    "github.com/AgoraIO/agora-agents-go/client"
    "github.com/AgoraIO/agora-agents-go/option"
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

## AgentKit GetTurns Pagination

`AgentSession.GetTurns` uses page-number pagination rather than `core.Page`. Pass `GetTurnsOptions` to fetch a specific page, or call `GetAllTurns` to aggregate every page.

```go
pageIndex := 1
pageSize := 50

turnsPage, err := session.GetTurns(ctx, agentkit.GetTurnsOptions{
    PageIndex: &pageIndex,
    PageSize:  &pageSize,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("turns on page:", len(turnsPage.Turns))
```

```go
turnsResponse, err := session.GetAllTurns(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("all turns:", len(turnsResponse.Turns))
```

If you also subscribe to notifications, event `112` indicates the session turns have finished and are ready to query.
