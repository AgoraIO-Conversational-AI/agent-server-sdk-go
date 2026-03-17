package vendors

type SpeechmaticsSTTOptions struct {
	APIKey           string
	Language         string
	Model            string
	AdditionalParams map[string]interface{}
}

type SpeechmaticsSTT struct {
	options SpeechmaticsSTTOptions
}

func NewSpeechmaticsSTT(opts SpeechmaticsSTTOptions) *SpeechmaticsSTT {
	if opts.APIKey == "" {
		panic("SpeechmaticsSTT requires APIKey")
	}
	if opts.Language == "" {
		panic("SpeechmaticsSTT requires Language")
	}
	return &SpeechmaticsSTT{options: opts}
}

func (s *SpeechmaticsSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{
		"api_key":  s.options.APIKey,
		"language": s.options.Language,
	}
	if s.options.Model != "" {
		params["model"] = s.options.Model
	}
	for k, v := range s.options.AdditionalParams {
		if _, exists := params[k]; !exists {
			params[k] = v
		}
	}

	return map[string]interface{}{
		"vendor":   "speechmatics",
		"language": s.options.Language,
		"params":   params,
	}
}

type DeepgramSTTOptions struct {
	APIKey           string
	Model            string
	Language         string
	SmartFormat      *bool
	Punctuation      *bool
	AdditionalParams map[string]interface{}
}

type DeepgramSTT struct {
	options DeepgramSTTOptions
}

func NewDeepgramSTT(opts DeepgramSTTOptions) *DeepgramSTT {
	return &DeepgramSTT{options: opts}
}

func (d *DeepgramSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range d.options.AdditionalParams {
		params[k] = v
	}
	if d.options.APIKey != "" {
		params["api_key"] = d.options.APIKey
	}
	if d.options.Model != "" {
		params["model"] = d.options.Model
	}
	if d.options.Language != "" {
		params["language"] = d.options.Language
	}
	if d.options.SmartFormat != nil {
		params["smart_format"] = *d.options.SmartFormat
	}
	if d.options.Punctuation != nil {
		params["punctuation"] = *d.options.Punctuation
	}

	config := map[string]interface{}{
		"vendor": "deepgram",
		"params": params,
	}
	if d.options.Language != "" {
		config["language"] = d.options.Language
	}
	return config
}

type MicrosoftSTTOptions struct {
	Key              string
	Region           string
	Language         string
	AdditionalParams map[string]interface{}
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
	params := map[string]interface{}{}
	for k, v := range m.options.AdditionalParams {
		params[k] = v
	}
	params["key"] = m.options.Key
	params["region"] = m.options.Region
	if m.options.Language != "" {
		params["language"] = m.options.Language
	}

	config := map[string]interface{}{
		"vendor": "microsoft",
		"params": params,
	}
	if m.options.Language != "" {
		config["language"] = m.options.Language
	}
	return config
}

type OpenAISTTOptions struct {
	APIKey           string
	Model            string
	Language         string
	AdditionalParams map[string]interface{}
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
	params := map[string]interface{}{}
	for k, v := range o.options.AdditionalParams {
		params[k] = v
	}
	params["api_key"] = o.options.APIKey
	if o.options.Model != "" {
		params["model"] = o.options.Model
	}
	if o.options.Language != "" {
		params["language"] = o.options.Language
	}

	config := map[string]interface{}{
		"vendor": "openai",
		"params": params,
	}
	if o.options.Language != "" {
		config["language"] = o.options.Language
	}
	return config
}

type GoogleSTTOptions struct {
	APIKey           string
	Language         string
	Model            string
	AdditionalParams map[string]interface{}
}

type GoogleSTT struct {
	options GoogleSTTOptions
}

func NewGoogleSTT(opts GoogleSTTOptions) *GoogleSTT {
	if opts.APIKey == "" {
		panic("GoogleSTT requires APIKey")
	}
	return &GoogleSTT{options: opts}
}

func (g *GoogleSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range g.options.AdditionalParams {
		params[k] = v
	}
	params["api_key"] = g.options.APIKey
	if g.options.Language != "" {
		params["language"] = g.options.Language
	}
	if g.options.Model != "" {
		params["model"] = g.options.Model
	}

	config := map[string]interface{}{
		"vendor": "google",
		"params": params,
	}
	if g.options.Language != "" {
		config["language"] = g.options.Language
	}
	return config
}

type AmazonSTTOptions struct {
	AccessKey        string
	SecretKey        string
	Region           string
	Language         string
	AdditionalParams map[string]interface{}
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
	params := map[string]interface{}{}
	for k, v := range a.options.AdditionalParams {
		params[k] = v
	}
	params["access_key"] = a.options.AccessKey
	params["secret_key"] = a.options.SecretKey
	params["region"] = a.options.Region
	if a.options.Language != "" {
		params["language"] = a.options.Language
	}

	config := map[string]interface{}{
		"vendor": "amazon",
		"params": params,
	}
	if a.options.Language != "" {
		config["language"] = a.options.Language
	}
	return config
}

type AssemblyAISTTOptions struct {
	APIKey           string
	Language         string
	AdditionalParams map[string]interface{}
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
	params := map[string]interface{}{}
	for k, v := range a.options.AdditionalParams {
		params[k] = v
	}
	params["api_key"] = a.options.APIKey

	config := map[string]interface{}{
		"vendor": "assemblyai",
		"params": params,
	}
	if a.options.Language != "" {
		config["language"] = a.options.Language
	}
	return config
}

type AresSTTOptions struct {
	Language         string
	AdditionalParams map[string]interface{}
}

type AresSTT struct {
	options AresSTTOptions
}

func NewAresSTT(opts AresSTTOptions) *AresSTT {
	return &AresSTT{options: opts}
}

func (a *AresSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range a.options.AdditionalParams {
		params[k] = v
	}
	if a.options.Language != "" {
		params["language"] = a.options.Language
	}

	config := map[string]interface{}{
		"vendor": "ares",
	}
	if len(params) > 0 {
		config["params"] = params
	}
	if a.options.Language != "" {
		config["language"] = a.options.Language
	}
	return config
}

type SarvamSTTOptions struct {
	APIKey           string
	Language         string
	Model            string
	AdditionalParams map[string]interface{}
}

type SarvamSTT struct {
	options SarvamSTTOptions
}

func NewSarvamSTT(opts SarvamSTTOptions) *SarvamSTT {
	if opts.APIKey == "" {
		panic("SarvamSTT requires APIKey")
	}
	if opts.Language == "" {
		panic("SarvamSTT requires Language")
	}
	return &SarvamSTT{options: opts}
}

func (s *SarvamSTT) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range s.options.AdditionalParams {
		params[k] = v
	}
	params["api_key"] = s.options.APIKey
	params["language"] = s.options.Language
	if s.options.Model != "" {
		params["model"] = s.options.Model
	}

	return map[string]interface{}{
		"vendor":   "sarvam",
		"language": s.options.Language,
		"params":   params,
	}
}
