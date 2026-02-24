package vendors

type SampleRate int

const (
	SampleRate8kHz  SampleRate = 8000
	SampleRate16kHz SampleRate = 16000
	SampleRate22kHz SampleRate = 22050
	SampleRate24kHz SampleRate = 24000
	SampleRate44kHz SampleRate = 44100
	SampleRate48kHz SampleRate = 48000
)

type LLM interface {
	ToConfig() map[string]interface{}
}

type TTS interface {
	ToConfig() map[string]interface{}
	GetSampleRate() *SampleRate
}

type STT interface {
	ToConfig() map[string]interface{}
}

type MLLM interface {
	ToConfig() map[string]interface{}
}

type Avatar interface {
	ToConfig() map[string]interface{}
	RequiredSampleRate() SampleRate
}
