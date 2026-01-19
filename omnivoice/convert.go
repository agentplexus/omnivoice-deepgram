package omnivoice

import (
	"time"

	"github.com/agentplexus/omnivoice/stt"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
)

// ConfigToLiveTranscriptionOptions converts OmniVoice TranscriptionConfig to Deepgram options.
func ConfigToLiveTranscriptionOptions(config stt.TranscriptionConfig) *interfaces.LiveTranscriptionOptions {
	opts := &interfaces.LiveTranscriptionOptions{
		// Audio format
		Encoding:   mapEncoding(config.Encoding),
		SampleRate: config.SampleRate,
		Channels:   config.Channels,

		// Model and language
		Model:    config.Model,
		Language: config.Language,

		// Features
		Punctuate:   config.EnablePunctuation,
		SmartFormat: true, // Enable smart formatting
	}

	// Set defaults for telephony if not specified
	if opts.SampleRate == 0 {
		opts.SampleRate = 8000
	}
	if opts.Channels == 0 {
		opts.Channels = 1
	}
	if opts.Model == "" {
		opts.Model = "nova-2" // Best general model
	}
	if opts.Language == "" {
		opts.Language = "en-US"
	}

	// Enable interim results for streaming
	opts.InterimResults = true

	// Enable utterance detection for natural turn-taking
	opts.UtteranceEndMs = "1000" // 1 second silence = end of utterance

	// Enable diarization if requested
	if config.EnableSpeakerDiarization {
		opts.Diarize = true
		if config.MaxSpeakers > 0 {
			opts.DiarizeVersion = "latest"
		}
	}

	// Add keywords for boosting
	if len(config.Keywords) > 0 {
		opts.Keywords = config.Keywords
	}

	return opts
}

// mapEncoding maps OmniVoice encoding names to Deepgram encoding strings.
func mapEncoding(encoding string) string {
	switch encoding {
	case "mulaw", "ulaw", "g711u", "pcm_mulaw":
		return "mulaw"
	case "alaw", "g711a", "pcm_alaw":
		return "alaw"
	case "linear16", "pcm", "pcm_s16le":
		return "linear16"
	case "flac":
		return "flac"
	case "opus":
		return "opus"
	case "speex":
		return "speex"
	case "mp3":
		return "mp3"
	case "webm":
		return "webm"
	default:
		// Default to linear16 for PCM
		if encoding == "" {
			return "linear16"
		}
		return encoding
	}
}

// MessageResponseToStreamEvent converts a Deepgram MessageResponse to an OmniVoice stream event.
func MessageResponseToStreamEvent(result *MessageResponse) stt.StreamEvent {
	if result == nil || len(result.Channel.Alternatives) == 0 {
		return stt.StreamEvent{Type: stt.EventTranscript}
	}

	alt := result.Channel.Alternatives[0]

	event := stt.StreamEvent{
		Transcript: alt.Transcript,
		IsFinal:    result.IsFinal,
		Type:       stt.EventTranscript,
	}

	// Convert words if available
	if len(alt.Words) > 0 {
		segment := &stt.Segment{
			Text:       alt.Transcript,
			Confidence: float64(alt.Confidence),
		}

		for _, w := range alt.Words {
			word := stt.Word{
				Text:       w.Word,
				Confidence: float64(w.Confidence),
				StartTime:  time.Duration(w.Start * float64(time.Second)),
				EndTime:    time.Duration(w.End * float64(time.Second)),
			}

			// Include speaker if diarization is enabled
			if w.Speaker != nil {
				word.Speaker = formatSpeaker(*w.Speaker)
			}

			segment.Words = append(segment.Words, word)
		}

		// Set segment timing from first and last word
		if len(segment.Words) > 0 {
			segment.StartTime = segment.Words[0].StartTime
			segment.EndTime = segment.Words[len(segment.Words)-1].EndTime
		}

		event.Segment = segment
	}

	return event
}

// MessageResponse mirrors the Deepgram MessageResponse structure.
// This allows us to decouple from Deepgram's internal types.
type MessageResponse struct {
	Channel  Channel `json:"channel,omitempty"`
	IsFinal  bool    `json:"is_final,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Start    float64 `json:"start,omitempty"`
}

// Channel represents a transcription channel.
type Channel struct {
	Alternatives []Alternative `json:"alternatives,omitempty"`
}

// Alternative represents a transcription alternative.
type Alternative struct {
	Transcript string  `json:"transcript,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Words      []Word  `json:"words,omitempty"`
}

// Word represents a transcribed word with timing.
type Word struct {
	Word       string  `json:"word,omitempty"`
	Start      float64 `json:"start,omitempty"`
	End        float64 `json:"end,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Speaker    *int    `json:"speaker,omitempty"`
}

// formatSpeaker formats a speaker ID for OmniVoice.
func formatSpeaker(speaker int) string {
	return "speaker_" + itoa(speaker)
}

// itoa converts an int to string without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	var b [20]byte
	n := len(b)
	for i > 0 {
		n--
		b[n] = byte('0' + i%10)
		i /= 10
	}
	return string(b[n:])
}
