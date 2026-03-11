---
sidebar_position: 3
title: Session Reference
description: Complete API reference for agentkit.AgentSession — lifecycle methods, state machine, and events.
---

# Session Reference

Package: `github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit`

## NewAgentSession

```go
func NewAgentSession(opts AgentSessionOptions) *AgentSession
```

Creates a new session. If `Name` is empty, defaults to `agent-<unix_timestamp>`. The session starts in `StatusIdle`.

### AgentSessionOptions

```go
type AgentSessionOptions struct {
    Client          *agents.Client
    Agent           *Agent
    AppID           string
    AppCertificate  string
    Name            string
    Channel         string
    Token           string
    AgentUID        string
    RemoteUIDs      []string
    IdleTimeout     *int
    EnableStringUID *bool
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `Client` | `*agents.Client` | Yes | Fern-generated agents sub-client (from `c.Agents`) |
| `Agent` | `*Agent` | Yes | Agent configuration |
| `AppID` | `string` | Yes | Agora App ID |
| `AppCertificate` | `string` | Conditional | Required if `Token` is not set |
| `Name` | `string` | No | Session name (default: `agent-<unix_timestamp>`) |
| `Channel` | `string` | Yes | Agora channel name |
| `Token` | `string` | Conditional | Pre-generated RTC token |
| `AgentUID` | `string` | Yes | Agent's UID in the channel |
| `RemoteUIDs` | `[]string` | Yes | Remote participant UIDs |
| `IdleTimeout` | `*int` | No | Idle timeout in seconds |
| `EnableStringUID` | `*bool` | No | Enable string UID mode |

## SessionStatus

```go
type SessionStatus string

const (
    StatusIdle     SessionStatus = "idle"
    StatusStarting SessionStatus = "starting"
    StatusRunning  SessionStatus = "running"
    StatusStopping SessionStatus = "stopping"
    StatusStopped  SessionStatus = "stopped"
    StatusError    SessionStatus = "error"
)
```

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
     │               │                │ success
     │               ▼                ▼
     │          ┌─────────┐      ┌─────────┐
     └─────────>│  (can   │      │ stopped │
     (restart)  │ restart)│      └─────────┘
                └─────────┘
```

## Methods

### Start

```go
func (s *AgentSession) Start(ctx context.Context) (string, error)
```

Starts the agent session. Returns the agent ID assigned by the API.

- **Valid from:** `idle`, `stopped`, `error`
- **Transitions to:** `starting` -> `running` (success) or `error` (failure)
- **Emits:** `"started"` on success, `"error"` on failure
- **Validates:** Avatar/TTS sample rate match before making the API call

### Stop

```go
func (s *AgentSession) Stop(ctx context.Context) error
```

Stops the running agent.

- **Valid from:** `running`
- **Transitions to:** `stopping` -> `stopped` (success) or `error` (failure)
- **Emits:** `"stopped"` on success, `"error"` on failure

### Say

```go
func (s *AgentSession) Say(ctx context.Context, text string, priority *Agora.SpeakAgentsRequestPriority, interruptable *bool) error
```

Sends text for the agent to speak.

- **Valid from:** `running`
- **Parameters:**
  - `text` — the text to speak
  - `priority` — optional priority level (pass `nil` for default)
  - `interruptable` — whether the utterance can be interrupted (pass `nil` for default)

### Interrupt

```go
func (s *AgentSession) Interrupt(ctx context.Context) error
```

Interrupts the agent's current speech.

- **Valid from:** `running`

### Update

```go
func (s *AgentSession) Update(ctx context.Context, properties *Agora.UpdateAgentsRequestProperties) error
```

Updates the agent's properties while running.

- **Valid from:** `running`

### GetHistory

```go
func (s *AgentSession) GetHistory(ctx context.Context) (*Agora.GetHistoryAgentsResponse, error)
```

Retrieves conversation history.

- **Requires:** Valid agent ID (any state after successful `Start`)

### GetInfo

```go
func (s *AgentSession) GetInfo(ctx context.Context) (*Agora.GetAgentsResponse, error)
```

Gets the current agent status from the API.

- **Requires:** Valid agent ID

## Getters

```go
func (s *AgentSession) ID() string
```
Returns the agent ID (empty string before `Start` succeeds).

```go
func (s *AgentSession) Status() SessionStatus
```
Returns the current session state.

```go
func (s *AgentSession) Agent() *Agent
```
Returns the agent configuration.

```go
func (s *AgentSession) AppID() string
```
Returns the App ID.

```go
func (s *AgentSession) Raw() *agents.Client
```
Returns the underlying Fern-generated agents client for direct API access.

## Event System

### On

```go
func (s *AgentSession) On(event string, handler EventHandler)
```

Registers an event handler. Multiple handlers can be registered for the same event.

### EventHandler

```go
type EventHandler func(data interface{})
```

### Events

| Event | Data Type | When |
|---|---|---|
| `"started"` | `map[string]string{"agent_id": "..."}` | `Start()` succeeds |
| `"stopped"` | `map[string]string{"agent_id": "..."}` | `Stop()` succeeds |
| `"error"` | `error` | `Start()` or `Stop()` fails |

Handlers run synchronously. Panics in handlers are recovered and silently discarded. Register handlers before calling `Start()`.

## Thread Safety

All state access is protected by `sync.RWMutex`. The session is safe for concurrent use across goroutines.
