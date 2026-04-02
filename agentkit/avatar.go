package agentkit

import "fmt"

func IsHeyGenAvatar(vendor string) bool {
	return vendor == "heygen"
}

func IsAkoolAvatar(vendor string) bool {
	return vendor == "akool"
}

func IsLiveAvatarAvatar(vendor string) bool {
	return vendor == "liveavatar"
}

func IsAnamAvatar(vendor string) bool {
	return vendor == "anam"
}

func ValidateAvatarConfig(vendor string, params map[string]interface{}) error {
	if IsHeyGenAvatar(vendor) || IsLiveAvatarAvatar(vendor) {
		label := "HeyGen"
		if IsLiveAvatarAvatar(vendor) {
			label = "LiveAvatar"
		}
		if params == nil {
			return fmt.Errorf("%s avatar requires params", label)
		}
		if _, ok := params["api_key"]; !ok {
			return fmt.Errorf("%s avatar requires api_key", label)
		}
		if q, ok := params["quality"]; !ok {
			return fmt.Errorf("%s avatar requires quality (low, medium, or high)", label)
		} else {
			qs, _ := q.(string)
			if qs != "low" && qs != "medium" && qs != "high" {
				return fmt.Errorf("invalid quality for %s: %v. Must be one of: low, medium, high", label, q)
			}
		}
		if _, ok := params["agora_uid"]; !ok {
			return fmt.Errorf("%s avatar requires agora_uid", label)
		}
	} else if IsAkoolAvatar(vendor) {
		if params == nil {
			return fmt.Errorf("Akool avatar requires params")
		}
		if _, ok := params["api_key"]; !ok {
			return fmt.Errorf("Akool avatar requires api_key")
		}
	} else if IsAnamAvatar(vendor) {
		if params == nil {
			return fmt.Errorf("Anam avatar requires params")
		}
		if _, ok := params["api_key"]; !ok {
			return fmt.Errorf("Anam avatar requires api_key")
		}
	}
	return nil
}

func ValidateTtsSampleRate(avatarVendor string, sampleRate int) error {
	if IsHeyGenAvatar(avatarVendor) || IsLiveAvatarAvatar(avatarVendor) {
		label := "HeyGen"
		docURL := "https://docs.agora.io/en/conversational-ai/models/avatar/heygen"
		if IsLiveAvatarAvatar(avatarVendor) {
			label = "LiveAvatar"
			docURL = "https://docs.agora.io/en/conversational-ai/models/avatar/overview"
		}
		if sampleRate != 24000 {
			return fmt.Errorf(
				"%s avatars ONLY support 24,000 Hz sample rate. "+
					"Your TTS is configured with %d Hz. "+
					"Please update your TTS configuration to use 24kHz sample rate. "+
					"See: %s",
				label, sampleRate, docURL,
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
