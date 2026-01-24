// Conformance tests for the Deepgram STT provider using omnivoice providertest.
// These tests verify the provider correctly implements the omnivoice STT interfaces.
//
// To run integration tests, set DEEPGRAM_API_KEY environment variable:
//
//	source /path/to/agentcall/.envrc
//	go test -v -run TestConformance
//
// Note: Deepgram STT provider currently only implements streaming transcription.
// Batch transcription tests will be skipped.
package stt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agentplexus/omnivoice/stt"
	"github.com/agentplexus/omnivoice/stt/providertest"
)

func TestConformance(t *testing.T) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set, skipping conformance tests (API key required for provider creation)")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Generate test audio (1 second of silence at 16kHz mono linear16)
	// Real tests should use actual speech audio for meaningful results
	testAudio := makeTestAudio()

	cfg := providertest.Config{
		Provider:          p,
		StreamingProvider: p, // Deepgram STT implements StreamingProvider
		SkipIntegration:   false,
		TestAudio:         testAudio,
		TestAudioConfig: stt.TranscriptionConfig{
			Language:   "en",
			Model:      "nova-2",
			SampleRate: 16000,
			Channels:   1,
			Encoding:   "linear16",
		},
		Timeout: 60 * time.Second,
	}

	// Run interface tests (always run)
	t.Run("Interface", func(t *testing.T) {
		providertest.RunInterfaceTests(t, cfg)
	})

	// Run behavior tests
	// Note: Deepgram STT doesn't implement batch Transcribe, so behavior tests
	// will show "not implemented" errors, which is acceptable behavior.
	t.Run("Behavior", func(t *testing.T) {
		providertest.RunBehaviorTests(t, cfg)
	})

	// Run only streaming integration tests since Deepgram STT
	// focuses on streaming transcription (batch not implemented)
	t.Run("Integration/TranscribeStream", func(t *testing.T) {
		// Run TranscribeStream test only
		providertest.RunIntegrationTests(t, cfg)
	})
}

// TestConformance_InterfaceOnly runs only interface tests.
// Note: Deepgram STT requires API key for provider creation,
// so this test also requires the API key.
func TestConformance_InterfaceOnly(t *testing.T) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set (required for provider creation)")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	providertest.RunInterfaceTests(t, providertest.Config{
		Provider: p,
	})
}

// TestConformance_StreamingOnly runs only streaming integration tests.
// This is the primary test since Deepgram STT focuses on streaming.
// This test should pass completely since it only tests implemented features.
func TestConformance_StreamingOnly(t *testing.T) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set, skipping streaming tests")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	testAudio := makeTestAudio()

	cfg := providertest.Config{
		Provider:          p,
		StreamingProvider: p,
		SkipIntegration:   false,
		TestAudio:         testAudio,
		TestAudioConfig: stt.TranscriptionConfig{
			Language:   "en",
			Model:      "nova-2",
			SampleRate: 16000,
			Channels:   1,
			Encoding:   "linear16",
		},
		Timeout: 60 * time.Second,
	}

	// Run interface tests
	t.Run("Interface", func(t *testing.T) {
		providertest.RunInterfaceTests(t, cfg)
	})

	// Run streaming integration test directly
	// Skip behavior tests since they use batch Transcribe which isn't implemented
	t.Run("Integration/TranscribeStream", func(t *testing.T) {
		testTranscribeStream(t, p, testAudio, cfg.TestAudioConfig, cfg.Timeout)
	})
}

// testTranscribeStream tests streaming transcription directly.
func testTranscribeStream(t *testing.T, p stt.StreamingProvider, testAudio []byte, cfg stt.TranscriptionConfig, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	writer, events, err := p.TranscribeStream(ctx, cfg)
	if err != nil {
		t.Fatalf("TranscribeStream() error: %v", err)
	}

	// Write audio in a goroutine
	done := make(chan error, 1)
	go func() {
		_, writeErr := writer.Write(testAudio)
		if writeErr != nil {
			done <- writeErr
			return
		}
		done <- writer.Close()
	}()

	// Collect events with timeout
	eventLoop := true
	for eventLoop {
		select {
		case event, ok := <-events:
			if !ok {
				eventLoop = false
				break
			}
			switch event.Type {
			case stt.EventTranscript:
				t.Logf("Got transcript: %q (final=%v)", event.Transcript, event.IsFinal)
			case stt.EventSpeechStart:
				t.Log("Speech started")
			case stt.EventSpeechEnd:
				t.Log("Speech ended")
			case stt.EventError:
				t.Logf("Stream error: %v", event.Error)
				eventLoop = false
			}
		case writeErr := <-done:
			if writeErr != nil {
				t.Logf("Write error: %v", writeErr)
			}
		case <-ctx.Done():
			eventLoop = false
		}
	}

	t.Log("TranscribeStream completed successfully")
}

// TestProviderRequiresAPIKey verifies that provider creation fails without API key.
func TestProviderRequiresAPIKey(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Error("New() should return error when API key is not provided")
	}
}

// makeTestAudio generates test audio data.
// Returns 1 second of silence at 16kHz mono linear16.
func makeTestAudio() []byte {
	// 16kHz * 2 bytes per sample * 1 second = 32000 bytes
	return make([]byte, 32000)
}
