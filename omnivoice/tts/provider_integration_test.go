// Integration tests for the TTS provider.
// These tests require DEEPGRAM_API_KEY environment variable to be set.

package tts

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agentplexus/omnivoice/tts"
)

func getAPIKey(t *testing.T) string {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set, skipping integration test")
	}
	return apiKey
}

func TestSynthesizeIntegration(t *testing.T) {
	apiKey := getAPIKey(t)

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config := tts.SynthesisConfig{
		Model:        "aura-asteria-en",
		OutputFormat: "mp3",
		SampleRate:   24000,
	}

	result, err := p.Synthesize(ctx, "Hello, world! This is a test of the Deepgram TTS integration.", config)
	if err != nil {
		t.Fatalf("Synthesize() error = %v", err)
	}

	if result == nil {
		t.Fatal("Synthesize() returned nil result")
	}

	if len(result.Audio) == 0 {
		t.Error("Synthesize() returned empty audio")
	}

	if result.CharacterCount == 0 {
		t.Error("Synthesize() returned zero character count")
	}

	t.Logf("Synthesize() returned %d bytes of audio, %d characters processed", len(result.Audio), result.CharacterCount)
}

func TestSynthesizeStreamIntegration(t *testing.T) {
	apiKey := getAPIKey(t)

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config := tts.SynthesisConfig{
		Model:        "aura-asteria-en",
		OutputFormat: "linear16",
		SampleRate:   24000,
	}

	chunkCh, err := p.SynthesizeStream(ctx, "Hello, this is a streaming test.", config)
	if err != nil {
		t.Fatalf("SynthesizeStream() error = %v", err)
	}

	var totalBytes int
	var chunkCount int
	var gotFinal bool
	var streamErr error

	for chunk := range chunkCh {
		if chunk.Error != nil {
			streamErr = chunk.Error
			break
		}
		if len(chunk.Audio) > 0 {
			totalBytes += len(chunk.Audio)
			chunkCount++
		}
		if chunk.IsFinal {
			gotFinal = true
			cancel() // Cancel context to close the stream
		}
	}

	if streamErr != nil {
		t.Fatalf("SynthesizeStream() stream error = %v", streamErr)
	}

	if totalBytes == 0 {
		t.Error("SynthesizeStream() received no audio bytes")
	}

	if chunkCount == 0 {
		t.Error("SynthesizeStream() received no audio chunks")
	}

	if !gotFinal {
		t.Error("SynthesizeStream() did not receive final chunk")
	}

	t.Logf("SynthesizeStream() received %d bytes in %d chunks", totalBytes, chunkCount)
}

func TestSynthesizeWithDifferentVoices(t *testing.T) {
	apiKey := getAPIKey(t)

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	voices := []string{
		"aura-asteria-en",
		"aura-orion-en",
	}

	for _, voiceID := range voices {
		t.Run(voiceID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			config := tts.SynthesisConfig{
				VoiceID:      voiceID,
				OutputFormat: "mp3",
			}

			result, err := p.Synthesize(ctx, "Testing voice selection.", config)
			if err != nil {
				t.Errorf("Synthesize() with voice %s error = %v", voiceID, err)
				return
			}

			if len(result.Audio) == 0 {
				t.Errorf("Synthesize() with voice %s returned empty audio", voiceID)
			}
		})
	}
}

func TestSynthesizeWithDifferentFormats(t *testing.T) {
	apiKey := getAPIKey(t)

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	formats := []string{
		"mp3",
		"linear16",
		"mulaw",
	}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			config := tts.SynthesisConfig{
				Model:        "aura-asteria-en",
				OutputFormat: format,
			}

			result, err := p.Synthesize(ctx, "Testing output format.", config)
			if err != nil {
				t.Errorf("Synthesize() with format %s error = %v", format, err)
				return
			}

			if len(result.Audio) == 0 {
				t.Errorf("Synthesize() with format %s returned empty audio", format)
			}

			t.Logf("Format %s: %d bytes", format, len(result.Audio))
		})
	}
}
