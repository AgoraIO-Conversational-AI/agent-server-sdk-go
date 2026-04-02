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
	config := vendors.NewGeminiLive(vendors.GeminiLiveOptions{
		APIKey:           "google-key",
		Model:            "gemini-live-2.5-flash",
		Instructions:     "Be concise.",
		Voice:            "Aoede",
		GreetingMessage:  "Hello from Gemini",
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
		"input_modalities":  []string{"audio"},
		"output_modalities": []string{"text", "audio"},
	}, config)
}

func TestToPropertiesDoesNotBubbleUnsupportedMLLMFields(t *testing.T) {
	enableMllm := true
	maxHistory := 9
	agent := NewAgent(
		WithGreeting("Agent greeting"),
		WithFailureMessage("Agent failure"),
		WithMaxHistory(maxHistory),
		WithAdvancedFeatures(&AdvancedFeatures{EnableMllm: &enableMllm}),
	).WithMllm(vendors.NewOpenAIRealtime(vendors.OpenAIRealtimeOptions{
		APIKey: "openai-key",
		Model:  "gpt-4o-realtime-preview",
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
	assert.NotContains(t, string(payload), "failure_message")
	assert.NotContains(t, string(payload), "max_history")
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
