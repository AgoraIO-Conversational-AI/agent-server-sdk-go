package agentkit

import (
	"encoding/json"
	"fmt"
	"time"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
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

// =============================================================================
// Top-level configuration aliases
// =============================================================================

type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type SalMode = Agora.StartAgentsRequestPropertiesSalSalMode
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters
type GeofenceConfig = Agora.StartAgentsRequestPropertiesGeofence
type RtcConfig = Agora.StartAgentsRequestPropertiesRtc
type FillerWordsConfig = Agora.StartAgentsRequestPropertiesFillerWords

// =============================================================================
// SOS/EOS turn detection aliases (preferred — replaces deprecated types below)
// =============================================================================

// TurnDetectionNestedConfig is the detailed nested config within TurnDetectionConfig.Config.
type TurnDetectionNestedConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfig

// StartOfSpeechConfig configures when the agent detects the start of a user's speech.
type StartOfSpeechConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeech

// StartOfSpeechMode is the detection mode: "vad" | "keywords" | "disabled".
type StartOfSpeechMode = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeechMode

// StartOfSpeechVadConfig holds VAD settings for SoS detection.
type StartOfSpeechVadConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeechVadConfig

// StartOfSpeechKeywordsConfig holds keyword-trigger settings for SoS detection.
type StartOfSpeechKeywordsConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeechKeywordsConfig

// StartOfSpeechDisabledConfig holds settings when SoS detection is disabled.
type StartOfSpeechDisabledConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeechDisabledConfig

// StartOfSpeechDisabledConfigStrategy is the voice processing strategy when SoS is disabled: "append" | "ignored".
type StartOfSpeechDisabledConfigStrategy = Agora.StartAgentsRequestPropertiesTurnDetectionConfigStartOfSpeechDisabledConfigStrategy

// EndOfSpeechConfig configures when the agent detects the end of a user's speech.
type EndOfSpeechConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeech

// EndOfSpeechMode is the detection mode: "vad" | "semantic".
type EndOfSpeechMode = Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeechMode

// EndOfSpeechVadConfig holds VAD settings for EoS detection.
type EndOfSpeechVadConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeechVadConfig

// EndOfSpeechSemanticConfig holds semantic model settings for EoS detection.
type EndOfSpeechSemanticConfig = Agora.StartAgentsRequestPropertiesTurnDetectionConfigEndOfSpeechSemanticConfig

// =============================================================================
// Deprecated turn detection aliases
// =============================================================================

// Deprecated: Use TurnDetectionConfig with TurnDetectionNestedConfig.StartOfSpeech and
// TurnDetectionNestedConfig.EndOfSpeech instead. The Type field and agora_vad / server_vad /
// semantic_vad values are being removed in a future release.
type TurnDetectionType = Agora.StartAgentsRequestPropertiesTurnDetectionType

// Deprecated: Use StartOfSpeechConfig with Mode "vad" | "keywords" | "disabled" and the
// corresponding VadConfig, KeywordsConfig, or DisabledConfig instead.
type InterruptMode = Agora.StartAgentsRequestPropertiesTurnDetectionInterruptMode

// Deprecated: Only applies to server_vad / semantic_vad modes with OpenAI Realtime API (MLLM).
// Has no equivalent in the standard ASR + LLM + TTS pipeline.
type Eagerness = Agora.StartAgentsRequestPropertiesTurnDetectionEagerness

// =============================================================================
// LLM sub-type aliases
// =============================================================================

// LlmGreetingConfigs configures how greeting messages are broadcast when multiple
// remote users are in the channel (llm.greeting_configs).
type LlmGreetingConfigs = Agora.StartAgentsRequestPropertiesLlmGreetingConfigs

// LlmGreetingConfigsMode is the greeting broadcast mode: "single_every" | "single_first".
type LlmGreetingConfigsMode = Agora.StartAgentsRequestPropertiesLlmGreetingConfigsMode

// McpServersItem is a single MCP server config entry (llm.mcp_servers[]).
type McpServersItem = Agora.StartAgentsRequestPropertiesLlmMcpServersItem

// =============================================================================
// Parameters (SessionParams) sub-type aliases
// =============================================================================

// SilenceConfig configures agent behaviour during user silence (parameters.silence_config).
type SilenceConfig = Agora.StartAgentsRequestPropertiesParametersSilenceConfig

// SilenceAction is the action taken after the silence timeout elapses.
type SilenceAction = Agora.StartAgentsRequestPropertiesParametersSilenceConfigAction

// FarewellConfig configures graceful hang-up behaviour (parameters.farewell_config).
type FarewellConfig = Agora.StartAgentsRequestPropertiesParametersFarewellConfig

// ParametersDataChannel is the agent data-transmission channel: "rtm" | "datastream".
type ParametersDataChannel = Agora.StartAgentsRequestPropertiesParametersDataChannel

// =============================================================================
// Geofence sub-type aliases
// =============================================================================

// GeofenceArea is an allowed geographic region for server access.
type GeofenceArea = Agora.StartAgentsRequestPropertiesGeofenceArea

// GeofenceExcludeArea is a geographic region to exclude when Area is "GLOBAL".
type GeofenceExcludeArea = Agora.StartAgentsRequestPropertiesGeofenceExcludeArea

// =============================================================================
// Concrete API payload config aliases (for constructing or inspecting ToProperties output)
// =============================================================================

// LlmConfig is the concrete LLM configuration payload (start_agents_request_properties.llm).
type LlmConfig = Agora.StartAgentsRequestPropertiesLlm

// MllmConfig is the concrete MLLM configuration payload (start_agents_request_properties.mllm).
type MllmConfig = Agora.StartAgentsRequestPropertiesMllm

// AsrConfig is the concrete ASR/STT configuration payload (start_agents_request_properties.asr).
type AsrConfig = Agora.StartAgentsRequestPropertiesAsr

// TtsConfig is the concrete TTS configuration payload (start_agents_request_properties.tts).
type TtsConfig = Agora.Tts

// AvatarConfig is the concrete Avatar configuration payload (start_agents_request_properties.avatar).
type AvatarConfig = Agora.StartAgentsRequestPropertiesAvatar

// =============================================================================
// FillerWords sub-type aliases
// =============================================================================

// FillerWordsTrigger configures when filler words are played (filler_words.trigger).
type FillerWordsTrigger = Agora.StartAgentsRequestPropertiesFillerWordsTrigger

// FillerWordsTriggerFixedTimeConfig holds the fixed-time trigger threshold (trigger.fixed_time_config).
type FillerWordsTriggerFixedTimeConfig = Agora.StartAgentsRequestPropertiesFillerWordsTriggerFixedTimeConfig

// FillerWordsContent configures the source and selection of filler words (filler_words.content).
type FillerWordsContent = Agora.StartAgentsRequestPropertiesFillerWordsContent

// FillerWordsContentStaticConfig configures a static list of filler words (content.static_config).
type FillerWordsContentStaticConfig = Agora.StartAgentsRequestPropertiesFillerWordsContentStaticConfig

// FillerWordsContentSelectionRule is the filler word selection rule: "shuffle" | "round_robin".
type FillerWordsContentSelectionRule = Agora.StartAgentsRequestPropertiesFillerWordsContentStaticConfigSelectionRule

type Agent struct {
	name                     string
	instructions             string
	greeting                 string
	failureMessage           string
	maxHistory               *int
	llm                      map[string]interface{}
	tts                      map[string]interface{}
	stt                      map[string]interface{}
	mllm                     map[string]interface{}
	ttsSampleRate            *vendors.SampleRate
	avatar                   map[string]interface{}
	avatarRequiredSampleRate *vendors.SampleRate
	turnDetection            *TurnDetectionConfig
	sal                      *SalConfig
	advancedFeatures         *AdvancedFeatures
	parameters               *SessionParams
	geofence                 *GeofenceConfig
	labels                   map[string]string
	rtc                      *RtcConfig
	fillerWords              *FillerWordsConfig
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

func (a *Agent) WithSal(sal *SalConfig) *Agent {
	clone := a.clone()
	clone.sal = sal
	return clone
}

func (a *Agent) WithAdvancedFeatures(af *AdvancedFeatures) *Agent {
	clone := a.clone()
	clone.advancedFeatures = af
	return clone
}

func (a *Agent) WithParameters(params *SessionParams) *Agent {
	clone := a.clone()
	clone.parameters = params
	return clone
}

func (a *Agent) WithFailureMessage(msg string) *Agent {
	clone := a.clone()
	clone.failureMessage = msg
	return clone
}

func (a *Agent) WithMaxHistory(n int) *Agent {
	clone := a.clone()
	clone.maxHistory = &n
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

func (a *Agent) Name() string                                  { return a.name }
func (a *Agent) Instructions() string                          { return a.instructions }
func (a *Agent) Greeting() string                              { return a.greeting }
func (a *Agent) LlmConfig() map[string]interface{}             { return a.llm }
func (a *Agent) TtsConfig() map[string]interface{}             { return a.tts }
func (a *Agent) SttConfig() map[string]interface{}             { return a.stt }
func (a *Agent) MllmConfig() map[string]interface{}            { return a.mllm }
func (a *Agent) TtsSampleRate() *vendors.SampleRate            { return a.ttsSampleRate }
func (a *Agent) AvatarRequiredSampleRate() *vendors.SampleRate { return a.avatarRequiredSampleRate }
func (a *Agent) FailureMessage() string                        { return a.failureMessage }
func (a *Agent) MaxHistory() *int                              { return a.maxHistory }
func (a *Agent) Avatar() map[string]interface{}                { return a.avatar }
func (a *Agent) TurnDetection() *TurnDetectionConfig           { return a.turnDetection }
func (a *Agent) Sal() *SalConfig                               { return a.sal }
func (a *Agent) AdvancedFeatures() *AdvancedFeatures           { return a.advancedFeatures }
func (a *Agent) Parameters() *SessionParams                    { return a.parameters }
func (a *Agent) Geofence() *GeofenceConfig                     { return a.geofence }
func (a *Agent) Labels() map[string]string                     { return a.labels }
func (a *Agent) Rtc() *RtcConfig                               { return a.rtc }
func (a *Agent) FillerWords() *FillerWordsConfig               { return a.fillerWords }

type CreateSessionOptions struct {
	Name            string
	Channel         string
	Token           string
	AgentUID        string
	RemoteUIDs      []string
	IdleTimeout     *int
	EnableStringUID *bool
	ExpiresIn       int
	Preset          []string
	PipelineID      string
	Debug           bool
	Warn            func(string)
}

func (a *Agent) CreateSession(client *AgoraClient, opts CreateSessionOptions) *AgentSession {
	name := opts.Name
	if name == "" {
		if a.name != "" {
			name = a.name
		} else {
			name = fmt.Sprintf("agent-%d", time.Now().UnixMilli())
		}
	}

	return NewAgentSession(AgentSessionOptions{
		Client:                   client.Agents,
		Agent:                    a,
		AppID:                    client.AppID,
		AppCertificate:           client.AppCertificate,
		Name:                     name,
		Channel:                  opts.Channel,
		Token:                    opts.Token,
		AgentUID:                 opts.AgentUID,
		RemoteUIDs:               opts.RemoteUIDs,
		IdleTimeout:              opts.IdleTimeout,
		EnableStringUID:          opts.EnableStringUID,
		ExpiresIn:                opts.ExpiresIn,
		UseAppCredentialsForREST: client.AuthMode == AuthModeAppCredentials,
		Preset:                   opts.Preset,
		PipelineID:               opts.PipelineID,
		Debug:                    opts.Debug,
		Warn:                     opts.Warn,
	})
}

func (a *Agent) ToProperties(opts ToPropertiesOptions) (*Agora.StartAgentsRequestProperties, error) {
	token := opts.Token
	if token == "" {
		if opts.AppID == "" || opts.AppCertificate == "" {
			return nil, fmt.Errorf("either token or app_id+app_certificate must be provided")
		}
		expiry := opts.ExpiresIn
		if expiry > 0 {
			var err error
			expiry, err = ValidateExpiresIn(expiry)
			if err != nil {
				return nil, fmt.Errorf("invalid expiresIn: %w", err)
			}
		}
		// Use GenerateConvoAIToken (RTC + RTM) so the token works whether or
		// not the caller enables advanced_features.enable_rtm.
		var err error
		token, err = GenerateConvoAIToken(GenerateConvoAITokenOptions{
			AppID:          opts.AppID,
			AppCertificate: opts.AppCertificate,
			ChannelName:    opts.Channel,
			Account:        opts.AgentUID,
			TokenExpire:    expiry,
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
		mllmConfig := make(map[string]interface{})
		for k, v := range a.mllm {
			mllmConfig[k] = v
		}
		if a.greeting != "" {
			if _, exists := mllmConfig["greeting_message"]; !exists {
				mllmConfig["greeting_message"] = a.greeting
			}
		}
		if a.failureMessage != "" {
			if _, exists := mllmConfig["failure_message"]; !exists {
				mllmConfig["failure_message"] = a.failureMessage
			}
		}
		if a.maxHistory != nil {
			if _, exists := mllmConfig["max_history"]; !exists {
				mllmConfig["max_history"] = *a.maxHistory
			}
		}
		var mllm Agora.StartAgentsRequestPropertiesMllm
		if err := mapToStruct(mllmConfig, &mllm); err != nil {
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
		params := *a.parameters
		props.Parameters = &params
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

	if a.advancedFeatures != nil && a.advancedFeatures.EnableRtm != nil && *a.advancedFeatures.EnableRtm {
		if props.Parameters == nil {
			props.Parameters = &Agora.StartAgentsRequestPropertiesParameters{}
		}
		if props.Parameters.DataChannel == nil {
			rtm := Agora.StartAgentsRequestPropertiesParametersDataChannel("rtm")
			props.Parameters.DataChannel = &rtm
		}
	}

	if !opts.SkipVendorValidation {
		if a.tts == nil {
			return nil, fmt.Errorf("TTS configuration is required; use WithTts() to set it")
		}
		if a.llm == nil {
			return nil, fmt.Errorf("LLM configuration is required; use WithLlm() to set it")
		}
	}

	if a.llm != nil {
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
	}

	if a.tts != nil {
		var tts Agora.Tts
		if err := mapToStruct(a.tts, &tts); err != nil {
			return nil, fmt.Errorf("failed to convert TTS config: %w", err)
		}
		props.Tts = &tts
	}

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
	Channel        string
	AgentUID       string
	RemoteUIDs     []string
	Token          string
	AppID          string
	AppCertificate string
	// ExpiresIn is the token lifetime in seconds (default: 86400 = 24 hours, Agora maximum).
	// Valid range: 1–86400. Use ExpiresInHours() / ExpiresInMinutes() for clarity.
	ExpiresIn            int
	IdleTimeout          *int
	EnableStringUID      *bool
	SkipVendorValidation bool
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
