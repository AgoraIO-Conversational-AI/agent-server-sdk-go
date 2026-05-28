---
sidebar_position: 8
title: Advanced
description: Timeouts, cancellation, stopping by agent ID, raw responses, and generated-client escape hatches.
---

# Advanced

The default path for realtime agents is still AgentKit with app credentials:

```go
client := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
    Area:           option.AreaUS,
    AppID:          "your-app-id",
    AppCertificate: "your-app-certificate",
})
```

That client keeps auth out of your request code: session methods mint ConvoAI REST auth and RTC join tokens from `AppID` and `AppCertificate`.

## Timeouts

Use `context.WithTimeout` for request-scoped timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

agentID, err := session.Start(ctx)
if err != nil {
    return err
}
_ = agentID
```

## Cancellation

Use `context.WithCancel` when your server needs to abort work after a client disconnects or workflow cancellation:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    <-done
    cancel()
}()

if err := session.Stop(ctx); err != nil {
    return err
}
```

## Stop Without A Session Handle

Use `StopAgent` when a later request handler only has the agent session ID. It uses the same app credentials to generate the required ConvoAI REST auth header.

```go
if err := client.StopAgent(ctx, agentID); err != nil {
    return err
}
```

## Raw Responses

The generated clients expose `WithRawResponse` for advanced debugging of status codes and headers. Keep session lifecycle code on AgentKit unless you specifically need a generated-client escape hatch.

```go
resp, err := client.Agents.WithRawResponse.Get(ctx, req)
if err != nil {
    return err
}

fmt.Println(resp.StatusCode)
fmt.Println(resp.Header)
fmt.Println(resp.Body)
```

## Generated-Client Options

Generated-client request options such as `option.WithMaxAttempts`, `option.WithHTTPHeader`, and `option.WithQueryParameters` are available for low-level REST calls. They are advanced tools and are not needed for the recommended AgentKit session flow.
