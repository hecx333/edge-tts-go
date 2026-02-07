package main

import (
	"log"
	"os"

	edgetts "github.com/hecx333/edge-tts-go"
)

func main() {
	// Simple usage
	tts := edgetts.NewTTS()
	text := "Hello, this is a simple example."
	audio, err := tts.Speak(text)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	os.WriteFile("simple.mp3", audio, 0644)
	log.Println("simple.mp3 created")

	// Advanced usage
	advancedTTS := edgetts.NewTTS(
		edgetts.WithVoice("zh-CN-XiaoxiaoNeural"),
		edgetts.WithRate("+10%"),
		edgetts.WithVolume("+10%"),
		edgetts.WithPitch("+5Hz"),
	)
	chineseText := "你好，这是一个高级示例。"
	audio2, err := advancedTTS.Speak(chineseText)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	os.WriteFile("advanced.mp3", audio2, 0644)
	log.Println("advanced.mp3 created")
}
