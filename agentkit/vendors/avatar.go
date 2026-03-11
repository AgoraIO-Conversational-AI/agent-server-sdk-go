package vendors

const (
	HeyGenRequiredSampleRate = SampleRate24kHz
	AkoolRequiredSampleRate  = SampleRate16kHz
)

type HeyGenAvatarOptions struct {
	APIKey     string
	Quality    string
	AgoraUID   string
	AvatarName string
	VoiceID    string
	Language   string
	Version    string
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
	if h.options.AvatarName != "" {
		params["avatar_name"] = h.options.AvatarName
	}
	if h.options.VoiceID != "" {
		params["voice_id"] = h.options.VoiceID
	}
	if h.options.Language != "" {
		params["language"] = h.options.Language
	}
	if h.options.Version != "" {
		params["version"] = h.options.Version
	}

	return map[string]interface{}{
		"vendor": "heygen",
		"params": params,
	}
}

type AkoolAvatarOptions struct {
	APIKey string
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
	return map[string]interface{}{
		"vendor": "akool",
		"params": map[string]interface{}{
			"api_key": a.options.APIKey,
		},
	}
}
