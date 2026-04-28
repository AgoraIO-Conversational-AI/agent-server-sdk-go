package agentkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
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
	preset          []string
	pipelineID      string
	debug           bool
	warn            func(string)

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
	Preset                   []string
	PipelineID               string
	Debug                    bool
	Warn                     func(string)
}

func NewAgentSession(opts AgentSessionOptions) *AgentSession {
	name := opts.Name
	if name == "" {
		name = fmt.Sprintf("agent-%d", time.Now().UnixMilli())
	}

	return &AgentSession{
		client:          opts.Client,
		agent:           opts.Agent,
		appID:           opts.AppID,
		appCertificate:  opts.AppCertificate,
		name:            name,
		channel:         opts.Channel,
		token:           opts.Token,
		agentUID:        opts.AgentUID,
		remoteUIDs:      opts.RemoteUIDs,
		idleTimeout:     opts.IdleTimeout,
		enableStringUID: opts.EnableStringUID,
		expiresIn:       opts.ExpiresIn,
		useAppCredsREST: opts.UseAppCredentialsForREST,
		preset:          append([]string(nil), opts.Preset...),
		pipelineID:      opts.PipelineID,
		debug:           opts.Debug,
		warn:            opts.Warn,
		status:          StatusIdle,
		handlers:        make(map[string][]EventHandler),
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

func (s *AgentSession) warnf(message string) {
	if s.warn != nil {
		s.warn(message)
		return
	}
	log.Printf("agentkit: %s", message)
}

func (s *AgentSession) validateAvatarConfig() error {
	if s.agent == nil || s.agent.avatar == nil {
		return nil
	}

	vendor, _ := s.agent.avatar["vendor"].(string)
	if vendor == "" {
		return nil
	}
	params, _ := s.agent.avatar["params"].(map[string]interface{})
	if err := ValidateAvatarConfig(vendor, params); err != nil {
		return err
	}

	if s.agent.ttsSampleRate != nil {
		if err := ValidateTtsSampleRate(vendor, int(*s.agent.ttsSampleRate)); err != nil {
			return err
		}
		return nil
	}

	switch vendor {
	case "heygen":
		s.warnf("Warning: HeyGen avatar detected but TTS sample_rate is not explicitly set. HeyGen requires 24,000 Hz. Please ensure your TTS provider is configured for 24kHz.")
	case "liveavatar":
		s.warnf("Warning: LiveAvatar avatar detected but TTS sample_rate is not explicitly set. LiveAvatar requires 24,000 Hz. Please ensure your TTS provider is configured for 24kHz.")
	case "akool":
		s.warnf("Warning: Akool avatar detected but TTS sample_rate is not explicitly set. Akool requires 16,000 Hz. Please ensure your TTS provider is configured for 16kHz.")
	}
	return nil
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
	if err := s.validateAvatarConfig(); err != nil {
		s.mu.Unlock()
		return "", err
	}

	s.status = StatusStarting
	s.mu.Unlock()

	propOpts := ToPropertiesOptions{
		Channel:              s.channel,
		AgentUID:             s.agentUID,
		RemoteUIDs:           s.remoteUIDs,
		Token:                s.token,
		AppID:                s.appID,
		AppCertificate:       s.appCertificate,
		ExpiresIn:            s.expiresIn,
		IdleTimeout:          s.idleTimeout,
		EnableStringUID:      s.enableStringUID,
		SkipVendorValidation: len(s.preset) > 0 || s.pipelineID != "",
	}

	properties, err := s.agent.ToPropertiesMap(propOpts)
	if err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		s.emit("error", err)
		return "", err
	}

	resolvedPreset, resolvedProperties, err := ResolveSessionPresetsMap(s.preset, properties)
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
		Preset:     stringPtrOrNil(resolvedPreset),
		PipelineID: stringPtrOrNil(s.pipelineID),
	}
	if s.debug {
		debugPayload := map[string]interface{}{
			"name":       s.name,
			"properties": resolvedProperties,
		}
		if resolvedPreset != "" {
			debugPayload["preset"] = resolvedPreset
		}
		if s.pipelineID != "" {
			debugPayload["pipeline_id"] = s.pipelineID
		}
		if payload, err := json.Marshal(debugPayload); err == nil {
			log.Printf("[Agora Debug] Starting agent session: %s", payload)
		} else {
			s.warnf(fmt.Sprintf("debug logging failed to marshal start request: %v", err))
		}
	}

	reqOpts := s.convoAIRequestOpts(ctx)
	resp, err := s.client.StartWithPropertiesMap(ctx, req, resolvedProperties, reqOpts...)
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

func (s *AgentSession) GetTurns(ctx context.Context) (*Agora.GetTurnsAgentsResponse, error) {
	s.mu.RLock()
	if s.agentID == "" {
		s.mu.RUnlock()
		return nil, fmt.Errorf("no agent ID available")
	}
	s.mu.RUnlock()

	reqOpts := s.convoAIRequestOpts(ctx)
	return s.client.GetTurns(ctx, &Agora.GetTurnsAgentsRequest{
		Appid:   s.appID,
		AgentID: s.agentID,
	}, reqOpts...)
}

func (s *AgentSession) On(event string, handler EventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[event] = append(s.handlers[event], handler)
}

func (s *AgentSession) Off(event string, handler EventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	handlers := s.handlers[event]
	if len(handlers) == 0 {
		return
	}

	target := reflect.ValueOf(handler).Pointer()
	filtered := handlers[:0]
	for _, registered := range handlers {
		if reflect.ValueOf(registered).Pointer() == target {
			continue
		}
		filtered = append(filtered, registered)
	}
	if len(filtered) == 0 {
		delete(s.handlers, event)
		return
	}
	s.handlers[event] = filtered
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
					s.warnf(fmt.Sprintf("recovered panic in %q event handler: %v", event, r))
				}
			}()
			h(data)
		}()
	}
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
