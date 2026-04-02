package vendors

type OpenAIRealtimeOptions struct {
	APIKey           string
	Model            string
	URL              string
	GreetingMessage  string
	InputModalities  []string
	OutputModalities []string
	Messages         []map[string]interface{}
	Params           map[string]interface{}
}

type OpenAIRealtime struct {
	options OpenAIRealtimeOptions
}

func NewOpenAIRealtime(opts OpenAIRealtimeOptions) *OpenAIRealtime {
	if opts.APIKey == "" {
		panic("OpenAIRealtime requires APIKey")
	}
	if opts.Model == "" {
		opts.Model = "gpt-4o-realtime-preview"
	}
	return &OpenAIRealtime{options: opts}
}

func (o *OpenAIRealtime) ToConfig() map[string]interface{} {
	var params map[string]interface{}
	if o.options.Model != "" || o.options.Params != nil {
		params = map[string]interface{}{}
		for k, v := range o.options.Params {
			params[k] = v
		}
		if o.options.Model != "" {
			params["model"] = o.options.Model
		}
	}

	config := map[string]interface{}{
		"vendor":  "openai",
		"style":   "openai",
		"api_key": o.options.APIKey,
	}
	if o.options.URL != "" {
		config["url"] = o.options.URL
	}
	if params != nil {
		config["params"] = params
	}

	if o.options.GreetingMessage != "" {
		config["greeting_message"] = o.options.GreetingMessage
	}
	if o.options.InputModalities != nil {
		config["input_modalities"] = o.options.InputModalities
	}
	if o.options.OutputModalities != nil {
		config["output_modalities"] = o.options.OutputModalities
	}
	if o.options.Messages != nil {
		config["messages"] = o.options.Messages
	}

	return config
}

type GeminiLiveOptions struct {
	APIKey           string
	Model            string
	Instructions     string
	Voice            string
	GreetingMessage  string
	InputModalities  []string
	OutputModalities []string
	Messages         []map[string]interface{}
	AdditionalParams map[string]interface{}
}

type GeminiLive struct {
	options GeminiLiveOptions
}

func NewGeminiLive(opts GeminiLiveOptions) *GeminiLive {
	if opts.APIKey == "" {
		panic("GeminiLive requires APIKey")
	}
	if opts.Model == "" {
		panic("GeminiLive requires Model")
	}
	return &GeminiLive{options: opts}
}

func (g *GeminiLive) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range g.options.AdditionalParams {
		params[k] = v
	}
	params["model"] = g.options.Model
	if g.options.Instructions != "" {
		params["instructions"] = g.options.Instructions
	}
	if g.options.Voice != "" {
		params["voice"] = g.options.Voice
	}

	config := map[string]interface{}{
		"vendor":  "gemini",
		"style":   "openai",
		"api_key": g.options.APIKey,
		"params":  params,
	}
	if g.options.GreetingMessage != "" {
		config["greeting_message"] = g.options.GreetingMessage
	}
	if g.options.InputModalities != nil {
		config["input_modalities"] = g.options.InputModalities
	}
	if g.options.OutputModalities != nil {
		config["output_modalities"] = g.options.OutputModalities
	}
	if g.options.Messages != nil {
		config["messages"] = g.options.Messages
	}
	return config
}

type VertexAIOptions struct {
	ProjectID           string
	Location            string
	Model               string
	Voice               string
	Instructions        string
	Messages            []map[string]interface{}
	ADCredentialsString string
	AdditionalParams    map[string]interface{}
	GreetingMessage     string
	InputModalities     []string
	OutputModalities    []string
}

type VertexAI struct {
	options VertexAIOptions
}

func NewVertexAI(opts VertexAIOptions) *VertexAI {
	if opts.ProjectID == "" {
		panic("VertexAI requires ProjectID")
	}
	if opts.ADCredentialsString == "" {
		panic("VertexAI requires ADCredentialsString")
	}
	if opts.Location == "" {
		opts.Location = "us-central1"
	}
	if opts.Model == "" {
		opts.Model = "gemini-2.0-flash-exp"
	}
	return &VertexAI{options: opts}
}

func (v *VertexAI) ToConfig() map[string]interface{} {
	params := map[string]interface{}{}
	for k, val := range v.options.AdditionalParams {
		params[k] = val
	}
	params["project_id"] = v.options.ProjectID
	params["location"] = v.options.Location
	params["model"] = v.options.Model
	params["adc_credentials_string"] = v.options.ADCredentialsString
	if v.options.Voice != "" {
		params["voice"] = v.options.Voice
	}
	if v.options.Instructions != "" {
		params["instructions"] = v.options.Instructions
	}

	config := map[string]interface{}{
		"vendor": "vertexai",
		"style":  "openai",
		"params": params,
	}

	if v.options.GreetingMessage != "" {
		config["greeting_message"] = v.options.GreetingMessage
	}
	if v.options.InputModalities != nil {
		config["input_modalities"] = v.options.InputModalities
	}
	if v.options.OutputModalities != nil {
		config["output_modalities"] = v.options.OutputModalities
	}
	if v.options.Messages != nil {
		config["messages"] = v.options.Messages
	}

	return config
}
