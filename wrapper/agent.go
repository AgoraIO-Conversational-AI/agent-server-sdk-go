package wrapper

import (
	"fmt"
	"strings"

	Agora "github.com/fern-demo/agoraio-go-sdk"
)

var vendorURLs = map[string]string{
	"openai":    "https://api.openai.com/v1/chat/completions",
	"anthropic": "https://api.anthropic.com/v1/messages",
	"azure":     "https://YOUR_RESOURCE.openai.azure.com/openai/deployments/YOUR_DEPLOYMENT/chat/completions",
	"gemini":    "https://generativelanguage.googleapis.com/v1beta/models",
}

var sttVendorMap = map[string]Agora.StartAgentsRequestPropertiesAsrVendor{
	"ares":        "ares",
	"microsoft":   "microsoft",
	"deepgram":    "deepgram",
	"openai":      "openai",
	"google":      "google",
	"amazon":      "amazon",
	"assemblyai":  "assemblyai",
	"speechmatics": "speechmatics",
}

type LlmConfig = Agora.StartAgentsRequestPropertiesLlm
type SttConfig = Agora.StartAgentsRequestPropertiesAsr
type TtsConfig = Agora.Tts
type MllmConfig = Agora.StartAgentsRequestPropertiesMllm
type TurnDetectionConfig = Agora.StartAgentsRequestPropertiesTurnDetection
type SalConfig = Agora.StartAgentsRequestPropertiesSal
type AvatarConfig = Agora.StartAgentsRequestPropertiesAvatar
type AdvancedFeatures = Agora.StartAgentsRequestPropertiesAdvancedFeatures
type SessionParams = Agora.StartAgentsRequestPropertiesParameters

type Agent struct {
	name             string
	instructions     string
	greeting         string
	failureMessage   string
	maxHistory       *int
	llm              *LlmConfig
	tts              *TtsConfig
	stt              *SttConfig
	mllm             *MllmConfig
	turnDetection    *TurnDetectionConfig
	sal              *SalConfig
	avatar           *AvatarConfig
	advancedFeatures *AdvancedFeatures
	parameters       *SessionParams
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

func WithLlmConfig(llm *LlmConfig) AgentOption {
	return func(a *Agent) {
		a.llm = llm
	}
}

func WithLlmShorthand(shorthand string) AgentOption {
	return func(a *Agent) {
		a.llm = parseLlmShorthand(shorthand)
	}
}

func WithTtsConfig(tts *TtsConfig) AgentOption {
	return func(a *Agent) {
		a.tts = tts
	}
}

func WithSttConfig(stt *SttConfig) AgentOption {
	return func(a *Agent) {
		a.stt = stt
	}
}

func WithSttShorthand(shorthand string) AgentOption {
	return func(a *Agent) {
		a.stt = parseSttShorthand(shorthand)
	}
}

func WithMllmConfig(mllm *MllmConfig) AgentOption {
	return func(a *Agent) {
		a.mllm = mllm
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

func WithAvatarConfig(avatar *AvatarConfig) AgentOption {
	return func(a *Agent) {
		a.avatar = avatar
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

func (a *Agent) SetLlm(llm *LlmConfig) *Agent {
	clone := a.clone()
	clone.llm = llm
	return clone
}

func (a *Agent) SetLlmShorthand(shorthand string) *Agent {
	clone := a.clone()
	clone.llm = parseLlmShorthand(shorthand)
	return clone
}

func (a *Agent) SetTts(tts *TtsConfig) *Agent {
	clone := a.clone()
	clone.tts = tts
	return clone
}

func (a *Agent) SetStt(stt *SttConfig) *Agent {
	clone := a.clone()
	clone.stt = stt
	return clone
}

func (a *Agent) SetSttShorthand(shorthand string) *Agent {
	clone := a.clone()
	clone.stt = parseSttShorthand(shorthand)
	return clone
}

func (a *Agent) SetMllm(mllm *MllmConfig) *Agent {
	clone := a.clone()
	clone.mllm = mllm
	return clone
}

func (a *Agent) SetTurnDetection(td *TurnDetectionConfig) *Agent {
	clone := a.clone()
	clone.turnDetection = td
	return clone
}

func (a *Agent) SetInstructions(instructions string) *Agent {
	clone := a.clone()
	clone.instructions = instructions
	return clone
}

func (a *Agent) SetGreeting(greeting string) *Agent {
	clone := a.clone()
	clone.greeting = greeting
	return clone
}

func (a *Agent) SetName(name string) *Agent {
	clone := a.clone()
	clone.name = name
	return clone
}

func (a *Agent) Name() string           { return a.name }
func (a *Agent) Instructions() string    { return a.instructions }
func (a *Agent) Greeting() string        { return a.greeting }
func (a *Agent) Llm() *LlmConfig        { return a.llm }
func (a *Agent) Tts() *TtsConfig        { return a.tts }
func (a *Agent) Stt() *SttConfig        { return a.stt }
func (a *Agent) Mllm() *MllmConfig      { return a.mllm }

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
		props.Mllm = a.mllm
	}
	if a.turnDetection != nil {
		props.TurnDetection = a.turnDetection
	}
	if a.sal != nil {
		props.Sal = a.sal
	}
	if a.avatar != nil {
		props.Avatar = a.avatar
	}
	if a.advancedFeatures != nil {
		props.AdvancedFeatures = a.advancedFeatures
	}
	if a.parameters != nil {
		props.Parameters = a.parameters
	}

	isMllmMode := a.advancedFeatures != nil && a.advancedFeatures.EnableMllm != nil && *a.advancedFeatures.EnableMllm
	if isMllmMode {
		return props, nil
	}

	if a.tts == nil {
		return nil, fmt.Errorf("TTS configuration is required; use WithTtsConfig() to set it")
	}
	if a.llm == nil {
		return nil, fmt.Errorf("LLM configuration is required; use WithLlmConfig() or WithLlmShorthand() to set it")
	}

	llm := *a.llm
	if a.instructions != "" {
		llm.SystemMessages = []map[string]interface{}{
			{"role": "system", "content": a.instructions},
		}
	}
	if a.greeting != "" && llm.GreetingMessage == nil {
		llm.GreetingMessage = &a.greeting
	}
	if a.failureMessage != "" && llm.FailureMessage == nil {
		llm.FailureMessage = &a.failureMessage
	}
	if a.maxHistory != nil && llm.MaxHistory == nil {
		llm.MaxHistory = a.maxHistory
	}

	props.Llm = &llm
	props.Tts = a.tts
	if a.stt != nil {
		props.Asr = a.stt
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
	return &clone
}

func parseLlmShorthand(shorthand string) *LlmConfig {
	parts := strings.SplitN(shorthand, "/", 2)
	vendor := strings.ToLower(parts[0])
	var model string
	if len(parts) > 1 {
		model = parts[1]
	}

	url, ok := vendorURLs[vendor]
	if !ok {
		url = vendorURLs["openai"]
	}

	var style *Agora.StartAgentsRequestPropertiesLlmStyle
	switch vendor {
	case "gemini":
		s := Agora.StartAgentsRequestPropertiesLlmStyle("gemini")
		style = &s
	case "anthropic":
		s := Agora.StartAgentsRequestPropertiesLlmStyle("anthropic")
		style = &s
	default:
		s := Agora.StartAgentsRequestPropertiesLlmStyle("openai")
		style = &s
	}

	llm := &LlmConfig{
		URL:   url,
		Style: style,
	}

	if vendor == "azure" {
		v := "azure"
		llm.Vendor = &v
	}

	if model != "" {
		llm.Params = map[string]interface{}{"model": model}
	}

	return llm
}

func parseSttShorthand(shorthand string) *SttConfig {
	parts := strings.SplitN(shorthand, "/", 2)
	vendor := strings.ToLower(parts[0])

	v, ok := sttVendorMap[vendor]
	if !ok {
		v = "deepgram"
	}

	return &SttConfig{
		Vendor: &v,
	}
}
