package agentkit

import "fmt"

func IsHeyGenAvatar(vendor string) bool {
	return vendor == "heygen"
}

func IsAkoolAvatar(vendor string) bool {
	return vendor == "akool"
}

func ValidateAvatarConfig(vendor string, params map[string]interface{}) error {
	if IsHeyGenAvatar(vendor) {
		if params == nil {
			return fmt.Errorf("HeyGen avatar requires params")
		}
		if _, ok := params["api_key"]; !ok {
			return fmt.Errorf("HeyGen avatar requires api_key")
		}
		if q, ok := params["quality"]; !ok {
			return fmt.Errorf("HeyGen avatar requires quality (low, medium, or high)")
		} else {
			qs, _ := q.(string)
			if qs != "low" && qs != "medium" && qs != "high" {
				return fmt.Errorf("invalid quality for HeyGen: %v. Must be one of: low, medium, high", q)
			}
		}
		if _, ok := params["agora_uid"]; !ok {
			return fmt.Errorf("HeyGen avatar requires agora_uid")
		}
	} else if IsAkoolAvatar(vendor) {
		if params == nil {
			return fmt.Errorf("Akool avatar requires params")
		}
		if _, ok := params["api_key"]; !ok {
			return fmt.Errorf("Akool avatar requires api_key")
		}
	}
	return nil
}

func ValidateTtsSampleRate(avatarVendor string, sampleRate int) error {
	if IsHeyGenAvatar(avatarVendor) {
		if sampleRate != 24000 {
			return fmt.Errorf(
				"HeyGen avatars ONLY support 24,000 Hz sample rate. "+
					"Your TTS is configured with %d Hz. "+
					"Please update your TTS configuration to use 24kHz sample rate. "+
					"See: https://docs.agora.io/en/conversational-ai/models/avatar/heygen",
				sampleRate,
			)
		}
	} else if IsAkoolAvatar(avatarVendor) {
		if sampleRate != 16000 {
			return fmt.Errorf(
				"Akool avatars ONLY support 16,000 Hz sample rate. "+
					"Your TTS is configured with %d Hz. "+
					"Please update your TTS configuration to use 16kHz sample rate. "+
					"See: https://docs.agora.io/en/conversational-ai/models/avatar/akool",
				sampleRate,
			)
		}
	}
	return nil
}
