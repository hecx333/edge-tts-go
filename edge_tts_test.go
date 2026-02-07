package edgetts

import (
	"os"
	"testing"
)

// TestNewTTS validates the default configuration.
func TestNewTTS(t *testing.T) {
	tts := NewTTS()

	if tts.cfg.Voice != DefaultVoice {
		t.Errorf("Expected default voice %s, got %s", DefaultVoice, tts.cfg.Voice)
	}
	if tts.cfg.Rate != DefaultRate {
		t.Errorf("Expected default rate %s, got %s", DefaultRate, tts.cfg.Rate)
	}
}

// TestWithConfig validates custom configuration application.
func TestWithConfig(t *testing.T) {
	voice := "en-GB-RyanNeural"
	rate := "+10%"
	volume := "+20%"
	pitch := "+5Hz"

	tts := NewTTS(
		WithVoice(voice),
		WithRate(rate),
		WithVolume(volume),
		WithPitch(pitch),
	)

	if tts.cfg.Voice != voice {
		t.Errorf("Expected voice %s, got %s", voice, tts.cfg.Voice)
	}
	if tts.cfg.Rate != rate {
		t.Errorf("Expected rate %s, got %s", rate, tts.cfg.Rate)
	}
	if tts.cfg.Volume != volume {
		t.Errorf("Expected volume %s, got %s", volume, tts.cfg.Volume)
	}
	if tts.cfg.Pitch != pitch {
		t.Errorf("Expected pitch %s, got %s", pitch, tts.cfg.Pitch)
	}
}

// TestSynthesize checks if the synthesis returns data.
// Note: This is an integration test that hits the actual Edge TTS service.
// It might fail without internet connection.
func TestSynthesize(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") != "" {
		t.Skip("Skipping integration test")
	}

	text := "Hello, this is a test from Go Edge TTS."
	tts := NewTTS()

	audioData, err := tts.Speak(text)
	if err != nil {
		t.Fatalf("Synthesize failed: %v", err)
	}

	if len(audioData) == 0 {
		t.Fatal("Received 0 bytes of audio data")
	}

	// Verify it looks like an MP3 (simple check for ID3 or frame header is complex, so length > 100 is good enough for now)
	if len(audioData) < 100 {
		t.Errorf("Audio data seems too small: %d bytes", len(audioData))
	}
}
