// Conformance tests for the Deepgram TTS provider using omnivoice providertest.
// These tests verify the provider correctly implements the omnivoice TTS interfaces.
//
// To run tests, set DEEPGRAM_API_KEY environment variable:
//
//	source /path/to/agentcall/.envrc
//	go test -v -run TestConformance
//
// Note: Deepgram TTS provider requires API key for creation,
// so all tests require the DEEPGRAM_API_KEY environment variable.
package tts

import (
	"os"
	"testing"

	"github.com/agentplexus/omnivoice/tts/providertest"
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

	providertest.RunAll(t, providertest.Config{
		Provider:          p,
		StreamingProvider: p, // Deepgram TTS implements StreamingProvider
		SkipIntegration:   false,
		TestVoiceID:       "aura-asteria-en", // Default Deepgram voice
		TestText:          "Hello, this is a test of the Deepgram text to speech provider.",
	})
}

// TestConformance_InterfaceOnly runs only interface tests.
// Note: Deepgram TTS requires API key for provider creation.
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

// TestConformance_Behavior runs behavior tests (may require API).
func TestConformance_Behavior(t *testing.T) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set, skipping behavior tests")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	providertest.RunBehaviorTests(t, providertest.Config{
		Provider:        p,
		SkipIntegration: false,
		TestVoiceID:     "aura-asteria-en",
	})
}

// TestConformance_Integration runs only integration tests (requires API).
func TestConformance_Integration(t *testing.T) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		t.Skip("DEEPGRAM_API_KEY not set, skipping integration tests")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	providertest.RunIntegrationTests(t, providertest.Config{
		Provider:          p,
		StreamingProvider: p,
		SkipIntegration:   false,
		TestVoiceID:       "aura-asteria-en",
		TestText:          "Hello, this is a test of the Deepgram text to speech provider.",
	})
}

// TestProviderRequiresAPIKey verifies that provider creation fails without API key.
func TestProviderRequiresAPIKey(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Error("New() should return error when API key is not provided")
	}
}
