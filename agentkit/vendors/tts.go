package vendors

import "strings"

type ElevenLabsTTSOptions struct {
	Key                      string
	ModelID                  string
	VoiceID                  string
	BaseURL                  string
	SampleRate               *SampleRate
	OptimizeStreamingLatency *int
	Stability                *float64
	SimilarityBoost          *float64
	Style                    *float64
	UseSpeakerBoost          *bool
	SkipPatterns             []int
}

type ElevenLabsTTS struct {
	options ElevenLabsTTSOptions
}

func NewElevenLabsTTS(opts ElevenLabsTTSOptions) *ElevenLabsTTS {
	if opts.Key == "" {
		panic("ElevenLabsTTS requires Key")
	}
	if opts.ModelID == "" {
		panic("ElevenLabsTTS requires ModelID")
	}
	if opts.VoiceID == "" {
		panic("ElevenLabsTTS requires VoiceID")
	}
	return &ElevenLabsTTS{options: opts}
}

func (e *ElevenLabsTTS) GetSampleRate() *SampleRate {
	return e.options.SampleRate
}

func (e *ElevenLabsTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":      e.options.Key,
		"model_id": e.options.ModelID,
		"voice_id": e.options.VoiceID,
	}
	if e.options.BaseURL != "" {
		params["base_url"] = e.options.BaseURL
	}
	if e.options.SampleRate != nil {
		params["sample_rate"] = int(*e.options.SampleRate)
	}
	if e.options.OptimizeStreamingLatency != nil {
		params["optimize_streaming_latency"] = *e.options.OptimizeStreamingLatency
	}
	if e.options.Stability != nil {
		params["stability"] = *e.options.Stability
	}
	if e.options.SimilarityBoost != nil {
		params["similarity_boost"] = *e.options.SimilarityBoost
	}
	if e.options.Style != nil {
		params["style"] = *e.options.Style
	}
	if e.options.UseSpeakerBoost != nil {
		params["use_speaker_boost"] = *e.options.UseSpeakerBoost
	}

	config := map[string]interface{}{
		"vendor": "elevenlabs",
		"params": params,
	}
	if e.options.SkipPatterns != nil {
		config["skip_patterns"] = e.options.SkipPatterns
	}
	return config
}

type MicrosoftTTSOptions struct {
	Key          string
	Region       string
	VoiceName    string
	SampleRate   *SampleRate
	SkipPatterns []int
}

type MicrosoftTTS struct {
	options MicrosoftTTSOptions
}

func NewMicrosoftTTS(opts MicrosoftTTSOptions) *MicrosoftTTS {
	if opts.Key == "" {
		panic("MicrosoftTTS requires Key")
	}
	if opts.Region == "" {
		panic("MicrosoftTTS requires Region")
	}
	if opts.VoiceName == "" {
		panic("MicrosoftTTS requires VoiceName")
	}
	return &MicrosoftTTS{options: opts}
}

func (m *MicrosoftTTS) GetSampleRate() *SampleRate {
	return m.options.SampleRate
}

func (m *MicrosoftTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":        m.options.Key,
		"region":     m.options.Region,
		"voice_name": m.options.VoiceName,
	}
	if m.options.SampleRate != nil {
		params["sample_rate"] = int(*m.options.SampleRate)
	}

	config := map[string]interface{}{
		"vendor": "microsoft",
		"params": params,
	}
	if m.options.SkipPatterns != nil {
		config["skip_patterns"] = m.options.SkipPatterns
	}
	return config
}

type OpenAITTSOptions struct {
	APIKey         string
	Voice          string
	Model          string
	ResponseFormat string
	Speed          *float64
	SkipPatterns   []int
}

type OpenAITTS struct {
	options OpenAITTSOptions
}

func NewOpenAITTS(opts OpenAITTSOptions) *OpenAITTS {
	if opts.Voice == "" {
		panic("OpenAITTS requires Voice")
	}
	if opts.APIKey == "" {
		model := strings.ToLower(strings.TrimSpace(opts.Model))
		if model != "" && model != "tts-1" {
			panic("OpenAITTS requires APIKey unless using the preset-backed tts-1 model")
		}
	}
	return &OpenAITTS{options: opts}
}

func (o *OpenAITTS) GetSampleRate() *SampleRate {
	sr := SampleRate24kHz
	return &sr
}

func (o *OpenAITTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"voice": o.options.Voice,
	}
	if o.options.APIKey != "" {
		params["api_key"] = o.options.APIKey
	}
	if o.options.Model != "" {
		params["model"] = o.options.Model
	}
	if o.options.ResponseFormat != "" {
		params["response_format"] = o.options.ResponseFormat
	}
	if o.options.Speed != nil {
		params["speed"] = *o.options.Speed
	}

	config := map[string]interface{}{
		"vendor": "openai",
		"params": params,
	}
	if o.options.SkipPatterns != nil {
		config["skip_patterns"] = o.options.SkipPatterns
	}
	return config
}

type CartesiaTTSOptions struct {
	APIKey       string
	VoiceID      string
	ModelID      string
	SampleRate   *SampleRate
	SkipPatterns []int
}

type CartesiaTTS struct {
	options CartesiaTTSOptions
}

func NewCartesiaTTS(opts CartesiaTTSOptions) *CartesiaTTS {
	if opts.APIKey == "" {
		panic("CartesiaTTS requires APIKey")
	}
	if opts.VoiceID == "" {
		panic("CartesiaTTS requires VoiceID")
	}
	return &CartesiaTTS{options: opts}
}

func (c *CartesiaTTS) GetSampleRate() *SampleRate {
	return c.options.SampleRate
}

func (c *CartesiaTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": c.options.APIKey,
		"voice":   map[string]interface{}{"mode": "id", "id": c.options.VoiceID},
	}
	if c.options.ModelID != "" {
		params["model_id"] = c.options.ModelID
	}
	if c.options.SampleRate != nil {
		params["sample_rate"] = int(*c.options.SampleRate)
	}

	config := map[string]interface{}{
		"vendor": "cartesia",
		"params": params,
	}
	if c.options.SkipPatterns != nil {
		config["skip_patterns"] = c.options.SkipPatterns
	}
	return config
}

type GoogleTTSOptions struct {
	Key          string
	VoiceName    string
	LanguageCode string
	SampleRate   *SampleRate
	SkipPatterns []int
}

type GoogleTTS struct {
	options GoogleTTSOptions
}

func NewGoogleTTS(opts GoogleTTSOptions) *GoogleTTS {
	if opts.Key == "" {
		panic("GoogleTTS requires Key")
	}
	if opts.VoiceName == "" {
		panic("GoogleTTS requires VoiceName")
	}
	return &GoogleTTS{options: opts}
}

func (g *GoogleTTS) GetSampleRate() *SampleRate {
	return g.options.SampleRate
}

func (g *GoogleTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":        g.options.Key,
		"voice_name": g.options.VoiceName,
	}
	if g.options.LanguageCode != "" {
		params["language_code"] = g.options.LanguageCode
	}
	if g.options.SampleRate != nil {
		params["sample_rate_hertz"] = int(*g.options.SampleRate)
	}

	config := map[string]interface{}{
		"vendor": "google",
		"params": params,
	}
	if g.options.SkipPatterns != nil {
		config["skip_patterns"] = g.options.SkipPatterns
	}
	return config
}

type AmazonTTSOptions struct {
	AccessKey    string
	SecretKey    string
	Region       string
	VoiceID      string
	SkipPatterns []int
}

type AmazonTTS struct {
	options AmazonTTSOptions
}

func NewAmazonTTS(opts AmazonTTSOptions) *AmazonTTS {
	if opts.AccessKey == "" {
		panic("AmazonTTS requires AccessKey")
	}
	if opts.SecretKey == "" {
		panic("AmazonTTS requires SecretKey")
	}
	if opts.Region == "" {
		panic("AmazonTTS requires Region")
	}
	if opts.VoiceID == "" {
		panic("AmazonTTS requires VoiceID")
	}
	return &AmazonTTS{options: opts}
}

func (a *AmazonTTS) GetSampleRate() *SampleRate {
	return nil
}

func (a *AmazonTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"access_key": a.options.AccessKey,
		"secret_key": a.options.SecretKey,
		"region":     a.options.Region,
		"voice_id":   a.options.VoiceID,
	}

	config := map[string]interface{}{
		"vendor": "amazon",
		"params": params,
	}
	if a.options.SkipPatterns != nil {
		config["skip_patterns"] = a.options.SkipPatterns
	}
	return config
}

type HumeAITTSOptions struct {
	Key          string
	ConfigID     string
	SkipPatterns []int
}

type HumeAITTS struct {
	options HumeAITTSOptions
}

func NewHumeAITTS(opts HumeAITTSOptions) *HumeAITTS {
	if opts.Key == "" {
		panic("HumeAITTS requires Key")
	}
	return &HumeAITTS{options: opts}
}

func (h *HumeAITTS) GetSampleRate() *SampleRate {
	return nil
}

func (h *HumeAITTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key": h.options.Key,
	}
	if h.options.ConfigID != "" {
		params["config_id"] = h.options.ConfigID
	}

	config := map[string]interface{}{
		"vendor": "humeai",
		"params": params,
	}
	if h.options.SkipPatterns != nil {
		config["skip_patterns"] = h.options.SkipPatterns
	}
	return config
}

type RimeTTSOptions struct {
	Key          string
	Speaker      string
	ModelID      string
	Lang         string
	SamplingRate *int
	SpeedAlpha   *float64
	SkipPatterns []int
}

type RimeTTS struct {
	options RimeTTSOptions
}

func NewRimeTTS(opts RimeTTSOptions) *RimeTTS {
	if opts.Key == "" {
		panic("RimeTTS requires Key")
	}
	if opts.Speaker == "" {
		panic("RimeTTS requires Speaker")
	}
	return &RimeTTS{options: opts}
}

func (r *RimeTTS) GetSampleRate() *SampleRate {
	return nil
}

func (r *RimeTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":     r.options.Key,
		"speaker": r.options.Speaker,
	}
	if r.options.ModelID != "" {
		params["model_id"] = r.options.ModelID
	}
	if r.options.Lang != "" {
		params["lang"] = r.options.Lang
	}
	if r.options.SamplingRate != nil {
		params["samplingRate"] = *r.options.SamplingRate
	}
	if r.options.SpeedAlpha != nil {
		params["speedAlpha"] = *r.options.SpeedAlpha
	}

	config := map[string]interface{}{
		"vendor": "rime",
		"params": params,
	}
	if r.options.SkipPatterns != nil {
		config["skip_patterns"] = r.options.SkipPatterns
	}
	return config
}

type FishAudioTTSOptions struct {
	Key          string
	ReferenceID  string
	SkipPatterns []int
}

type FishAudioTTS struct {
	options FishAudioTTSOptions
}

func NewFishAudioTTS(opts FishAudioTTSOptions) *FishAudioTTS {
	if opts.Key == "" {
		panic("FishAudioTTS requires Key")
	}
	if opts.ReferenceID == "" {
		panic("FishAudioTTS requires ReferenceID")
	}
	return &FishAudioTTS{options: opts}
}

func (f *FishAudioTTS) GetSampleRate() *SampleRate {
	return nil
}

func (f *FishAudioTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":          f.options.Key,
		"reference_id": f.options.ReferenceID,
	}

	config := map[string]interface{}{
		"vendor": "fishaudio",
		"params": params,
	}
	if f.options.SkipPatterns != nil {
		config["skip_patterns"] = f.options.SkipPatterns
	}
	return config
}

type MiniMaxTTSOptions struct {
	Key          string
	GroupID      string
	Model        string
	VoiceID      string
	URL          string
	SkipPatterns []int
}

type MiniMaxTTS struct {
	options MiniMaxTTSOptions
}

func NewMiniMaxTTS(opts MiniMaxTTSOptions) *MiniMaxTTS {
	if opts.Key == "" {
		model := strings.ToLower(strings.TrimSpace(opts.Model))
		switch model {
		case "speech-2.6-turbo", "speech_2_6_turbo", "speech-2.8-turbo", "speech_2_8_turbo":
			if opts.GroupID != "" || opts.VoiceID != "" || opts.URL != "" {
				panic("MiniMaxTTS preset-backed mode does not allow GroupID, VoiceID, or URL")
			}
		default:
			panic("MiniMaxTTS requires Key unless using a supported preset-backed model")
		}
	}
	if opts.Key != "" && opts.GroupID == "" {
		panic("MiniMaxTTS requires GroupID")
	}
	return &MiniMaxTTS{options: opts}
}

func (m *MiniMaxTTS) GetSampleRate() *SampleRate {
	return nil
}

func (m *MiniMaxTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"model": m.options.Model,
	}
	if m.options.Key != "" {
		params["key"] = m.options.Key
	}
	if m.options.GroupID != "" {
		params["group_id"] = m.options.GroupID
	}
	if m.options.VoiceID != "" {
		params["voice_setting"] = map[string]interface{}{"voice_id": m.options.VoiceID}
	}
	if m.options.URL != "" {
		params["url"] = m.options.URL
	}

	config := map[string]interface{}{
		"vendor": "minimax",
		"params": params,
	}
	if m.options.SkipPatterns != nil {
		config["skip_patterns"] = m.options.SkipPatterns
	}
	return config
}

type SarvamTTSOptions struct {
	Key                string
	Speaker            string
	TargetLanguageCode string
	SkipPatterns       []int
}

type SarvamTTS struct {
	options SarvamTTSOptions
}

func NewSarvamTTS(opts SarvamTTSOptions) *SarvamTTS {
	if opts.Key == "" {
		panic("SarvamTTS requires Key")
	}
	if opts.Speaker == "" {
		panic("SarvamTTS requires Speaker")
	}
	if opts.TargetLanguageCode == "" {
		panic("SarvamTTS requires TargetLanguageCode")
	}
	return &SarvamTTS{options: opts}
}

func (s *SarvamTTS) GetSampleRate() *SampleRate {
	return nil
}

func (s *SarvamTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":                  s.options.Key,
		"speaker":              s.options.Speaker,
		"target_language_code": s.options.TargetLanguageCode,
	}

	config := map[string]interface{}{
		"vendor": "sarvam",
		"params": params,
	}
	if s.options.SkipPatterns != nil {
		config["skip_patterns"] = s.options.SkipPatterns
	}
	return config
}

type MurfTTSOptions struct {
	Key          string
	VoiceID      string
	Style        string
	SkipPatterns []int
}

type MurfTTS struct {
	options MurfTTSOptions
}

func NewMurfTTS(opts MurfTTSOptions) *MurfTTS {
	if opts.Key == "" {
		panic("MurfTTS requires Key")
	}
	if opts.VoiceID == "" {
		panic("MurfTTS requires VoiceID")
	}
	return &MurfTTS{options: opts}
}

func (m *MurfTTS) GetSampleRate() *SampleRate {
	return nil
}

func (m *MurfTTS) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":      m.options.Key,
		"voice_id": m.options.VoiceID,
	}
	if m.options.Style != "" {
		params["style"] = m.options.Style
	}

	config := map[string]interface{}{
		"vendor": "murf",
		"params": params,
	}
	if m.options.SkipPatterns != nil {
		config["skip_patterns"] = m.options.SkipPatterns
	}
	return config
}
