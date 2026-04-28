package agentkit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agentkit/vendors"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToPropertiesSupportsPresetFlowAndRTMDefault(t *testing.T) {
	enableRTM := true
	agent := NewAgent(
		WithInstructions("Preset flow"),
		WithAdvancedFeatures(&AdvancedFeatures{EnableRtm: &enableRTM}),
	)

	props, err := agent.ToProperties(ToPropertiesOptions{
		Channel:              "room-1",
		Token:                "rtc-token",
		AgentUID:             "1",
		RemoteUIDs:           []string{"100"},
		SkipVendorValidation: true,
	})
	require.NoError(t, err)
	require.NotNil(t, props)
	assert.Equal(t, "room-1", props.Channel)
	assert.NotNil(t, props.Parameters)
	require.NotNil(t, props.Parameters.DataChannel)
	assert.Equal(t, "rtm", string(*props.Parameters.DataChannel))
	assert.Nil(t, props.Llm)
	assert.Nil(t, props.Tts)
}

func TestCreateSessionStartIncludesPresetPipelineAndGetTurns(t *testing.T) {
	var started int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/projects/appid/join":
			var req map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			assert.Equal(t, "support-agent", req["name"])
			assert.Equal(t, "deepgram_nova_3,openai_gpt_4o_mini,openai_tts_1", req["preset"])
			assert.Equal(t, "pipeline_123", req["pipeline_id"])

			props := req["properties"].(map[string]interface{})
			assert.Equal(t, "room-1", props["channel"])
			assert.Equal(t, "1", props["agent_rtc_uid"])

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"agent_id":"agent_123","status":"RUNNING"}`))
			atomic.StoreInt32(&started, 1)
		case "/v2/projects/appid/agents/agent_123/turns":
			assert.Equal(t, int32(1), atomic.LoadInt32(&started))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"turns":[{"agent_id":"agent_123","turn_id":1}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	rawClient := client.NewClient(
		option.WithBaseURL(server.URL),
		option.WithBasicAuth("user", "pass"),
		option.WithMaxAttempts(1),
	)
	agoraClient := &AgoraClient{
		Agents:         rawClient.Agents,
		AppID:          "appid",
		AppCertificate: "app-cert",
		AuthMode:       AuthModeBasic,
	}

	agent := NewAgent(WithName("support-agent"))
	session := agent.CreateSession(agoraClient, CreateSessionOptions{
		Channel:    "room-1",
		AgentUID:   "1",
		RemoteUIDs: []string{"100"},
		Token:      "rtc-token",
		Preset: []string{
			AgentPresets.Asr.DeepgramNova3,
			AgentPresets.Llm.OpenAIGpt4oMini,
			AgentPresets.Tts.OpenAITts1,
		},
		PipelineID: "pipeline_123",
	})

	agentID, err := session.Start(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "agent_123", agentID)

	turns, err := session.GetTurns(context.Background())
	require.NoError(t, err)
	require.Len(t, turns.Turns, 1)
	assert.Equal(t, "agent_123", *turns.Turns[0].AgentID)
}

func TestCreateSessionStartSendsManagedPresetPayloadWithoutGeneratedEmptyFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v2/projects/appid/join", r.URL.Path)
		var req map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "deepgram_nova_3,openai_gpt_4o_mini,minimax_speech_2_6_turbo", req["preset"])

		payload, err := json.Marshal(req)
		require.NoError(t, err)
		payloadText := string(payload)
		assert.NotContains(t, payloadText, `"url":""`)
		assert.NotContains(t, payloadText, `"api_key":""`)
		assert.NotContains(t, payloadText, `"key":""`)
		assert.NotContains(t, payloadText, `"group_id":""`)

		props := req["properties"].(map[string]interface{})
		llm := props["llm"].(map[string]interface{})
		tts := props["tts"].(map[string]interface{})
		asr := props["asr"].(map[string]interface{})
		assert.NotContains(t, llm, "url")
		if params, ok := llm["params"].(map[string]interface{}); ok {
			assert.NotContains(t, params, "model")
		}
		assert.NotContains(t, tts["params"].(map[string]interface{}), "model")
		assert.NotContains(t, asr["params"].(map[string]interface{}), "model")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent_id":"agent_123","status":"RUNNING"}`))
	}))
	defer server.Close()

	rawClient := client.NewClient(
		option.WithBaseURL(server.URL),
		option.WithBasicAuth("user", "pass"),
		option.WithMaxAttempts(1),
	)
	agoraClient := &AgoraClient{
		Agents:         rawClient.Agents,
		AppID:          "appid",
		AppCertificate: "app-cert",
		AuthMode:       AuthModeBasic,
	}

	agent := NewAgent(WithName("managed-agent")).
		WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
			Model:    "nova-3",
			Language: "en",
		})).
		WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
			Model: "gpt-4o-mini",
		})).
		WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
			Model:   "speech_2_6_turbo",
			VoiceID: "English_captivating_female1",
		}))

	session := agent.CreateSession(agoraClient, CreateSessionOptions{
		Channel:    "room-1",
		AgentUID:   "1",
		RemoteUIDs: []string{"100"},
		Token:      "rtc-token",
	})

	agentID, err := session.Start(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "agent_123", agentID)
}

func TestOffRemovesRegisteredHandler(t *testing.T) {
	session := NewAgentSession(AgentSessionOptions{
		Client:     nil,
		Agent:      NewAgent(),
		AppID:      "appid",
		Name:       "agent",
		Channel:    "room",
		AgentUID:   "1",
		RemoteUIDs: []string{"2"},
	})

	var count int
	handler := func(data interface{}) { count++ }
	session.On("started", handler)
	session.Off("started", handler)
	session.emit("started", map[string]string{"agent_id": "agent"})
	assert.Equal(t, 0, count)
}

func TestGeminiLiveMatchesTypeScriptShape(t *testing.T) {
	maxHistory := 8
	config := vendors.NewGeminiLive(vendors.GeminiLiveOptions{
		APIKey:           "google-key",
		Model:            "gemini-live-2.5-flash",
		URL:              "wss://generativelanguage.googleapis.com/ws",
		Instructions:     "Be concise.",
		Voice:            "Aoede",
		GreetingMessage:  "Hello from Gemini",
		FailureMessage:   "Please try again.",
		MaxHistory:       &maxHistory,
		PredefinedTools:  []string{"_publish_message"},
		InputModalities:  []string{"audio"},
		OutputModalities: []string{"text", "audio"},
		Messages: []map[string]interface{}{
			{"role": "system", "content": "short memory"},
		},
		AdditionalParams: map[string]interface{}{
			"temperature": 0.2,
		},
	}).ToConfig()

	assert.Equal(t, map[string]interface{}{
		"vendor":  "gemini",
		"style":   "openai",
		"api_key": "google-key",
		"url":     "wss://generativelanguage.googleapis.com/ws",
		"params": map[string]interface{}{
			"temperature":  0.2,
			"model":        "gemini-live-2.5-flash",
			"instructions": "Be concise.",
			"voice":        "Aoede",
		},
		"messages": []map[string]interface{}{
			{"role": "system", "content": "short memory"},
		},
		"greeting_message":  "Hello from Gemini",
		"failure_message":   "Please try again.",
		"max_history":       8,
		"predefined_tools":  []string{"_publish_message"},
		"input_modalities":  []string{"audio"},
		"output_modalities": []string{"text", "audio"},
	}, config)
}

func TestMLLMWrappersIncludeOptionalFields(t *testing.T) {
	openAIMaxHistory := 3
	openAIConfig := vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
		APIKey:          "key",
		URL:             "wss://openai.example.com/realtime",
		PredefinedTools: []string{"_publish_message"},
		FailureMessage:  "Retry",
		MaxHistory:      &openAIMaxHistory,
	}).ToConfig()
	assert.Equal(t, "wss://openai.example.com/realtime", openAIConfig["url"])
	assert.Equal(t, []string{"_publish_message"}, openAIConfig["predefined_tools"])
	assert.Equal(t, "Retry", openAIConfig["failure_message"])
	assert.Equal(t, 3, openAIConfig["max_history"])

	vertexMaxHistory := 5
	vertexConfig := vendors.NewVertexAI(vendors.VertexAIOptions{
		Model:               "gemini-live",
		URL:                 "wss://vertex.example.com/realtime",
		ProjectID:           "project",
		Location:            "us-central1",
		ADCredentialsString: "adc",
		PredefinedTools:     []string{"_publish_message"},
		FailureMessage:      "Try again",
		MaxHistory:          &vertexMaxHistory,
	}).ToConfig()
	assert.Equal(t, "wss://vertex.example.com/realtime", vertexConfig["url"])
	assert.Equal(t, []string{"_publish_message"}, vertexConfig["predefined_tools"])
	assert.Equal(t, "Try again", vertexConfig["failure_message"])
	assert.Equal(t, 5, vertexConfig["max_history"])
}

func TestPresetBackedOpenAIVendorsAllowMissingKeys(t *testing.T) {
	agent := NewAgent(WithInstructions("Preset-backed flow")).
		WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
			Model: "nova-3",
		})).
		WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
			Model: "gpt-5-mini",
		})).
		WithTts(vendors.NewOpenAITTS(vendors.OpenAITTSOptions{
			Voice: "alloy",
		}))

	props, err := agent.ToProperties(ToPropertiesOptions{
		Channel:              "room-1",
		AgentUID:             "1",
		RemoteUIDs:           []string{"100"},
		AppID:                "0123456789abcdef0123456789abcdef",
		AppCertificate:       "fedcba9876543210fedcba9876543210",
		SkipVendorValidation: true,
	})
	require.NoError(t, err)

	preset, resolved, err := ResolveSessionPresets(nil, props)
	require.NoError(t, err)
	assert.Equal(t, "deepgram_nova_3,openai_gpt_5_mini,openai_tts_1", preset)
	require.NotNil(t, resolved)

	payload, err := json.Marshal(resolved)
	require.NoError(t, err)
	assert.NotContains(t, string(payload), "api_key")
}

func TestPresetBackedMiniMaxTTSAllowsMissingKey(t *testing.T) {
	tts := vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
		Model:   "speech-2.6-turbo",
		VoiceID: "English_captivating_female1",
		URL:     "wss://api-uw.minimax.io/ws/v1/t2a_v2",
	}).ToConfig()

	assert.Equal(t, "minimax", tts["vendor"])
	params := tts["params"].(map[string]interface{})
	assert.Equal(t, "speech-2.6-turbo", params["model"])
	assert.NotContains(t, params, "key")
	assert.NotContains(t, params, "group_id")
	assert.Equal(t, map[string]interface{}{"voice_id": "English_captivating_female1"}, params["voice_setting"])
	assert.Equal(t, "wss://api-uw.minimax.io/ws/v1/t2a_v2", params["url"])
}

func TestManagedPresetPayloadOmitsProviderOwnedFields(t *testing.T) {
	smartFormat := true
	punctuation := true
	maxHistory := 15
	agent := NewAgent(
		WithInstructions("Preset flow"),
		WithGreeting("hello"),
	).WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
		Model:       "nova-3",
		Language:    "en",
		SmartFormat: &smartFormat,
		Punctuation: &punctuation,
	})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
		Model:      "gpt-4o-mini",
		MaxHistory: &maxHistory,
		Params: map[string]interface{}{
			"max_tokens":  1024,
			"temperature": 0.7,
			"top_p":       0.95,
		},
	})).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
		Model:   "speech_2_6_turbo",
		VoiceID: "English_captivating_female1",
	}))

	props, err := agent.ToPropertiesMap(ToPropertiesOptions{
		Channel:              "room-1",
		Token:                "rtc-token",
		AgentUID:             "1",
		RemoteUIDs:           []string{"100"},
		SkipVendorValidation: true,
	})
	require.NoError(t, err)

	preset, resolved, err := ResolveSessionPresetsMap(nil, props)
	require.NoError(t, err)
	assert.Equal(t, "deepgram_nova_3,openai_gpt_4o_mini,minimax_speech_2_6_turbo", preset)

	payload, err := json.Marshal(resolved)
	require.NoError(t, err)
	payloadText := string(payload)
	assert.NotContains(t, payloadText, `"url":""`)
	assert.NotContains(t, payloadText, `"api_key":""`)
	assert.NotContains(t, payloadText, `"key":""`)
	assert.NotContains(t, payloadText, `"group_id":""`)
	assert.NotContains(t, payloadText, `"model":"gpt-4o-mini"`)
	assert.NotContains(t, payloadText, `"model":"speech_2_6_turbo"`)
	assert.NotContains(t, payloadText, `"model":"nova-3"`)

	asr := resolved["asr"].(map[string]interface{})
	asrParams := asr["params"].(map[string]interface{})
	assert.Equal(t, "deepgram", asr["vendor"])
	assert.Equal(t, "en", asrParams["language"])
	assert.Equal(t, true, asrParams["smart_format"])
	assert.Equal(t, true, asrParams["punctuation"])
	assert.NotContains(t, asrParams, "api_key")
	assert.NotContains(t, asrParams, "model")

	llm := resolved["llm"].(map[string]interface{})
	llmParams := llm["params"].(map[string]interface{})
	assert.Equal(t, "openai", llm["style"])
	assert.Equal(t, 1024, llmParams["max_tokens"])
	assert.Equal(t, 0.7, llmParams["temperature"])
	assert.Equal(t, 0.95, llmParams["top_p"])
	assert.Equal(t, 15, llm["max_history"])
	assert.NotContains(t, llm, "url")
	assert.NotContains(t, llm, "api_key")
	assert.NotContains(t, llmParams, "model")

	tts := resolved["tts"].(map[string]interface{})
	ttsParams := tts["params"].(map[string]interface{})
	assert.Equal(t, "minimax", tts["vendor"])
	assert.Equal(t, map[string]interface{}{"voice_id": "English_captivating_female1"}, ttsParams["voice_setting"])
	assert.NotContains(t, ttsParams, "key")
	assert.NotContains(t, ttsParams, "group_id")
	assert.NotContains(t, ttsParams, "url")
	assert.NotContains(t, ttsParams, "model")
}

func TestManagedOpenAITTSOmitKeyAndModel(t *testing.T) {
	agent := NewAgent().WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
		Model: "nova-2",
	})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
		Model: "gpt-5-mini",
	})).WithTts(vendors.NewOpenAITTS(vendors.OpenAITTSOptions{
		Model: "tts-1",
		Voice: "alloy",
	}))

	props, err := agent.ToPropertiesMap(ToPropertiesOptions{
		Channel:              "room-1",
		Token:                "rtc-token",
		AgentUID:             "1",
		RemoteUIDs:           []string{"100"},
		SkipVendorValidation: true,
	})
	require.NoError(t, err)

	preset, resolved, err := ResolveSessionPresetsMap(nil, props)
	require.NoError(t, err)
	assert.Equal(t, "deepgram_nova_2,openai_gpt_5_mini,openai_tts_1", preset)

	tts := resolved["tts"].(map[string]interface{})
	params := tts["params"].(map[string]interface{})
	assert.Equal(t, "alloy", params["voice"])
	assert.NotContains(t, params, "key")
	assert.NotContains(t, params, "api_key")
	assert.NotContains(t, params, "model")
}

func TestAresASRRemainsKeylessWithoutPreset(t *testing.T) {
	agent := NewAgent().WithStt(vendors.NewAresSTT(vendors.AresSTTOptions{
		Language: "en-US",
		AdditionalParams: map[string]interface{}{
			"sample_rate": 16000,
		},
	})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
		Model: "gpt-4o-mini",
	})).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
		Model:   "speech_2_8_turbo",
		VoiceID: "English_captivating_female1",
	}))

	props, err := agent.ToPropertiesMap(ToPropertiesOptions{
		Channel:              "room-1",
		Token:                "rtc-token",
		AgentUID:             "1",
		RemoteUIDs:           []string{"100"},
		SkipVendorValidation: true,
	})
	require.NoError(t, err)

	preset, resolved, err := ResolveSessionPresetsMap(nil, props)
	require.NoError(t, err)
	assert.Equal(t, "openai_gpt_4o_mini,minimax_speech_2_8_turbo", preset)

	asr := resolved["asr"].(map[string]interface{})
	params := asr["params"].(map[string]interface{})
	assert.Equal(t, "ares", asr["vendor"])
	assert.Equal(t, "en-US", asr["language"])
	assert.Equal(t, 16000, params["sample_rate"])
	assert.NotContains(t, params, "api_key")
	assert.NotContains(t, params, "key")
	assert.NotContains(t, params, "model")
}

func TestBYOKProvidersAreNotTreatedAsManagedPresets(t *testing.T) {
	agent := NewAgent().WithStt(vendors.NewDeepgramSTT(vendors.DeepgramSTTOptions{
		APIKey: "deepgram-key",
		Model:  "nova-3",
	})).WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
		APIKey: "openai-key",
		Model:  "gpt-4o-mini",
	})).WithTts(vendors.NewMiniMaxTTS(vendors.MiniMaxTTSOptions{
		Key:     "minimax-key",
		GroupID: "minimax-group",
		Model:   "speech_2_6_turbo",
		VoiceID: "English_captivating_female1",
	}))

	props, err := agent.ToPropertiesMap(ToPropertiesOptions{
		Channel:    "room-1",
		Token:      "rtc-token",
		AgentUID:   "1",
		RemoteUIDs: []string{"100"},
	})
	require.NoError(t, err)

	preset, resolved, err := ResolveSessionPresetsMap(nil, props)
	require.NoError(t, err)
	assert.Empty(t, preset)

	asrParams := resolved["asr"].(map[string]interface{})["params"].(map[string]interface{})
	llm := resolved["llm"].(map[string]interface{})
	ttsParams := resolved["tts"].(map[string]interface{})["params"].(map[string]interface{})
	assert.Equal(t, "deepgram-key", asrParams["api_key"])
	assert.Equal(t, "openai-key", llm["api_key"])
	assert.Equal(t, "minimax-key", ttsParams["key"])
	assert.Equal(t, "minimax-group", ttsParams["group_id"])
}

func TestToPropertiesBubblesMLLMFieldsAndPreservesVendorOverrides(t *testing.T) {
	enableMllm := true
	maxHistory := 9
	agent := NewAgent(
		WithGreeting("Agent greeting"),
		WithFailureMessage("Agent failure"),
		WithMaxHistory(maxHistory),
		WithAdvancedFeatures(&AdvancedFeatures{EnableMllm: &enableMllm}),
	).WithMllm(vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
		APIKey:          "openai-key",
		Model:           "gpt-4o-realtime-preview",
		URL:             "wss://openai.example.com/realtime",
		GreetingMessage: "Vendor greeting",
		PredefinedTools: []string{"_publish_message"},
	}))

	props, err := agent.ToProperties(ToPropertiesOptions{
		Channel:    "room",
		Token:      "rtc-token",
		AgentUID:   "1",
		RemoteUIDs: []string{"100"},
	})
	require.NoError(t, err)
	require.NotNil(t, props)
	require.NotNil(t, props.Mllm)

	payload, err := json.Marshal(props.Mllm)
	require.NoError(t, err)
	assert.Contains(t, string(payload), "greeting_message")
	assert.Contains(t, string(payload), "failure_message")
	assert.Contains(t, string(payload), "max_history")
	assert.Contains(t, string(payload), "predefined_tools")
	assert.Contains(t, string(payload), "url")

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(payload, &decoded))
	assert.Equal(t, "Vendor greeting", decoded["greeting_message"])
	assert.Equal(t, "Agent failure", decoded["failure_message"])
	assert.Equal(t, float64(9), decoded["max_history"])
	assert.Equal(t, []interface{}{"_publish_message"}, decoded["predefined_tools"])
	assert.Equal(t, "wss://openai.example.com/realtime", decoded["url"])
}

func TestAvatarHelpersCoverLiveAvatarAndAnam(t *testing.T) {
	assert.True(t, IsLiveAvatarAvatar("liveavatar"))
	assert.True(t, IsAnamAvatar("anam"))
	require.NoError(t, ValidateAvatarConfig("liveavatar", map[string]interface{}{
		"api_key":   "live-key",
		"quality":   "high",
		"agora_uid": "42",
	}))
	require.NoError(t, ValidateAvatarConfig("anam", map[string]interface{}{
		"api_key": "anam-key",
	}))
	require.NoError(t, ValidateTtsSampleRate("liveavatar", 24000))
	require.Error(t, ValidateTtsSampleRate("liveavatar", 16000))

	avatar := vendors.NewAnamAvatar(vendors.AnamAvatarOptions{
		APIKey:    "anam-key",
		PersonaID: "persona-1",
	}).ToConfig()
	assert.Equal(t, "anam", avatar["vendor"])
}

func TestSessionWarnsForAvatarWithoutExplicitSampleRateAndSupportsWarnHook(t *testing.T) {
	var warnings []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/projects/appid/join":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"agent_id":"agent_123","status":"RUNNING"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	rawClient := client.NewClient(
		option.WithBaseURL(server.URL),
		option.WithBasicAuth("user", "pass"),
		option.WithMaxAttempts(1),
	)

	agent := NewAgent(WithName("avatar-agent")).
		WithLlm(vendors.NewOpenAI(vendors.OpenAIOptions{
			APIKey: "openai-key",
			Model:  "gpt-4o-mini",
		})).
		WithTts(vendors.NewMicrosoftTTS(vendors.MicrosoftTTSOptions{
			Key:       "ms-key",
			Region:    "eastus",
			VoiceName: "en-US-JennyNeural",
		})).
		WithAvatar(vendors.NewLiveAvatarAvatar(vendors.LiveAvatarAvatarOptions{
			APIKey:   "live-key",
			Quality:  "high",
			AgoraUID: "42",
		}))

	session := NewAgentSession(AgentSessionOptions{
		Client:     rawClient.Agents,
		Agent:      agent,
		AppID:      "appid",
		Name:       "avatar-agent",
		Channel:    "room",
		Token:      "rtc-token",
		AgentUID:   "1",
		RemoteUIDs: []string{"2"},
		Warn: func(msg string) {
			warnings = append(warnings, msg)
		},
	})

	_, err := session.Start(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "LiveAvatar")
}

func TestSessionWarnHookReceivesHandlerPanics(t *testing.T) {
	var warnings []string

	session := NewAgentSession(AgentSessionOptions{
		Client:     nil,
		Agent:      NewAgent(),
		AppID:      "appid",
		Name:       "agent",
		Channel:    "room",
		AgentUID:   "1",
		RemoteUIDs: []string{"2"},
		Warn: func(msg string) {
			warnings = append(warnings, msg)
		},
	})

	session.On("started", func(data interface{}) {
		panic("boom")
	})
	session.emit("started", map[string]string{"agent_id": "agent"})
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "recovered panic")
}
