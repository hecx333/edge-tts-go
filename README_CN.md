# edge-tts-go

[English](README.md) | [中文](README_CN.md)

一个用于 Microsoft Edge 在线文本转语音服务的 Go 语言库。
本项目允许您免费使用 Microsoft Edge 的高质量神经 TTS 语音。



## 特性

- **Functional Options 模式**：提供清晰、符合 Go 语言习惯的 API。
- **多语音支持**：支持 Microsoft Edge 在线提供的多种高质量语音。
- **可定制性**：支持调节语速（Rate）、音调（Pitch）和音量（Volume）。
- **命令行工具**：内置开箱即用的命令行工具（CLI）。

## 安装

```bash
go get github.com/hecx333/edge-tts-go
```

## 使用指南

### 作为库使用

```go
package main

import (
	"log"
	"os"

	edgetts "github.com/hecx333/edge-tts-go"
)

func main() {
	// 简单用法，使用默认设置
	tts := edgetts.NewTTS()
	
	// 高级用法，配置自定义选项
	// tts := edgetts.NewTTS(
	// 	edgetts.WithVoice("zh-CN-XiaoxiaoNeural"),
	// 	edgetts.WithRate("+10%"),
	// 	edgetts.WithVolume("+20%"),
	// )

	text := "你好，世界"
	audio, err := tts.Speak(text)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("output.mp3", audio, 0644)
}
```


### 命令行工具 (CLI)

您可以直接使用 `go run` 运行内置的命令行工具，或者先编译再运行。

#### 快速运行

```bash
go run cmd/edge-tts/main.go -text "你好，世界" -voice "zh-CN-XiaoxiaoNeural"
```

#### 编译运行

```bash
# 编译工具
go build -o edge-tts cmd/edge-tts/main.go

# 运行工具
./edge-tts -text "你好，世界" -voice "zh-CN-XiaoxiaoNeural" -output "hello.mp3"
```

#### 可用参数

| 参数 | 描述 | 默认值 |
|------|-------------|---------|
| `-text` | 要转换的文本 | "Hello, Edge TTS!" |
| `-voice` | 使用的语音角色 | "en-US-GuyNeural" |
| `-output` | 输出文件路径 | "output.mp3" |
| `-rate` | 语速 (例如 +10%, -10%) | "+0%" |
| `-volume` | 音量 (例如 +10%, -10%) | "+0%" |
| `-pitch` | 音调 (例如 +5Hz, -5Hz) | "+0Hz" |

## 目录结构

- `cmd/`: 包含命令行应用程序入口。
- `internal/`: 包含内部使用的辅助逻辑（如 SSML 构建），不对外暴露。
- `examples/`:包含使用示例代码。
- `edge_tts.go`: 核心库代码。

## 致谢

- [https://github.com/rany2/edge-tts](https://github.com/rany2/edge-tts)

## 许可证

MIT

## 免责声明

本项目仅供学习和研究使用。无法保证服务的可用性和稳定性。用户需自行承担使用本工具的所有风险。作者不对因使用本工具而产生的任何后果负责。
