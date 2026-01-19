# OmniVoice Deepgram Provider

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

OmniVoice provider implementation for [Deepgram](https://deepgram.com/) speech-to-text services.

This package adapts the official [Deepgram Go SDK](https://github.com/deepgram/deepgram-go-sdk) to the [OmniVoice](https://github.com/agentplexus/omnivoice) interfaces, enabling Deepgram's speech-to-text capabilities within the OmniVoice framework.

## Features

- Real-time streaming transcription via WebSocket
- Support for telephony audio formats (mu-law, a-law)
- Interim and final transcription results
- Speech start/end detection for natural turn-taking
- Speaker diarization support
- Keyword boosting

## Installation

```bash
go get github.com/agentplexus/omnivoice-deepgram
```

## Usage

### Basic Streaming Transcription

```go
import (
    deepgramstt "github.com/agentplexus/omnivoice-deepgram/omnivoice/stt"
    "github.com/agentplexus/omnivoice/stt"
)

// Create provider with API key
provider, err := deepgramstt.New(deepgramstt.WithAPIKey("your-api-key"))
if err != nil {
    log.Fatal(err)
}

// Configure for telephony audio
config := stt.TranscriptionConfig{
    Model:      "nova-2",
    Language:   "en-US",
    Encoding:   "mulaw",
    SampleRate: 8000,
}

// Start streaming transcription
writer, events, err := provider.TranscribeStream(ctx, config)
if err != nil {
    log.Fatal(err)
}

// Send audio data
go func() {
    defer writer.Close()
    io.Copy(writer, audioSource)
}()

// Receive transcription events
for event := range events {
    switch event.Type {
    case stt.EventTranscript:
        if event.IsFinal {
            fmt.Println("Final:", event.Transcript)
        }
    case stt.EventSpeechStart:
        fmt.Println("Speech started")
    case stt.EventSpeechEnd:
        fmt.Println("Speech ended")
    case stt.EventError:
        log.Printf("Error: %v", event.Error)
    }
}
```

### With OmniVoice Pipeline

For a complete voice agent example using Deepgram STT with ElevenLabs TTS and Twilio Media Streams, see the [omnivoice-examples](https://github.com/agentplexus/omnivoice-examples) repository.

## Supported Audio Formats

| Format | Encoding Value | Typical Use |
|--------|---------------|-------------|
| mu-law | `mulaw` | Twilio, telephony |
| A-law | `alaw` | European telephony |
| Linear PCM | `linear16` | General audio |
| FLAC | `flac` | Compressed lossless |
| Opus | `opus` | WebRTC |
| MP3 | `mp3` | Compressed lossy |

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `Model` | Deepgram model | `nova-2` |
| `Language` | Language code | `en-US` |
| `SampleRate` | Audio sample rate | `8000` |
| `Channels` | Audio channels | `1` |
| `EnablePunctuation` | Add punctuation | `false` |
| `EnableSpeakerDiarization` | Identify speakers | `false` |
| `Keywords` | Words to boost | `[]` |

## Requirements

- Go 1.21 or later
- Deepgram API key ([get one here](https://console.deepgram.com/))

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Projects

- [omnivoice](https://github.com/agentplexus/omnivoice) - Voice agent framework interfaces
- [go-elevenlabs](https://github.com/agentplexus/go-elevenlabs) - ElevenLabs TTS provider
- [omnivoice-twilio](https://github.com/agentplexus/omnivoice-twilio) - Twilio Media Streams transport
- [omnivoice-examples](https://github.com/agentplexus/omnivoice-examples) - Complete voice agent examples

 [build-status-svg]: https://github.com/agentplexus/omnivoice-deepgram/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/agentplexus/omnivoice-deepgram/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/agentplexus/omnivoice-deepgram/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/agentplexus/omnivoice-deepgram/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/agentplexus/omnivoice-deepgram
 [goreport-url]: https://goreportcard.com/report/github.com/agentplexus/omnivoice-deepgram
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/agentplexus/omnivoice-deepgram
 [docs-godoc-url]: https://pkg.go.dev/github.com/agentplexus/omnivoice-deepgram
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/agentplexus/omnivoice-deepgram/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/agentplexus/omnivoice-deepgram/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/agentplexus/omnivoice-deepgram?badge
