package wrapper

import (
	"context"
	"fmt"
	"sync"
	"time"

	Agora "github.com/fern-demo/agoraio-go-sdk"
	"github.com/fern-demo/agoraio-go-sdk/agents"
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
	client         *agents.Client
	agent          *Agent
	appID          string
	appCertificate string
	name           string
	channel        string
	token          string
	agentUID       string
	remoteUIDs     []string
	idleTimeout    *int
	enableStringUID *bool

	agentID  string
	status   SessionStatus
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

type AgentSessionOptions struct {
	Client         *agents.Client
	Agent          *Agent
	AppID          string
	AppCertificate string
	Name           string
	Channel        string
	Token          string
	AgentUID       string
	RemoteUIDs     []string
	IdleTimeout    *int
	EnableStringUID *bool
}

func NewAgentSession(opts AgentSessionOptions) *AgentSession {
	name := opts.Name
	if name == "" {
		name = fmt.Sprintf("agent-%d", time.Now().Unix())
	}

	return &AgentSession{
		client:         opts.Client,
		agent:          opts.Agent,
		appID:          opts.AppID,
		appCertificate: opts.AppCertificate,
		name:           name,
		channel:        opts.Channel,
		token:          opts.Token,
		agentUID:       opts.AgentUID,
		remoteUIDs:     opts.RemoteUIDs,
		idleTimeout:    opts.IdleTimeout,
		enableStringUID: opts.EnableStringUID,
		status:         StatusIdle,
		handlers:       make(map[string][]EventHandler),
	}
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

	if s.agent.avatar != nil {
		vendor := ""
		if s.agent.avatar.Vendor != nil {
			vendor = string(*s.agent.avatar.Vendor)
		}
		if IsHeyGenAvatar(vendor) || IsAkoolAvatar(vendor) {
			if err := ValidateAvatarConfig(vendor, s.agent.avatar.Params); err != nil {
				s.mu.Unlock()
				return "", err
			}
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

	resp, err := s.client.Start(ctx, req)
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

	err := s.client.Stop(ctx, &Agora.StopAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	})
	if err != nil {
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

	_, err := s.client.Speak(ctx, req)
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

	_, err := s.client.Interrupt(ctx, &Agora.InterruptAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	})
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

	_, err := s.client.Update(ctx, &Agora.UpdateAgentsRequest{
		Appid:      s.appID,
		AgentID:    s.agentID,
		Properties: properties,
	})
	return err
}

func (s *AgentSession) GetHistory(ctx context.Context) (*Agora.GetHistoryAgentsResponse, error) {
	s.mu.RLock()
	if s.agentID == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	return s.client.GetHistory(ctx, &Agora.GetHistoryAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	})
}

func (s *AgentSession) GetInfo(ctx context.Context) (*Agora.GetAgentsResponse, error) {
	s.mu.RLock()
	if s.agentID == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	return s.client.Get(ctx, &Agora.GetAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	})
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
			defer func() { recover() }()
			h(data)
		}()
	}
}
