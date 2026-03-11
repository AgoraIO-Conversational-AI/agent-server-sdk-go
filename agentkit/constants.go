package agentkit

import Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"

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
	// TurnDetectionTypeServerVad uses server-side voice activity detection.
	TurnDetectionTypeServerVad = Agora.StartAgentsRequestPropertiesTurnDetectionTypeServerVad
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
