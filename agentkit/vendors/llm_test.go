package vendors

import "testing"

func TestGroqSerializesAsOpenAICompatible(t *testing.T) {
	config := NewGroq(GroqOptions{
		APIKey: "groq-key",
		Model:  "llama-3.3-70b-versatile",
	}).ToConfig()

	if config["url"] != "https://api.groq.com/openai/v1/chat/completions" {
		t.Fatalf("unexpected url: %v", config["url"])
	}
	if config["style"] != "openai" {
		t.Fatalf("unexpected style: %v", config["style"])
	}
}

func TestCustomLLMMarksRequestAsCustom(t *testing.T) {
	config := NewCustomLLM(CustomLLMOptions{
		APIKey:  "key",
		Model:   "model",
		BaseURL: "https://llm.example.com/chat",
	}).ToConfig()

	if config["vendor"] != "custom" {
		t.Fatalf("unexpected vendor: %v", config["vendor"])
	}
	if config["style"] != "openai" {
		t.Fatalf("unexpected style: %v", config["style"])
	}
}

func TestVertexAILLMIncludesProjectRouting(t *testing.T) {
	config := NewVertexAILLM(VertexAILLMOptions{
		GeminiOptions: GeminiOptions{
			APIKey: "vertex-token",
			Model:  "gemini-2.0-flash",
		},
		ProjectID: "project",
		Location:  "us-central1",
	}).ToConfig()

	params := config["params"].(map[string]interface{})
	if config["style"] != "gemini" {
		t.Fatalf("unexpected style: %v", config["style"])
	}
	if params["project_id"] != "project" || params["location"] != "us-central1" {
		t.Fatalf("unexpected params: %#v", params)
	}
}

func TestAmazonBedrockSerializesAsAnthropicStyle(t *testing.T) {
	config := NewAmazonBedrock(AmazonBedrockOptions{
		APIKey: "bedrock-key",
		URL:    "https://bedrock.example.com/messages",
		Model:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
	}).ToConfig()

	if config["style"] != "anthropic" {
		t.Fatalf("unexpected style: %v", config["style"])
	}
}

func TestDifySerializesConversationFields(t *testing.T) {
	config := NewDify(DifyOptions{
		APIKey:         "dify-key",
		URL:            "https://api.dify.ai/v1/chat-messages",
		User:           "user-1",
		ConversationID: "conversation-1",
	}).ToConfig()

	params := config["params"].(map[string]interface{})
	if config["style"] != "dify" {
		t.Fatalf("unexpected style: %v", config["style"])
	}
	if params["user"] != "user-1" || params["conversation_id"] != "conversation-1" {
		t.Fatalf("unexpected params: %#v", params)
	}
}
