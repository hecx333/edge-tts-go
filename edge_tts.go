package edgetts

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hecx333/edge-tts-go/internal/ssml"
)

// Define constants
const (
	DefaultWSSURL      = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	DefaultClientToken = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	DefaultVoice       = "en-US-GuyNeural"
	DefaultRate        = "+0%"
	DefaultVolume      = "+0%"
	DefaultPitch       = "+0Hz"
)

// Option is a functional option for configuring the TTS instance.
type Option func(*config)

type config struct {
	Voice  string
	Rate   string
	Volume string
	Pitch  string
}

// WithVoice sets the voice for the TTS.
func WithVoice(v string) Option {
	return func(c *config) {
		c.Voice = v
	}
}

// WithRate sets the speaking rate (e.g. "+50%", "-10%").
func WithRate(r string) Option {
	return func(c *config) {
		c.Rate = r
	}
}

// WithVolume sets the volume (e.g. "+50%", "-10%").
func WithVolume(v string) Option {
	return func(c *config) {
		c.Volume = v
	}
}

// WithPitch sets the pitch (e.g. "+5Hz", "-5Hz").
func WithPitch(p string) Option {
	return func(c *config) {
		c.Pitch = p
	}
}

// TTS handles the communication with Edge TTS service.
type TTS struct {
	cfg config
}

// NewTTS creates a new TTS instance with optional configurations.
func NewTTS(opts ...Option) *TTS {
	cfg := config{
		Voice:  DefaultVoice,
		Rate:   DefaultRate,
		Volume: DefaultVolume,
		Pitch:  DefaultPitch,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &TTS{
		cfg: cfg,
	}
}

// Speak connects to the Edge TTS service and returns the synthesized audio data for the given text.
func (t *TTS) Speak(text string) ([]byte, error) {
	connID := strings.ReplaceAll(uuid.New().String(), "-", "")
	secMSGEC := ssml.GenerateSecMSGEC(DefaultClientToken)
	muid := ssml.GenerateMUID()

	url := fmt.Sprintf("%s&ConnectionId=%s&Sec-MS-GEC=%s&Sec-MS-GEC-Version=1-143.0.3650.75",
		DefaultWSSURL, connID, secMSGEC)

	headers := http.Header{}
	headers.Set("Pragma", "no-cache")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 Edg/143.0.0.0")
	headers.Set("Origin", "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold")
	headers.Set("Cookie", fmt.Sprintf("muid=%s;", muid))
	headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	headers.Set("Accept-Language", "en-US,en;q=0.9")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}
	defer conn.Close()

	if err := t.sendConfig(conn); err != nil {
		return nil, err
	}

	if err := t.sendSSML(conn, text); err != nil {
		return nil, err
	}

	return t.readAudio(conn)
}

func (t *TTS) sendConfig(conn *websocket.Conn) error {
	configPayload := map[string]interface{}{
		"context": map[string]interface{}{
			"synthesis": map[string]interface{}{
				"audio": map[string]interface{}{
					"metadataoptions": map[string]string{
						"sentenceBoundaryEnabled": "false",
						"wordBoundaryEnabled":     "false",
					},
					"outputFormat": "audio-24khz-48kbitrate-mono-mp3",
				},
			},
		},
	}
	configJSON, err := json.Marshal(configPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal config payload: %w", err)
	}

	configMsg := fmt.Sprintf("X-Timestamp:%s\r\n"+
		"Content-Type:application/json; charset=utf-8\r\n"+
		"Path:speech.config\r\n\r\n%s",
		ssml.DateToString(), string(configJSON))

	if err := conn.WriteMessage(websocket.TextMessage, []byte(configMsg)); err != nil {
		return fmt.Errorf("failed to send speech.config: %w", err)
	}
	return nil
}

func (t *TTS) sendSSML(conn *websocket.Conn, text string) error {
	s := ssml.BuildSSML(text, t.cfg.Voice, t.cfg.Rate, t.cfg.Volume, t.cfg.Pitch)
	requestID := ssml.GenerateRequestID()
	ssmlMsg := fmt.Sprintf("X-RequestId:%s\r\n"+
		"Content-Type:application/ssml+xml\r\n"+
		"X-Timestamp:%sZ\r\n"+
		"Path:ssml\r\n\r\n"+
		"%s",
		requestID, ssml.DateToString(), s)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(ssmlMsg)); err != nil {
		return fmt.Errorf("failed to send ssml: %w", err)
	}
	return nil
}

func (t *TTS) readAudio(conn *websocket.Conn) ([]byte, error) {
	var audioData []byte
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			return nil, fmt.Errorf("read error: %w", err)
		}

		if messageType == websocket.TextMessage {
			msg := string(p)
			if strings.Contains(msg, "Path:turn.end") {
				break
			}
		} else if messageType == websocket.BinaryMessage {
			if len(p) < 2 {
				continue
			}
			headerLen := int(binary.BigEndian.Uint16(p[:2]))
			skip := 2 + headerLen
			if len(p) > skip {
				audioData = append(audioData, p[skip:]...)
			}
		}
	}

	if len(audioData) == 0 {
		return nil, fmt.Errorf("no audio data received")
	}

	return audioData, nil
}
