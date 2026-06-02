package agentkit

import Agora "github.com/AgoraIO/agora-agents-go/v2"

// =============================================================================
// SAL mode constants
// Use these instead of the raw Agora.StartAgentsRequestPropertiesSalSalMode* values.
// =============================================================================

var (
	// SalModeLocking locks onto the primary speaker, suppressing background voices.
	SalModeLocking = Agora.StartAgentsRequestPropertiesSalSalModeLocking

	// SalModeRecognition identifies different speakers and suppresses background noise.
	SalModeRecognition = Agora.StartAgentsRequestPropertiesSalSalModeRecognition
)

// =============================================================================
// Data channel constants
// =============================================================================

var (
	// DataChannelRtm routes agent data via the RTM (Signaling) service.
	// Requires enable_rtm: true in AdvancedFeatures.
	DataChannelRtm = Agora.StartAgentsRequestPropertiesParametersDataChannelRtm

	// DataChannelDatastream routes agent data via the RTC data stream.
	DataChannelDatastream = Agora.StartAgentsRequestPropertiesParametersDataChannelDatastream
)

var (
	// AudioScenarioDefault maps to the default RTC audio scenario.
	AudioScenarioDefault = ParametersAudioScenarioDefault

	// AudioScenarioChorus is optimized for ultra-low-latency chorus use cases.
	AudioScenarioChorus = ParametersAudioScenarioChorus

	// AudioScenarioAIServer is optimized for conversational AI interaction reliability.
	AudioScenarioAIServer = ParametersAudioScenarioAIServer
)

// =============================================================================
// Silence action constants
// =============================================================================

var (
	// SilenceActionSpeak plays the silence prompt via TTS when the timeout elapses.
	SilenceActionSpeak = Agora.StartAgentsRequestPropertiesParametersSilenceConfigActionSpeak

	// SilenceActionThink appends the silence prompt to the LLM context instead of speaking it.
	SilenceActionThink = Agora.StartAgentsRequestPropertiesParametersSilenceConfigActionThink
)

// =============================================================================
// Geofence area constants
// =============================================================================

var (
	GeofenceAreaGlobal       = Agora.StartAgentsRequestPropertiesGeofenceAreaGlobal
	GeofenceAreaNorthAmerica = Agora.StartAgentsRequestPropertiesGeofenceAreaNorthAmerica
	GeofenceAreaEurope       = Agora.StartAgentsRequestPropertiesGeofenceAreaEurope
	GeofenceAreaAsia         = Agora.StartAgentsRequestPropertiesGeofenceAreaAsia
	GeofenceAreaIndia        = Agora.StartAgentsRequestPropertiesGeofenceAreaIndia
	GeofenceAreaJapan        = Agora.StartAgentsRequestPropertiesGeofenceAreaJapan

	GeofenceExcludeAreaNorthAmerica = Agora.StartAgentsRequestPropertiesGeofenceExcludeAreaNorthAmerica
	GeofenceExcludeAreaEurope       = Agora.StartAgentsRequestPropertiesGeofenceExcludeAreaEurope
	GeofenceExcludeAreaAsia         = Agora.StartAgentsRequestPropertiesGeofenceExcludeAreaAsia
	GeofenceExcludeAreaIndia        = Agora.StartAgentsRequestPropertiesGeofenceExcludeAreaIndia
	GeofenceExcludeAreaJapan        = Agora.StartAgentsRequestPropertiesGeofenceExcludeAreaJapan
)

// =============================================================================
// Turn detection type constants (deprecated; use Config.EndOfSpeech instead)
// =============================================================================

var (
	// TurnDetectionTypeAgoraVad uses Agora voice activity detection.
	TurnDetectionTypeAgoraVad = Agora.StartAgentsRequestPropertiesTurnDetectionTypeAgoraVad

	// TurnDetectionTypeServerVad uses server-side voice activity detection.
	TurnDetectionTypeServerVad = Agora.StartAgentsRequestPropertiesTurnDetectionTypeServerVad

	// TurnDetectionTypeSemanticVad uses semantic voice activity detection.
	TurnDetectionTypeSemanticVad = Agora.StartAgentsRequestPropertiesTurnDetectionTypeSemanticVad
)

// =============================================================================
// Filler words selection rule constants
// =============================================================================

var (
	// FillerWordsSelectionRuleShuffle plays filler words in random order (no repeats until all used).
	FillerWordsSelectionRuleShuffle = Agora.StartAgentsRequestPropertiesFillerWordsContentStaticConfigSelectionRuleShuffle

	// FillerWordsSelectionRuleRoundRobin plays filler words sequentially.
	FillerWordsSelectionRuleRoundRobin = Agora.StartAgentsRequestPropertiesFillerWordsContentStaticConfigSelectionRuleRoundRobin
)

// =============================================================================
// Think action constants
// =============================================================================

var (
	// ThinkOnListeningActionInject injects the text into the current listening turn.
	ThinkOnListeningActionInject = Agora.AgentThinkAgentManagementRequestOnListeningActionInject

	// ThinkOnListeningActionInterrupt interrupts the current listening flow and starts a new turn.
	ThinkOnListeningActionInterrupt = Agora.AgentThinkAgentManagementRequestOnListeningActionInterrupt

	// ThinkOnListeningActionIgnore ignores the think request while listening.
	ThinkOnListeningActionIgnore = Agora.AgentThinkAgentManagementRequestOnListeningActionIgnore

	// ThinkOnThinkingActionInterrupt interrupts the current thinking state and starts a new turn.
	ThinkOnThinkingActionInterrupt = Agora.AgentThinkAgentManagementRequestOnThinkingActionInterrupt

	// ThinkOnThinkingActionIgnore ignores the think request while thinking.
	ThinkOnThinkingActionIgnore = Agora.AgentThinkAgentManagementRequestOnThinkingActionIgnore

	// ThinkOnSpeakingActionInterrupt interrupts the current speaking state and starts a new turn.
	ThinkOnSpeakingActionInterrupt = Agora.AgentThinkAgentManagementRequestOnSpeakingActionInterrupt

	// ThinkOnSpeakingActionIgnore ignores the think request while speaking.
	ThinkOnSpeakingActionIgnore = Agora.AgentThinkAgentManagementRequestOnSpeakingActionIgnore
)

// =============================================================================
// Interruption configuration constants
// =============================================================================

var (
	// InterruptionModeStartOfSpeech triggers interruption when the user starts speaking.
	InterruptionModeStartOfSpeech = Agora.StartAgentsRequestPropertiesInterruptionModeStartOfSpeech

	// InterruptionModeKeywords triggers interruption only when the user speaks a configured keyword.
	InterruptionModeKeywords = Agora.StartAgentsRequestPropertiesInterruptionModeKeywords

	// InterruptionDisabledStrategyAppend queues user speech until the agent finishes the current turn (interruption.enable=false).
	InterruptionDisabledStrategyAppend = Agora.StartAgentsRequestPropertiesInterruptionDisabledConfigStrategyAppend

	// InterruptionDisabledStrategyIgnore drops user speech while the agent is speaking or thinking (interruption.enable=false).
	InterruptionDisabledStrategyIgnore = Agora.StartAgentsRequestPropertiesInterruptionDisabledConfigStrategyIgnore
)

// =============================================================================
// Speak priority constants (used by AgentSession.Say)
// =============================================================================

var (
	// SpeakPriorityInterrupt interrupts the current agent output and starts speaking immediately.
	SpeakPriorityInterrupt = Agora.SpeakAgentsRequestPriorityInterrupt

	// SpeakPriorityAppend queues the message after the current agent output finishes.
	SpeakPriorityAppend = Agora.SpeakAgentsRequestPriorityAppend

	// SpeakPriorityIgnore drops the message if the agent is currently speaking.
	SpeakPriorityIgnore = Agora.SpeakAgentsRequestPriorityIgnore
)

// =============================================================================
// MLLM turn detection mode constants
// =============================================================================

var (
	// MllmTurnDetectionModeAgoraVad uses Agora-provided VAD for MLLM turn detection.
	MllmTurnDetectionModeAgoraVad = Agora.StartAgentsRequestPropertiesMllmTurnDetectionModeAgoraVad

	// MllmTurnDetectionModeServerVad delegates VAD to the MLLM vendor.
	MllmTurnDetectionModeServerVad = Agora.StartAgentsRequestPropertiesMllmTurnDetectionModeServerVad

	// MllmTurnDetectionModeSemanticVad uses semantic detection (OpenAI Realtime only).
	MllmTurnDetectionModeSemanticVad = Agora.StartAgentsRequestPropertiesMllmTurnDetectionModeSemanticVad
)
