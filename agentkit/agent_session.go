package agentkit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agents"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

type SessionStatus string

const (
	StatusIdle     SessionStatus = "idle"
	StatusStarting SessionStatus = "starting"
	StatusRunning  SessionStatus = "running"
	StatusStopping SessionStatus = "stopping"
	StatusStopped  SessionStatus = "stopped"
	StatusError    SessionStatus = "error"
)

type EventHandler func(data interface{})

type AgentSession struct {
	client          *agents.Client
	agent           *Agent
	appID           string
	appCertificate  string
	name            string
	channel         string
	token           string
	agentUID        string
	remoteUIDs      []string
	idleTimeout     *int
	enableStringUID *bool
	expiresIn       int  // Token lifetime in seconds (0 = use DefaultExpirySeconds)
	useAppCredsREST bool // When true, generate ConvoAI token per request for REST API auth

	agentID  string
	status   SessionStatus
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

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
	// ExpiresIn is the token lifetime in seconds (default: 86400 = 24 hours, Agora maximum).
	// Only applies when the SDK auto-generates a token. Valid range: 1–86400.
	// Use ExpiresInHours() / ExpiresInMinutes() for clarity.
	ExpiresIn int
	// UseAppCredentialsForREST when true, generates a ConvoAI token per request for REST API
	// authentication. Use when the client was created without Basic Auth or token (app-credentials mode).
	UseAppCredentialsForREST bool
}

func NewAgentSession(opts AgentSessionOptions) *AgentSession {
	name := opts.Name
	if name == "" {
		name = fmt.Sprintf("agent-%d", time.Now().Unix())
	}

	return &AgentSession{
		client:             opts.Client,
		agent:              opts.Agent,
		appID:              opts.AppID,
		appCertificate:     opts.AppCertificate,
		name:               name,
		channel:            opts.Channel,
		token:              opts.Token,
		agentUID:           opts.AgentUID,
		remoteUIDs:         opts.RemoteUIDs,
		idleTimeout:        opts.IdleTimeout,
		enableStringUID:    opts.EnableStringUID,
		expiresIn:          opts.ExpiresIn,
		useAppCredsREST:    opts.UseAppCredentialsForREST,
		status:             StatusIdle,
		handlers:           make(map[string][]EventHandler),
	}
}

// convoAIRequestOpts returns per-request options with ConvoAI token when using app credentials.
func (s *AgentSession) convoAIRequestOpts(ctx context.Context) []option.RequestOption {
	if !s.useAppCredsREST || s.appCertificate == "" {
		return nil
	}
	token, err := GenerateConvoAIToken(GenerateConvoAITokenOptions{
		AppID:          s.appID,
		AppCertificate: s.appCertificate,
		ChannelName:    s.channel,
		Account:        s.agentUID,
	})
	if err != nil {
		// Log and fall through without auth headers; the API call will fail with
		// an auth error, but this log surfaces the real cause.
		log.Printf("agentkit: failed to generate ConvoAI token: %v", err)
		return nil
	}
	h := make(http.Header)
	h.Set("Authorization", "agora token="+token)
	return []option.RequestOption{option.WithHTTPHeader(h)}
}

func (s *AgentSession) ID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.agentID
}

func (s *AgentSession) Status() SessionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *AgentSession) Agent() *Agent {
	return s.agent
}

func (s *AgentSession) AppID() string {
	return s.appID
}

func (s *AgentSession) Raw() *agents.Client {
	return s.client
}

func (s *AgentSession) Start(ctx context.Context) (string, error) {
	s.mu.Lock()
	if s.status != StatusIdle && s.status != StatusStopped && s.status != StatusError {
		s.mu.Unlock()
		return "", fmt.Errorf("cannot start session in %s state", s.status)
	}

	if s.agent.avatarRequiredSampleRate != nil && s.agent.ttsSampleRate != nil {
		if *s.agent.ttsSampleRate != *s.agent.avatarRequiredSampleRate {
			s.mu.Unlock()
			return "", fmt.Errorf(
				"avatar requires TTS sample rate of %d Hz, but TTS is configured with %d Hz",
				int(*s.agent.avatarRequiredSampleRate), int(*s.agent.ttsSampleRate),
			)
		}
	}

	s.status = StatusStarting
	s.mu.Unlock()

	propOpts := ToPropertiesOptions{
		Channel:         s.channel,
		AgentUID:        s.agentUID,
		RemoteUIDs:      s.remoteUIDs,
		Token:           s.token,
		AppID:           s.appID,
		AppCertificate:  s.appCertificate,
		ExpiresIn:       s.expiresIn,
		IdleTimeout:     s.idleTimeout,
		EnableStringUID: s.enableStringUID,
	}

	properties, err := s.agent.ToProperties(propOpts)
	if err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		s.emit("error", err)
		return "", err
	}

	req := &Agora.StartAgentsRequest{
		Appid:      s.appID,
		Name:       s.name,
		Properties: properties,
	}

	reqOpts := s.convoAIRequestOpts(ctx)
	resp, err := s.client.Start(ctx, req, reqOpts...)
	if err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		s.emit("error", err)
		return "", err
	}

	s.mu.Lock()
	if resp != nil && resp.AgentID != nil {
		s.agentID = *resp.AgentID
	}
	s.status = StatusRunning
	s.mu.Unlock()

	s.emit("started", map[string]string{"agent_id": s.agentID})
	return s.agentID, nil
}

func (s *AgentSession) Stop(ctx context.Context) error {
	s.mu.Lock()
	if s.status != StatusRunning {
		s.mu.Unlock()
		return fmt.Errorf("cannot stop session in %s state", s.status)
	}
	if s.agentID == "" {
		s.mu.Unlock()
		return fmt.Errorf("no agent ID available")
	}
	s.status = StatusStopping
	s.mu.Unlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	err := s.client.Stop(ctx, &Agora.StopAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	}, reqOpts...)
	if err != nil {
		// Handle 404 "task not found" gracefully — agent is already stopped
		var apiErr *core.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			s.mu.Lock()
			s.status = StatusStopped
			s.mu.Unlock()
			s.emit("stopped", map[string]string{"agent_id": s.agentID})
			return nil
		}
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		s.emit("error", err)
		return err
	}

	s.mu.Lock()
	s.status = StatusStopped
	s.mu.Unlock()
	s.emit("stopped", map[string]string{"agent_id": s.agentID})
	return nil
}

func (s *AgentSession) Say(ctx context.Context, text string, priority *Agora.SpeakAgentsRequestPriority, interruptable *bool) error {
	s.mu.RLock()
	if s.status != StatusRunning {
		s.mu.RUnlock()
		return fmt.Errorf("cannot say in %s state", s.status)
	}
	if s.agentID == "" {
		s.mu.RUnlock()
		return fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	req := &Agora.SpeakAgentsRequest{
		Appid:         s.appID,
		AgentID:       s.agentID,
		Text:          text,
		Priority:      priority,
		Interruptable: interruptable,
	}

	reqOpts := s.convoAIRequestOpts(ctx)
	_, err := s.client.Speak(ctx, req, reqOpts...)
	return err
}

func (s *AgentSession) Interrupt(ctx context.Context) error {
	s.mu.RLock()
	if s.status != StatusRunning {
		s.mu.RUnlock()
		return fmt.Errorf("cannot interrupt in %s state", s.status)
	}
	if s.agentID == "" {
		s.mu.RUnlock()
		return fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	_, err := s.client.Interrupt(ctx, &Agora.InterruptAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	}, reqOpts...)
	return err
}

func (s *AgentSession) Update(ctx context.Context, properties *Agora.UpdateAgentsRequestProperties) error {
	s.mu.RLock()
	if s.status != StatusRunning {
		s.mu.RUnlock()
		return fmt.Errorf("cannot update in %s state", s.status)
	}
	if s.agentID == "" {
		s.mu.RUnlock()
		return fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	_, err := s.client.Update(ctx, &Agora.UpdateAgentsRequest{
		Appid:      s.appID,
		AgentID:    s.agentID,
		Properties: properties,
	}, reqOpts...)
	return err
}

func (s *AgentSession) GetHistory(ctx context.Context) (*Agora.GetHistoryAgentsResponse, error) {
	s.mu.RLock()
	if s.agentID == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	return s.client.GetHistory(ctx, &Agora.GetHistoryAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	}, reqOpts...)
}

func (s *AgentSession) GetInfo(ctx context.Context) (*Agora.GetAgentsResponse, error) {
	s.mu.RLock()
	if s.agentID == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	return s.client.Get(ctx, &Agora.GetAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	}, reqOpts...)
}

func (s *AgentSession) On(event string, handler EventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[event] = append(s.handlers[event], handler)
}

func (s *AgentSession) emit(event string, data interface{}) {
	s.mu.RLock()
	handlers := s.handlers[event]
	s.mu.RUnlock()

	for _, h := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log and continue so a panicking handler does not prevent
					// remaining handlers or session lifecycle from completing.
					log.Printf("agentkit: recovered panic in %q event handler: %v", event, r)
				}
			}()
			h(data)
		}()
	}
}
