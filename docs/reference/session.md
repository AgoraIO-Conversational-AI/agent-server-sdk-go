---
sidebar_position: 3
title: Session Reference
description: Complete API reference for agentkit.AgentSession — lifecycle methods, state machine, and events.
---

# Session Reference

Package: `github.com/AgoraIO/agora-agents-go/agentkit`

## NewAgentSession

<!-- snippet: fragment -->
```go
func NewAgentSession(opts AgentSessionOptions) *AgentSession
```

Creates a new session. If `Name` is empty, defaults to `agent-<unix_timestamp>`. The session starts in `StatusIdle`.

### AgentSessionOptions

<!-- snippet: fragment -->
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
    ExpiresIn       int
    UseAppCredentialsForREST bool
    Preset          []string
    PipelineID      string
    Debug           bool
    Warn            func(string)
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
| `ExpiresIn` | `int` | No | Auto-generated token lifetime in seconds |
| `UseAppCredentialsForREST` | `bool` | No | Generate ConvoAI REST auth headers per request |
| `Preset` | `[]string` | No | Advanced feature flag for internal/explicit model routing. Leave unset for normal builder usage. |
| `PipelineID` | `string` | No | Published pipeline ID to send on session start |
| `Debug` | `bool` | No | Enable debug logging of the start request |
| `Warn` | `func(string)` | No | Custom warning sink; defaults to logger |

## SessionStatus

<!-- snippet: fragment -->
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

<!-- snippet: fragment -->
```go
func (s *AgentSession) Start(ctx context.Context) (string, error)
```

Starts the agent session. Returns the agent ID assigned by the API.

- **Valid from:** `idle`, `stopped`, `error`
- **Transitions to:** `starting` -> `running` (success) or `error` (failure)
- **Emits:** `"started"` on success, `"error"` on failure
- **Validates:** Avatar config and avatar/TTS sample rate match before making the API call

### Stop

<!-- snippet: fragment -->
```go
func (s *AgentSession) Stop(ctx context.Context) error
```

Stops the running agent.

- **Valid from:** `running`
- **Transitions to:** `stopping` -> `stopped` (success) or `error` (failure)
- **Emits:** `"stopped"` on success, `"error"` on failure

### Say

<!-- snippet: fragment -->
```go
func (s *AgentSession) Say(ctx context.Context, text string, priority *Agora.SpeakAgentsRequestPriority, interruptable *bool) error
```

Sends text for the agent to speak.

- **Valid from:** `running`
- **Parameters:**
  - `text` — the text to speak
  - `priority` — optional priority level (pass `nil` for default). Use the `agentkit.SpeakPriorityInterrupt.Ptr()`, `agentkit.SpeakPriorityAppend.Ptr()`, or `agentkit.SpeakPriorityIgnore.Ptr()` convenience constants instead of the raw generated enum.
  - `interruptable` — whether the utterance can be interrupted (pass `nil` for default)

### Interrupt

<!-- snippet: fragment -->
```go
func (s *AgentSession) Interrupt(ctx context.Context) error
```

Interrupts the agent's current speech.

- **Valid from:** `running`

### Update

<!-- snippet: fragment -->
```go
func (s *AgentSession) Update(ctx context.Context, properties *Agora.UpdateAgentsRequestProperties) error
```

Updates the agent's properties while running.

- **Valid from:** `running`

### GetHistory

<!-- snippet: fragment -->
```go
func (s *AgentSession) GetHistory(ctx context.Context) (*Agora.GetHistoryAgentsResponse, error)
```

Retrieves conversation history.

- **Requires:** Valid agent ID (any state after successful `Start`)

### GetInfo

<!-- snippet: fragment -->
```go
func (s *AgentSession) GetInfo(ctx context.Context) (*Agora.GetAgentsResponse, error)
```

Gets the current agent status from the API.

- **Requires:** Valid agent ID

### GetTurns

<!-- snippet: fragment -->
```go
func (s *AgentSession) GetTurns(ctx context.Context, opts ...GetTurnsOptions) (*Agora.GetTurnsAgentsResponse, error)
func (s *AgentSession) GetAllTurns(ctx context.Context, opts ...GetAllTurnsOptions) (*Agora.GetTurnsAgentsResponse, error)

type GetTurnsOptions struct {
    PageIndex *int
    PageSize  *int
}

type GetAllTurnsOptions struct {
    PageSize *int
}
```

Retrieves turn-by-turn analytics for the session. `PageIndex` starts at 1. Use `GetAllTurns` to iterate through every page with a default page size of 50 and return the final response with aggregated `Turns`.

- **Requires:** Valid agent ID

When you consume server notifications, event `112` means all turns for the session have finished and are ready to query.

### Think

<!-- snippet: fragment -->
```go
func (s *AgentSession) Think(ctx context.Context, text string, onListeningAction *Agora.AgentThinkAgentManagementRequestOnListeningAction, onThinkingAction *Agora.AgentThinkAgentManagementRequestOnThinkingAction, onSpeakingAction *Agora.AgentThinkAgentManagementRequestOnSpeakingAction, interruptable *bool, metadata map[string]string) (*Agora.AgentThinkAgentManagementResponse, error)
func (s *AgentSession) ThinkWithOptions(ctx context.Context, text string, opts *ThinkOptions) (*Agora.AgentThinkAgentManagementResponse, error)
```

Injects a thought or instruction into a running agent. In v2.7, omitting `on_listening_action` uses the server default `interrupt`. Set `agentkit.ThinkOnListeningActionInject.Ptr()` if you need legacy inject behavior. AgentKit also exposes `ThinkOnThinkingActionInterrupt`, `ThinkOnThinkingActionIgnore`, `ThinkOnSpeakingActionInterrupt`, and `ThinkOnSpeakingActionIgnore` convenience constants.

## Getters

<!-- snippet: fragment -->
```go
func (s *AgentSession) ID() string
```
Returns the agent ID (empty string before `Start` succeeds).

<!-- snippet: fragment -->
```go
func (s *AgentSession) Status() SessionStatus
```
Returns the current session state.

<!-- snippet: fragment -->
```go
func (s *AgentSession) Agent() *Agent
```
Returns the agent configuration.

<!-- snippet: fragment -->
```go
func (s *AgentSession) AppID() string
```
Returns the App ID.

<!-- snippet: fragment -->
```go
func (s *AgentSession) Raw() *agents.Client
```
Returns the underlying Fern-generated agents client for direct API access.

## Event System

### On

<!-- snippet: fragment -->
```go
func (s *AgentSession) On(event string, handler EventHandler)
```

Registers an event handler. Multiple handlers can be registered for the same event.

### Off

<!-- snippet: fragment -->
```go
func (s *AgentSession) Off(event string, handler EventHandler)
```

Unregisters a previously registered event handler.

### EventHandler

<!-- snippet: fragment -->
```go
type EventHandler func(data interface{})
```

### Events

| Event | Data Type | When |
|---|---|---|
| `"started"` | `map[string]string{"agent_id": "..."}` | `Start()` succeeds |
| `"stopped"` | `map[string]string{"agent_id": "..."}` | `Stop()` succeeds |
| `"error"` | `error` | `Start()` or `Stop()` fails |

Handlers run synchronously. Panics in handlers are recovered and reported through the session warning sink. Register handlers before calling `Start()`.

## Thread Safety

All state access is protected by `sync.RWMutex`. The session is safe for concurrent use across goroutines.
