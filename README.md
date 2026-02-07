# edge-tts-go

[English](README.md) | [中文](README_CN.md)

A Go library for Microsoft Edge's online text-to-speech service.
This library allows you to use the high-quality Neural TTS voices from Microsoft Edge for free.



## Features

- **Functional Options Pattern**: Clean and idiomatic Go API.
- **Support for Multiple Voices**: Access to all Microsoft Edge online voices.
- **Customizable**: Adjustable rate, pitch, and volume.
- **CLI Tool**: Includes a ready-to-use command-line interface.

## Installation

```bash
go get github.com/hecx333/edge-tts-go
```

## Usage

### Library

```go
package main

import (
	"log"
	"os"

	edgetts "github.com/hecx333/edge-tts-go"
)

func main() {
	// Simple usage with default settings
	tts := edgetts.NewTTS()
	
	// Advanced usage with options
	// tts := edgetts.NewTTS(
	// 	edgetts.WithVoice("zh-CN-XiaoxiaoNeural"),
	// 	edgetts.WithRate("+10%"),
	// 	edgetts.WithVolume("+20%"),
	// )

	audio, err := tts.Speak("Hello world")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("output.mp3", audio, 0644)
}
```


### CLI Tool

You can use the included CLI tool directly with `go run`, or build it first.

#### fast run

```bash
go run cmd/edge-tts/main.go -text "Hello, world" -voice "en-US-GuyNeural"
```

#### build and run

```bash
# Build the tool
go build -o edge-tts cmd/edge-tts/main.go

# Run it
./edge-tts -text "Hello world" -voice "en-US-GuyNeural" -output "hello.mp3"
```

#### Available Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-text` | Text to synthesize | "Hello, Edge TTS!" |
| `-voice` | Voice to use | "en-US-GuyNeural" |
| `-output` | Output file path | "output.mp3" |
| `-rate` | Speech rate (e.g. +10%, -10%) | "+0%" |
| `-volume` | Volume (e.g. +10%, -10%) | "+0%" |
| `-pitch` | Pitch (e.g. +5Hz, -5Hz) | "+0Hz" |

## Project Structure

- `cmd/`: Command-line application entry point.
- `internal/`: Internal logic (e.g., SSML building), not exposed.
- `examples/`: Example usage code.
- `edge_tts.go`: Core library code.

## Acknowledgements

- [https://github.com/rany2/edge-tts](https://github.com/rany2/edge-tts)

## License

MIT

## Disclaimer

This project is for educational and research purposes only. The availability and stability of the service cannot be guaranteed. Users assume all risks associated with the use of this tool. The authors are not responsible for any consequences arising from its use.
