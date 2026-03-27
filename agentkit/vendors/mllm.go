package vendors

type OpenAIRealtimeOptions struct {
	APIKey          string
	Model           string
	Voice           string
	Temperature     *float64
	MaxOutputTokens *int
	SystemMessage   string
	PredefinedTools []string
	FailureMessage  string
	MaxHistory      *int
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
	params := map[string]interface{}{
		"model": o.options.Model,
	}
	if o.options.Voice != "" {
		params["voice"] = o.options.Voice
	}
	if o.options.Temperature != nil {
		params["temperature"] = *o.options.Temperature
	}
	if o.options.MaxOutputTokens != nil {
		params["max_output_tokens"] = *o.options.MaxOutputTokens
	}

	config := map[string]interface{}{
		"vendor":  "openai",
		"api_key": o.options.APIKey,
		"params":  params,
	}

	if o.options.SystemMessage != "" {
		config["system_message"] = o.options.SystemMessage
	}
	if o.options.PredefinedTools != nil {
		config["predefined_tools"] = o.options.PredefinedTools
	}
	if o.options.FailureMessage != "" {
		config["failure_message"] = o.options.FailureMessage
	}
	if o.options.MaxHistory != nil {
		config["max_history"] = *o.options.MaxHistory
	}

	return config
}

type VertexAIOptions struct {
	ProjectID       string
	Location        string
	Model           string
	Voice           string
	Language        string
	SystemMessage   string
	PredefinedTools []string
	FailureMessage  string
	MaxHistory      *int
}

type VertexAI struct {
	options VertexAIOptions
}

func NewVertexAI(opts VertexAIOptions) *VertexAI {
	if opts.ProjectID == "" {
		panic("VertexAI requires ProjectID")
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
	params := map[string]interface{}{
		"project_id": v.options.ProjectID,
		"location":   v.options.Location,
		"model":      v.options.Model,
	}
	if v.options.Voice != "" {
		params["voice"] = v.options.Voice
	}
	if v.options.Language != "" {
		params["language"] = v.options.Language
	}

	config := map[string]interface{}{
		"vendor": "vertexai",
		"params": params,
	}

	if v.options.SystemMessage != "" {
		config["system_message"] = v.options.SystemMessage
	}
	if v.options.PredefinedTools != nil {
		config["predefined_tools"] = v.options.PredefinedTools
	}
	if v.options.FailureMessage != "" {
		config["failure_message"] = v.options.FailureMessage
	}
	if v.options.MaxHistory != nil {
		config["max_history"] = *v.options.MaxHistory
	}

	return config
}
