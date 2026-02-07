package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	edgetts "github.com/hecx333/edge-tts-go"
)

func main() {
	text := flag.String("text", "Hello, Edge TTS!", "Text to synthesize")
	voice := flag.String("voice", "en-US-GuyNeural", "Voice to use (e.g., en-US-GuyNeural, zh-CN-XiaoxiaoNeural)")
	output := flag.String("output", "output.mp3", "Output file path")
	rate := flag.String("rate", "+0%", "Speech rate (e.g. +10%, -10%)")
	volume := flag.String("volume", "+0%", "Volume (e.g. +10%, -10%)")
	pitch := flag.String("pitch", "+0Hz", "Pitch (e.g. +5Hz, -5Hz)")

	flag.Parse()

	if *text == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Synthesizing '%s' with voice '%s'...\n", *text, *voice)

	tts := edgetts.NewTTS(
		edgetts.WithVoice(*voice),
		edgetts.WithRate(*rate),
		edgetts.WithVolume(*volume),
		edgetts.WithPitch(*pitch),
	)

	audioData, err := tts.Speak(*text)
	if err != nil {
		log.Fatalf("Error synthesizing speech: %v", err)
	}

	if err := os.WriteFile(*output, audioData, 0644); err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	fmt.Printf("Audio saved to %s\n", *output)
}
