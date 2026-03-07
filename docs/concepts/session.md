---
sidebar_position: 3
title: Session
description: AgentSession lifecycle — state machine, methods, and event handling.
---

# Session

`agentkit.AgentSession` manages the lifecycle of a running agent. It wraps the Fern-generated `agents.Client` and provides state tracking, event emission, and convenience methods.

## Creating a Session

```go
session := agentkit.NewAgentSession(agentkit.AgentSessionOptions{
    Client:         c.Agents,
    Agent:          a,
    AppID:          "<your_app_id>",
    AppCertificate: "<your_app_certificate>",
    Channel:        "my-channel",
    AgentUID:       "1001",
    RemoteUIDs:     []string{"1002"},
})
```

### AgentSessionOptions Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `Client` | `*agents.Client` | Yes | The Fern-generated agents sub-client |
| `Agent` | `*Agent` | Yes | Agent configuration built with `NewAgent` |
| `AppID` | `string` | Yes | Agora App ID |
| `AppCertificate` | `string` | Conditional | Required if `Token` is not provided |
| `Name` | `string` | No | Session name (defaults to `agent-<unix_timestamp>`) |
| `Channel` | `string` | Yes | Agora channel name |
| `Token` | `string` | Conditional | Pre-generated RTC token (skips auto-generation) |
| `AgentUID` | `string` | Yes | Agent's UID in the channel |
| `RemoteUIDs` | `[]string` | Yes | Remote participant UIDs |
| `IdleTimeout` | `*int` | No | Idle timeout in seconds |
| `EnableStringUID` | `*bool` | No | Enable string UIDs |

## State Machine

```
         Start()           API success
  ┌──────┐      ┌──────────┐      ┌─────────┐
  │ idle │─────>│ starting │─────>│ running │
  └──┬───┘      └────┬─────┘      └────┬────┘
     │               │                  │
     │               │ error            │ Stop()
     │               ▼                  ▼
     │          ┌─────────┐      ┌──────────┐
     │          │  error  │      │ stopping │
     │          └────┬────┘      └────┬─────┘
     │               │                │
     │               │                │ API success
     │               ▼                ▼
     │          ┌─────────┐      ┌─────────┐
     └─────────>│  (can   │      │ stopped │
     (restart)  │ restart)│      └─────────┘
                └─────────┘
```

### Session States

| State | Value | Description |
|---|---|---|
| `StatusIdle` | `"idle"` | Initial state, ready to start |
| `StatusStarting` | `"starting"` | `Start()` called, waiting for API response |
| `StatusRunning` | `"running"` | Agent is active and processing |
| `StatusStopping` | `"stopping"` | `Stop()` called, waiting for API response |
| `StatusStopped` | `"stopped"` | Agent has stopped (can restart) |
| `StatusError` | `"error"` | An error occurred (can restart) |

### Valid Transitions

- `Start()` is valid from: `idle`, `stopped`, `error`
- `Stop()` is valid from: `running`
- `Say()`, `Interrupt()`, `Update()` are valid from: `running`
- `GetHistory()`, `GetInfo()` require a valid agent ID (any state after a successful start)

## Methods

All methods that make API calls take `context.Context` as the first argument and return an `error`.

### Start

```go
agentID, err := session.Start(ctx)
if err != nil {
    log.Fatalf("Failed to start: %v", err)
}
fmt.Println("Agent ID:", agentID)
```

Transitions: `idle`/`stopped`/`error` -> `starting` -> `running` (or `error`).
Returns the agent ID string assigned by the API.

### Stop

```go
err := session.Stop(ctx)
if err != nil {
    log.Fatalf("Failed to stop: %v", err)
}
```

Transitions: `running` -> `stopping` -> `stopped` (or `error`).

### Say

Send text for the agent to speak:

```go
err := session.Say(ctx, "Hello, welcome!", nil, nil)
if err != nil {
    log.Printf("Say failed: %v", err)
}
```

Parameters:
- `text string` — the text to speak
- `priority *Agora.SpeakAgentsRequestPriority` — optional priority level
- `interruptable *bool` — whether this utterance can be interrupted

### Interrupt

Interrupt the agent's current speech:

```go
err := session.Interrupt(ctx)
if err != nil {
    log.Printf("Interrupt failed: %v", err)
}
```

### Update

Update agent properties while running:

```go
err := session.Update(ctx, &Agora.UpdateAgentsRequestProperties{
    // Updated properties...
})
if err != nil {
    log.Printf("Update failed: %v", err)
}
```

### GetHistory

Retrieve conversation history:

```go
history, err := session.GetHistory(ctx)
if err != nil {
    log.Printf("GetHistory failed: %v", err)
}
```

### GetInfo

Get the current agent status from the API:

```go
info, err := session.GetInfo(ctx)
if err != nil {
    log.Printf("GetInfo failed: %v", err)
}
```

## Getters

These do not make API calls:

```go
session.ID() string              // Agent ID (set after Start)
session.Status() SessionStatus   // Current state
session.Agent() *Agent           // Agent configuration
session.AppID() string           // App ID
session.Raw() *agents.Client     // Underlying Fern-generated client
```

## Event System

`AgentSession` supports an event handler pattern for reacting to lifecycle changes:

```go
session.On("started", func(data interface{}) {
    info := data.(map[string]string)
    fmt.Println("Agent started with ID:", info["agent_id"])
})

session.On("stopped", func(data interface{}) {
    fmt.Println("Agent stopped")
})

session.On("error", func(data interface{}) {
    err := data.(error)
    log.Println("Agent error:", err)
})
```

### Available Events

| Event | Data Type | Emitted When |
|---|---|---|
| `"started"` | `map[string]string{"agent_id": "..."}` | `Start()` succeeds |
| `"stopped"` | `map[string]string{"agent_id": "..."}` | `Stop()` succeeds |
| `"error"` | `error` | Any API call fails during Start or Stop |

Event handlers run synchronously in the calling goroutine. Panics in handlers are recovered and silently discarded. Register handlers before calling `Start()`.

## Thread Safety

`AgentSession` is safe for concurrent use. All state mutations are protected by a `sync.RWMutex`. Multiple goroutines can safely call `Status()`, `ID()`, and other getters while another goroutine calls `Start()` or `Stop()`.
