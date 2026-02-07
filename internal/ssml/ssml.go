package ssml

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateSecMSGEC generates the security header value required by Edge TTS.
func GenerateSecMSGEC(clientToken string) string {
	ticks := time.Now().UTC().Unix()
	ticks += 11644473600
	ticks -= ticks % 300
	ticks *= 10000000

	strToHash := fmt.Sprintf("%d%s", ticks, clientToken)
	hash := sha256.Sum256([]byte(strToHash))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// GenerateMUID generates a random Message Unique ID.
func GenerateMUID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}

// GenerateRequestID generates a random Request ID.
func GenerateRequestID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}

// DateToString formats the current time for Edge TTS headers.
func DateToString() string {
	return time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)")
}

// BuildSSML constructs the SSML XML string for synthesis.
func BuildSSML(text, voice, rate, volume, pitch string) string {
	escapedText := html.EscapeString(text)
	return fmt.Sprintf(
		`<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'>`+
			`<voice name='%s'>`+
			`<prosody pitch='%s' rate='%s' volume='%s'>`+
			`%s`+
			`</prosody>`+
			`</voice>`+
			`</speak>`,
		voice, pitch, rate, volume, escapedText,
	)
}
