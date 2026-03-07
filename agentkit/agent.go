package agentkit

import (
	"encoding/json"
	"fmt"

	Agora "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk"
	"github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/agentkit/vendors"
)

func mapToStruct(m map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal config map: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}
	return nil
}

type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters
type GeofenceConfig = Agora.StartAgentsRequestPropertiesGeofence
type RtcConfig = Agora.StartAgentsRequestPropertiesRtc
type FillerWordsConfig = Agora.StartAgentsRequestPropertiesFillerWords

type Agent struct {
	name                    string
	instructions            string
	greeting                string
	failureMessage          string
	maxHistory              *int
	llm                     map[string]interface{}
	tts                     map[string]interface{}
	stt                     map[string]interface{}
	mllm                    map[string]interface{}
	ttsSampleRate           *vendors.SampleRate
	avatar                  map[string]interface{}
	avatarRequiredSampleRate *vendors.SampleRate
	turnDetection           *TurnDetectionConfig
	sal                     *SalConfig
	advancedFeatures        *AdvancedFeatures
	parameters              *SessionParams
	geofence                *GeofenceConfig
	labels                  map[string]string
	rtc                     *RtcConfig
	fillerWords             *FillerWordsConfig
}

type AgentOption func(*Agent)

func NewAgent(opts ...AgentOption) *Agent {
	a := &Agent{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithName(name string) AgentOption {
	return func(a *Agent) {
		a.name = name
	}
}

func WithInstructions(instructions string) AgentOption {
	return func(a *Agent) {
		a.instructions = instructions
	}
}

func WithGreeting(greeting string) AgentOption {
	return func(a *Agent) {
		a.greeting = greeting
	}
}

func WithFailureMessage(msg string) AgentOption {
	return func(a *Agent) {
		a.failureMessage = msg
	}
}

func WithMaxHistory(n int) AgentOption {
	return func(a *Agent) {
		a.maxHistory = &n
	}
}

func WithTurnDetectionConfig(td *TurnDetectionConfig) AgentOption {
	return func(a *Agent) {
		a.turnDetection = td
	}
}

func WithSalConfig(sal *SalConfig) AgentOption {
	return func(a *Agent) {
		a.sal = sal
	}
}

func WithAdvancedFeatures(af *AdvancedFeatures) AgentOption {
	return func(a *Agent) {
		a.advancedFeatures = af
	}
}

func WithParameters(params *SessionParams) AgentOption {
	return func(a *Agent) {
		a.parameters = params
	}
}

func WithGeofence(gf *GeofenceConfig) AgentOption {
	return func(a *Agent) {
		a.geofence = gf
	}
}

func WithLabels(labels map[string]string) AgentOption {
	return func(a *Agent) {
		a.labels = labels
	}
}

func WithRtc(rtc *RtcConfig) AgentOption {
	return func(a *Agent) {
		a.rtc = rtc
	}
}

func WithFillerWords(fw *FillerWordsConfig) AgentOption {
	return func(a *Agent) {
		a.fillerWords = fw
	}
}

func (a *Agent) WithLlm(vendor vendors.LLM) *Agent {
	clone := a.clone()
	clone.llm = vendor.ToConfig()
	return clone
}

func (a *Agent) WithTts(vendor vendors.TTS) *Agent {
	clone := a.clone()
	clone.tts = vendor.ToConfig()
	clone.ttsSampleRate = vendor.GetSampleRate()
	// If an avatar is already set, verify the new TTS sample rate matches.
	// Mirrors the check in WithAvatar so both call orderings fail fast.
	if clone.avatarRequiredSampleRate != nil && clone.ttsSampleRate != nil {
		if *clone.ttsSampleRate != *clone.avatarRequiredSampleRate {
			panic(fmt.Sprintf(
				"TTS sample rate %d Hz is incompatible with the configured avatar, which requires %d Hz. "+
					"Please update your TTS sample_rate to %d.",
				int(*clone.ttsSampleRate), int(*clone.avatarRequiredSampleRate), int(*clone.avatarRequiredSampleRate),
			))
		}
	}
	return clone
}

func (a *Agent) WithStt(vendor vendors.STT) *Agent {
	clone := a.clone()
	clone.stt = vendor.ToConfig()
	return clone
}

func (a *Agent) WithMllm(vendor vendors.MLLM) *Agent {
	clone := a.clone()
	clone.mllm = vendor.ToConfig()
	return clone
}

func (a *Agent) WithAvatar(vendor vendors.Avatar) *Agent {
	requiredSR := vendor.RequiredSampleRate()
	// If a TTS is already set, verify sample rate compatibility now.
	// Mirrors the check in WithTts so both call orderings fail fast.
	// AgentSession.Start also validates as a final safety net.
	if a.ttsSampleRate != nil && *a.ttsSampleRate != requiredSR {
		panic(fmt.Sprintf(
			"Avatar requires TTS sample rate of %d Hz, but TTS is configured with %d Hz. "+
				"Please update your TTS sample_rate to %d.",
			int(requiredSR), int(*a.ttsSampleRate), int(requiredSR),
		))
	}
	clone := a.clone()
	clone.avatar = vendor.ToConfig()
	clone.avatarRequiredSampleRate = &requiredSR
	return clone
}

func (a *Agent) WithTurnDetection(td *TurnDetectionConfig) *Agent {
	clone := a.clone()
	clone.turnDetection = td
	return clone
}

func (a *Agent) WithInstructions(instructions string) *Agent {
	clone := a.clone()
	clone.instructions = instructions
	return clone
}

func (a *Agent) WithGreeting(greeting string) *Agent {
	clone := a.clone()
	clone.greeting = greeting
	return clone
}

func (a *Agent) WithName(name string) *Agent {
	clone := a.clone()
	clone.name = name
	return clone
}

func (a *Agent) WithGeofence(gf *GeofenceConfig) *Agent {
	clone := a.clone()
	clone.geofence = gf
	return clone
}

func (a *Agent) WithLabels(labels map[string]string) *Agent {
	clone := a.clone()
	clone.labels = labels
	return clone
}

func (a *Agent) WithRtc(rtc *RtcConfig) *Agent {
	clone := a.clone()
	clone.rtc = rtc
	return clone
}

func (a *Agent) WithFillerWords(fw *FillerWordsConfig) *Agent {
	clone := a.clone()
	clone.fillerWords = fw
	return clone
}

func (a *Agent) Name() string                        { return a.name }
func (a *Agent) Instructions() string                 { return a.instructions }
func (a *Agent) Greeting() string                     { return a.greeting }
func (a *Agent) LlmConfig() map[string]interface{}    { return a.llm }
func (a *Agent) TtsConfig() map[string]interface{}    { return a.tts }
func (a *Agent) SttConfig() map[string]interface{}    { return a.stt }
func (a *Agent) MllmConfig() map[string]interface{}   { return a.mllm }
func (a *Agent) TtsSampleRate() *vendors.SampleRate   { return a.ttsSampleRate }
func (a *Agent) AvatarRequiredSampleRate() *vendors.SampleRate { return a.avatarRequiredSampleRate }

func (a *Agent) ToProperties(opts ToPropertiesOptions) (*Agora.StartAgentsRequestProperties, error) {
	token := opts.Token
	if token == "" {
		if opts.AppID == "" || opts.AppCertificate == "" {
			return nil, fmt.Errorf("either token or app_id+app_certificate must be provided")
		}
		expiry := opts.TokenExpirySeconds
		if expiry <= 0 {
			expiry = DefaultExpirySeconds
		}
		var err error
		token, err = GenerateRtcToken(GenerateTokenOptions{
			AppID:          opts.AppID,
			AppCertificate: opts.AppCertificate,
			Channel:        opts.Channel,
			UID:            opts.UID,
			ExpirySeconds:  expiry,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}
	}

	props := &Agora.StartAgentsRequestProperties{
		Channel:       opts.Channel,
		Token:         token,
		AgentRtcUID:   opts.AgentUID,
		RemoteRtcUIDs: opts.RemoteUIDs,
	}

	if opts.IdleTimeout != nil {
		props.IdleTimeout = opts.IdleTimeout
	}
	if opts.EnableStringUID != nil {
		props.EnableStringUID = opts.EnableStringUID
	}
	if a.mllm != nil {
		var mllm Agora.StartAgentsRequestPropertiesMllm
		if err := mapToStruct(a.mllm, &mllm); err != nil {
			return nil, fmt.Errorf("failed to convert MLLM config: %w", err)
		}
		props.Mllm = &mllm
	}
	if a.turnDetection != nil {
		props.TurnDetection = a.turnDetection
	}
	if a.sal != nil {
		props.Sal = a.sal
	}
	if a.avatar != nil {
		var avatar Agora.StartAgentsRequestPropertiesAvatar
		if err := mapToStruct(a.avatar, &avatar); err != nil {
			return nil, fmt.Errorf("failed to convert avatar config: %w", err)
		}
		props.Avatar = &avatar
	}
	if a.advancedFeatures != nil {
		props.AdvancedFeatures = a.advancedFeatures
	}
	if a.parameters != nil {
		props.Parameters = a.parameters
	}
	if a.geofence != nil {
		props.Geofence = a.geofence
	}
	if len(a.labels) > 0 {
		props.Labels = a.labels
	}
	if a.rtc != nil {
		props.Rtc = a.rtc
	}
	if a.fillerWords != nil {
		props.FillerWords = a.fillerWords
	}

	isMllmMode := a.advancedFeatures != nil && a.advancedFeatures.EnableMllm != nil && *a.advancedFeatures.EnableMllm
	if isMllmMode {
		return props, nil
	}

	if a.tts == nil {
		return nil, fmt.Errorf("TTS configuration is required; use WithTts() to set it")
	}
	if a.llm == nil {
		return nil, fmt.Errorf("LLM configuration is required; use WithLlm() to set it")
	}

	llmConfig := make(map[string]interface{})
	for k, v := range a.llm {
		llmConfig[k] = v
	}
	if a.instructions != "" {
		llmConfig["system_messages"] = []map[string]interface{}{
			{"role": "system", "content": a.instructions},
		}
	}
	if a.greeting != "" {
		if _, exists := llmConfig["greeting_message"]; !exists {
			llmConfig["greeting_message"] = a.greeting
		}
	}
	if a.failureMessage != "" {
		if _, exists := llmConfig["failure_message"]; !exists {
			llmConfig["failure_message"] = a.failureMessage
		}
	}
	if a.maxHistory != nil {
		if _, exists := llmConfig["max_history"]; !exists {
			llmConfig["max_history"] = *a.maxHistory
		}
	}

	var llm Agora.StartAgentsRequestPropertiesLlm
	if err := mapToStruct(llmConfig, &llm); err != nil {
		return nil, fmt.Errorf("failed to convert LLM config: %w", err)
	}
	props.Llm = &llm

	var tts Agora.Tts
	if err := mapToStruct(a.tts, &tts); err != nil {
		return nil, fmt.Errorf("failed to convert TTS config: %w", err)
	}
	props.Tts = &tts

	if a.stt != nil {
		var stt Agora.StartAgentsRequestPropertiesAsr
		if err := mapToStruct(a.stt, &stt); err != nil {
			return nil, fmt.Errorf("failed to convert STT config: %w", err)
		}
		props.Asr = &stt
	}

	return props, nil
}

type ToPropertiesOptions struct {
	Channel            string
	AgentUID           string
	RemoteUIDs         []string
	Token              string
	AppID              string
	AppCertificate     string
	UID                uint32
	TokenExpirySeconds int
	IdleTimeout        *int
	EnableStringUID    *bool
}

func (a *Agent) clone() *Agent {
	clone := *a
	if a.labels != nil {
		clone.labels = make(map[string]string, len(a.labels))
		for k, v := range a.labels {
			clone.labels[k] = v
		}
	}
	return &clone
}
