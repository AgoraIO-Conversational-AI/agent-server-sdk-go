package vendors

import "fmt"

type OpenAIOptions struct {
	APIKey          string
	Model           string
	BaseURL         string
	Temperature     *float64
	TopP            *float64
	MaxTokens       *int
	SystemMessages  []map[string]interface{}
	GreetingMessage string
	FailureMessage  string
	InputModalities []string
	Params          map[string]interface{}
}

type OpenAI struct {
	options OpenAIOptions
}

func NewOpenAI(opts OpenAIOptions) *OpenAI {
	if opts.APIKey == "" {
		panic("OpenAI requires APIKey")
	}
	if opts.Model == "" {
		opts.Model = "gpt-4o-mini"
	}
	return &OpenAI{options: opts}
}

func (o *OpenAI) ToConfig() map[string]interface{} {
	url := o.options.BaseURL
	if url == "" {
		url = "https://api.openai.com/v1/chat/completions"
	}

	params := o.options.Params
	if params == nil {
		params = map[string]interface{}{"model": o.options.Model}
	}
	if o.options.Temperature != nil {
		params["temperature"] = *o.options.Temperature
	}
	if o.options.TopP != nil {
		params["top_p"] = *o.options.TopP
	}
	if o.options.MaxTokens != nil {
		params["max_tokens"] = *o.options.MaxTokens
	}

	inputMod := o.options.InputModalities
	if inputMod == nil {
		inputMod = []string{"text"}
	}

	config := map[string]interface{}{
		"url":              url,
		"api_key":          o.options.APIKey,
		"params":           params,
		"style":            "openai",
		"input_modalities": inputMod,
	}

	if o.options.SystemMessages != nil {
		config["system_messages"] = o.options.SystemMessages
	}
	if o.options.GreetingMessage != "" {
		config["greeting_message"] = o.options.GreetingMessage
	}
	if o.options.FailureMessage != "" {
		config["failure_message"] = o.options.FailureMessage
	}

	return config
}

type AzureOpenAIOptions struct {
	APIKey          string
	Endpoint        string
	DeploymentName  string
	APIVersion      string
	Temperature     *float64
	TopP            *float64
	MaxTokens       *int
	SystemMessages  []map[string]interface{}
	GreetingMessage string
	FailureMessage  string
	InputModalities []string
}

type AzureOpenAI struct {
	options AzureOpenAIOptions
}

func NewAzureOpenAI(opts AzureOpenAIOptions) *AzureOpenAI {
	if opts.APIKey == "" {
		panic("AzureOpenAI requires APIKey")
	}
	if opts.Endpoint == "" {
		panic("AzureOpenAI requires Endpoint")
	}
	if opts.DeploymentName == "" {
		panic("AzureOpenAI requires DeploymentName")
	}
	if opts.APIVersion == "" {
		opts.APIVersion = "2024-08-01-preview"
	}
	return &AzureOpenAI{options: opts}
}

func (a *AzureOpenAI) ToConfig() map[string]interface{} {
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		a.options.Endpoint, a.options.DeploymentName, a.options.APIVersion)

	inputMod := a.options.InputModalities
	if inputMod == nil {
		inputMod = []string{"text"}
	}

	config := map[string]interface{}{
		"url":              url,
		"api_key":          a.options.APIKey,
		"vendor":           "azure",
		"style":            "openai",
		"input_modalities": inputMod,
	}

	params := map[string]interface{}{}
	if a.options.Temperature != nil {
		params["temperature"] = *a.options.Temperature
	}
	if a.options.TopP != nil {
		params["top_p"] = *a.options.TopP
	}
	if a.options.MaxTokens != nil {
		params["max_tokens"] = *a.options.MaxTokens
	}
	if len(params) > 0 {
		config["params"] = params
	}

	if a.options.SystemMessages != nil {
		config["system_messages"] = a.options.SystemMessages
	}
	if a.options.GreetingMessage != "" {
		config["greeting_message"] = a.options.GreetingMessage
	}
	if a.options.FailureMessage != "" {
		config["failure_message"] = a.options.FailureMessage
	}

	return config
}

type AnthropicOptions struct {
	APIKey          string
	Model           string
	MaxTokens       *int
	Temperature     *float64
	TopP            *float64
	SystemMessages  []map[string]interface{}
	GreetingMessage string
	FailureMessage  string
	InputModalities []string
}

type Anthropic struct {
	options AnthropicOptions
}

func NewAnthropic(opts AnthropicOptions) *Anthropic {
	if opts.APIKey == "" {
		panic("Anthropic requires APIKey")
	}
	if opts.Model == "" {
		opts.Model = "claude-3-5-sonnet-20241022"
	}
	return &Anthropic{options: opts}
}

func (a *Anthropic) ToConfig() map[string]interface{} {
	inputMod := a.options.InputModalities
	if inputMod == nil {
		inputMod = []string{"text"}
	}

	params := map[string]interface{}{"model": a.options.Model}
	if a.options.MaxTokens != nil {
		params["max_tokens"] = *a.options.MaxTokens
	}
	if a.options.Temperature != nil {
		params["temperature"] = *a.options.Temperature
	}
	if a.options.TopP != nil {
		params["top_p"] = *a.options.TopP
	}

	config := map[string]interface{}{
		"url":              "https://api.anthropic.com/v1/messages",
		"api_key":          a.options.APIKey,
		"params":           params,
		"style":            "anthropic",
		"input_modalities": inputMod,
	}

	if a.options.SystemMessages != nil {
		config["system_messages"] = a.options.SystemMessages
	}
	if a.options.GreetingMessage != "" {
		config["greeting_message"] = a.options.GreetingMessage
	}
	if a.options.FailureMessage != "" {
		config["failure_message"] = a.options.FailureMessage
	}

	return config
}

type GeminiOptions struct {
	APIKey          string
	Model           string
	Temperature     *float64
	TopP            *float64
	TopK            *int
	MaxOutputTokens *int
	SystemMessages  []map[string]interface{}
	GreetingMessage string
	FailureMessage  string
	InputModalities []string
}

type Gemini struct {
	options GeminiOptions
}

func NewGemini(opts GeminiOptions) *Gemini {
	if opts.APIKey == "" {
		panic("Gemini requires APIKey")
	}
	if opts.Model == "" {
		opts.Model = "gemini-2.0-flash-exp"
	}
	return &Gemini{options: opts}
}

func (g *Gemini) ToConfig() map[string]interface{} {
	inputMod := g.options.InputModalities
	if inputMod == nil {
		inputMod = []string{"text"}
	}

	params := map[string]interface{}{"model": g.options.Model}
	if g.options.Temperature != nil {
		params["temperature"] = *g.options.Temperature
	}
	if g.options.TopP != nil {
		params["top_p"] = *g.options.TopP
	}
	if g.options.TopK != nil {
		params["top_k"] = *g.options.TopK
	}
	if g.options.MaxOutputTokens != nil {
		params["max_output_tokens"] = *g.options.MaxOutputTokens
	}

	config := map[string]interface{}{
		"url":              "https://generativelanguage.googleapis.com/v1beta/models",
		"api_key":          g.options.APIKey,
		"params":           params,
		"style":            "gemini",
		"input_modalities": inputMod,
	}

	if g.options.SystemMessages != nil {
		config["system_messages"] = g.options.SystemMessages
	}
	if g.options.GreetingMessage != "" {
		config["greeting_message"] = g.options.GreetingMessage
	}
	if g.options.FailureMessage != "" {
		config["failure_message"] = g.options.FailureMessage
	}

	return config
}
