package vendors

const (
	HeyGenRequiredSampleRate     = SampleRate24kHz
	LiveAvatarRequiredSampleRate = SampleRate24kHz
	AkoolRequiredSampleRate      = SampleRate16kHz
)

type HeyGenAvatarOptions struct {
	APIKey              string
	Quality             string
	AgoraUID            string
	AgoraToken          string
	AvatarID            string
	Enable              *bool
	DisableIdleTimeout  *bool
	ActivityIdleTimeout *int
	AdditionalParams    map[string]interface{}
}

type HeyGenAvatar struct {
	options HeyGenAvatarOptions
}

func NewHeyGenAvatar(opts HeyGenAvatarOptions) *HeyGenAvatar {
	if opts.APIKey == "" {
		panic("HeyGenAvatar requires APIKey")
	}
	if opts.Quality == "" {
		panic("HeyGenAvatar requires Quality (low, medium, or high)")
	}
	if opts.Quality != "low" && opts.Quality != "medium" && opts.Quality != "high" {
		panic("HeyGenAvatar Quality must be one of: low, medium, high")
	}
	if opts.AgoraUID == "" {
		panic("HeyGenAvatar requires AgoraUID")
	}
	return &HeyGenAvatar{options: opts}
}

func (h *HeyGenAvatar) RequiredSampleRate() SampleRate {
	return HeyGenRequiredSampleRate
}

func (h *HeyGenAvatar) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key":   h.options.APIKey,
		"quality":   h.options.Quality,
		"agora_uid": h.options.AgoraUID,
	}
	for k, v := range h.options.AdditionalParams {
		params[k] = v
	}
	if h.options.AgoraToken != "" {
		params["agora_token"] = h.options.AgoraToken
	}
	if h.options.AvatarID != "" {
		params["avatar_id"] = h.options.AvatarID
	}
	if h.options.DisableIdleTimeout != nil {
		params["disable_idle_timeout"] = *h.options.DisableIdleTimeout
	}
	if h.options.ActivityIdleTimeout != nil {
		params["activity_idle_timeout"] = *h.options.ActivityIdleTimeout
	}

	enable := true
	if h.options.Enable != nil {
		enable = *h.options.Enable
	}
	return map[string]interface{}{
		"enable": enable,
		"vendor": "heygen",
		"params": params,
	}
}

type AkoolAvatarOptions struct {
	APIKey           string
	AvatarID         string
	Enable           *bool
	AdditionalParams map[string]interface{}
}

type AkoolAvatar struct {
	options AkoolAvatarOptions
}

func NewAkoolAvatar(opts AkoolAvatarOptions) *AkoolAvatar {
	if opts.APIKey == "" {
		panic("AkoolAvatar requires APIKey")
	}
	return &AkoolAvatar{options: opts}
}

func (a *AkoolAvatar) RequiredSampleRate() SampleRate {
	return AkoolRequiredSampleRate
}

func (a *AkoolAvatar) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": a.options.APIKey,
	}
	for k, v := range a.options.AdditionalParams {
		params[k] = v
	}
	if a.options.AvatarID != "" {
		params["avatar_id"] = a.options.AvatarID
	}
	enable := true
	if a.options.Enable != nil {
		enable = *a.options.Enable
	}
	return map[string]interface{}{
		"enable": enable,
		"vendor": "akool",
		"params": params,
	}
}

type LiveAvatarAvatarOptions = HeyGenAvatarOptions

type LiveAvatarAvatar struct {
	options LiveAvatarAvatarOptions
}

func NewLiveAvatarAvatar(opts LiveAvatarAvatarOptions) *LiveAvatarAvatar {
	if opts.APIKey == "" {
		panic("LiveAvatarAvatar requires APIKey")
	}
	if opts.Quality == "" {
		panic("LiveAvatarAvatar requires Quality (low, medium, or high)")
	}
	if opts.Quality != "low" && opts.Quality != "medium" && opts.Quality != "high" {
		panic("LiveAvatarAvatar Quality must be one of: low, medium, high")
	}
	if opts.AgoraUID == "" {
		panic("LiveAvatarAvatar requires AgoraUID")
	}
	return &LiveAvatarAvatar{options: opts}
}

func (l *LiveAvatarAvatar) RequiredSampleRate() SampleRate {
	return LiveAvatarRequiredSampleRate
}

func (l *LiveAvatarAvatar) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key":   l.options.APIKey,
		"quality":   l.options.Quality,
		"agora_uid": l.options.AgoraUID,
	}
	for k, v := range l.options.AdditionalParams {
		params[k] = v
	}
	if l.options.AgoraToken != "" {
		params["agora_token"] = l.options.AgoraToken
	}
	if l.options.AvatarID != "" {
		params["avatar_id"] = l.options.AvatarID
	}
	if l.options.DisableIdleTimeout != nil {
		params["disable_idle_timeout"] = *l.options.DisableIdleTimeout
	}
	if l.options.ActivityIdleTimeout != nil {
		params["activity_idle_timeout"] = *l.options.ActivityIdleTimeout
	}

	enable := true
	if l.options.Enable != nil {
		enable = *l.options.Enable
	}
	return map[string]interface{}{
		"enable": enable,
		"vendor": "liveavatar",
		"params": params,
	}
}

type AnamAvatarOptions struct {
	APIKey           string
	PersonaID        string
	Enable           *bool
	AdditionalParams map[string]interface{}
}

type AnamAvatar struct {
	options AnamAvatarOptions
}

func NewAnamAvatar(opts AnamAvatarOptions) *AnamAvatar {
	if opts.APIKey == "" {
		panic("AnamAvatar requires APIKey")
	}
	return &AnamAvatar{options: opts}
}

func (a *AnamAvatar) RequiredSampleRate() SampleRate {
	return 0
}

func (a *AnamAvatar) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": a.options.APIKey,
	}
	for k, v := range a.options.AdditionalParams {
		params[k] = v
	}
	if a.options.PersonaID != "" {
		params["persona_id"] = a.options.PersonaID
	}
	enable := true
	if a.options.Enable != nil {
		enable = *a.options.Enable
	}
	return map[string]interface{}{
		"enable": enable,
		"vendor": "anam",
		"params": params,
	}
}
