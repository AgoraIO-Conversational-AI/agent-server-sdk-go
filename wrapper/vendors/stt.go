package vendors

type SpeechmaticsSTTOptions struct {
	APIKey   string
	Language string
	Model    string
}

type SpeechmaticsSTT struct {
	options SpeechmaticsSTTOptions
}

func NewSpeechmaticsSTT(opts SpeechmaticsSTTOptions) *SpeechmaticsSTT {
	if opts.APIKey == "" {
		panic("SpeechmaticsSTT requires APIKey")
	}
	return &SpeechmaticsSTT{options: opts}
}

func (s *SpeechmaticsSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": s.options.APIKey,
	}
	if s.options.Language != "" {
		params["language"] = s.options.Language
	}
	if s.options.Model != "" {
		params["model"] = s.options.Model
	}

	return map[string]interface{}{
		"vendor": "speechmatics",
		"params": params,
	}
}

type DeepgramSTTOptions struct {
	APIKey   string
	Model    string
	Language string
}

type DeepgramSTT struct {
	options DeepgramSTTOptions
}

func NewDeepgramSTT(opts DeepgramSTTOptions) *DeepgramSTT {
	if opts.APIKey == "" {
		panic("DeepgramSTT requires APIKey")
	}
	return &DeepgramSTT{options: opts}
}

func (d *DeepgramSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": d.options.APIKey,
	}
	if d.options.Model != "" {
		params["model"] = d.options.Model
	}
	if d.options.Language != "" {
		params["language"] = d.options.Language
	}

	return map[string]interface{}{
		"vendor": "deepgram",
		"params": params,
	}
}

type MicrosoftSTTOptions struct {
	Key      string
	Region   string
	Language string
}

type MicrosoftSTT struct {
	options MicrosoftSTTOptions
}

func NewMicrosoftSTT(opts MicrosoftSTTOptions) *MicrosoftSTT {
	if opts.Key == "" {
		panic("MicrosoftSTT requires Key")
	}
	if opts.Region == "" {
		panic("MicrosoftSTT requires Region")
	}
	return &MicrosoftSTT{options: opts}
}

func (m *MicrosoftSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key":    m.options.Key,
		"region": m.options.Region,
	}
	if m.options.Language != "" {
		params["language"] = m.options.Language
	}

	return map[string]interface{}{
		"vendor": "microsoft",
		"params": params,
	}
}

type OpenAISTTOptions struct {
	APIKey   string
	Model    string
	Language string
}

type OpenAISTT struct {
	options OpenAISTTOptions
}

func NewOpenAISTT(opts OpenAISTTOptions) *OpenAISTT {
	if opts.APIKey == "" {
		panic("OpenAISTT requires APIKey")
	}
	return &OpenAISTT{options: opts}
}

func (o *OpenAISTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": o.options.APIKey,
	}
	if o.options.Model != "" {
		params["model"] = o.options.Model
	}
	if o.options.Language != "" {
		params["language"] = o.options.Language
	}

	return map[string]interface{}{
		"vendor": "openai",
		"params": params,
	}
}

type GoogleSTTOptions struct {
	Key      string
	Language string
	Model    string
}

type GoogleSTT struct {
	options GoogleSTTOptions
}

func NewGoogleSTT(opts GoogleSTTOptions) *GoogleSTT {
	if opts.Key == "" {
		panic("GoogleSTT requires Key")
	}
	return &GoogleSTT{options: opts}
}

func (g *GoogleSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"key": g.options.Key,
	}
	if g.options.Language != "" {
		params["language"] = g.options.Language
	}
	if g.options.Model != "" {
		params["model"] = g.options.Model
	}

	return map[string]interface{}{
		"vendor": "google",
		"params": params,
	}
}

type AmazonSTTOptions struct {
	AccessKey string
	SecretKey string
	Region    string
	Language  string
}

type AmazonSTT struct {
	options AmazonSTTOptions
}

func NewAmazonSTT(opts AmazonSTTOptions) *AmazonSTT {
	if opts.AccessKey == "" {
		panic("AmazonSTT requires AccessKey")
	}
	if opts.SecretKey == "" {
		panic("AmazonSTT requires SecretKey")
	}
	if opts.Region == "" {
		panic("AmazonSTT requires Region")
	}
	return &AmazonSTT{options: opts}
}

func (a *AmazonSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"access_key": a.options.AccessKey,
		"secret_key": a.options.SecretKey,
		"region":     a.options.Region,
	}
	if a.options.Language != "" {
		params["language"] = a.options.Language
	}

	return map[string]interface{}{
		"vendor": "amazon",
		"params": params,
	}
}

type AssemblyAISTTOptions struct {
	APIKey string
}

type AssemblyAISTT struct {
	options AssemblyAISTTOptions
}

func NewAssemblyAISTT(opts AssemblyAISTTOptions) *AssemblyAISTT {
	if opts.APIKey == "" {
		panic("AssemblyAISTT requires APIKey")
	}
	return &AssemblyAISTT{options: opts}
}

func (a *AssemblyAISTT) ToConfig() map[string]interface{} {
	return map[string]interface{}{
		"vendor": "assemblyai",
		"params": map[string]interface{}{
			"api_key": a.options.APIKey,
		},
	}
}

type AresSTTOptions struct {
	APIKey string
}

type AresSTT struct {
	options AresSTTOptions
}

func NewAresSTT(opts AresSTTOptions) *AresSTT {
	if opts.APIKey == "" {
		panic("AresSTT requires APIKey")
	}
	return &AresSTT{options: opts}
}

func (a *AresSTT) ToConfig() map[string]interface{} {
	return map[string]interface{}{
		"vendor": "ares",
		"params": map[string]interface{}{
			"api_key": a.options.APIKey,
		},
	}
}

type SonioxSTTOptions struct {
	APIKey string
}

type SonioxSTT struct {
	options SonioxSTTOptions
}

func NewSonioxSTT(opts SonioxSTTOptions) *SonioxSTT {
	if opts.APIKey == "" {
		panic("SonioxSTT requires APIKey")
	}
	return &SonioxSTT{options: opts}
}

func (s *SonioxSTT) ToConfig() map[string]interface{} {
	return map[string]interface{}{
		"vendor": "soniox",
		"params": map[string]interface{}{
			"api_key": s.options.APIKey,
		},
	}
}

type SarvamSTTOptions struct {
	APIKey   string
	Language string
	Model    string
}

type SarvamSTT struct {
	options SarvamSTTOptions
}

func NewSarvamSTT(opts SarvamSTTOptions) *SarvamSTT {
	if opts.APIKey == "" {
		panic("SarvamSTT requires APIKey")
	}
	return &SarvamSTT{options: opts}
}

func (s *SarvamSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key": s.options.APIKey,
	}
	if s.options.Language != "" {
		params["language"] = s.options.Language
	}
	if s.options.Model != "" {
		params["model"] = s.options.Model
	}

	return map[string]interface{}{
		"vendor": "sarvam",
		"params": params,
	}
}
