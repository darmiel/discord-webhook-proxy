package discord

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
)

const (
	webhookURLMinLen = 33
	webhookURLMaxLen = 200
	dataMinLen       = 2 // {}
	dataMaxLen       = 2000
)

var (
	errorWebhookURLTooShort = errors.New("empty webhook too short")
	errorWebhookURLTooLong  = errors.New("empty webhook too long")
	errorUnknownWebhookURL  = errors.New("unknown webhook url")
	errorDataTooShort       = errors.New("data too short")
	errorDataTooLong        = errors.New("data too long")
)

var (
	discordURLRegex *regexp.Regexp
)

func init() {
	var err error
	discordURLRegex, err = regexp.Compile("https://((ptb|canary)\\.)?discord(app)?\\.com/api/webhooks/[0-9]+/[A-Za-z0-9-]+")
	if err != nil {
		log.Fatalln("Error compiling regex expression:", err)
		return
	}
	log.Println("Compiled regex:", discordURLRegex)
}

func (w *Webhook) CheckValidity(sendTest bool) (err error) {
	// check webhook length
	l := len(w.WebhookURL)
	if l < webhookURLMinLen {
		return errorWebhookURLTooShort
	} else if l > webhookURLMaxLen {
		return errorWebhookURLTooLong
	}

	// check webhook url
	if !discordURLRegex.Match([]byte(w.WebhookURL)) {
		return errorUnknownWebhookURL
	}
	log.Println("Matched regex exp")

	// check json validity
	if data, err := json.Marshal(w); err != nil {
		return err
	} else {
		// check data length
		l = len(data)
		if l < dataMinLen {
			return errorDataTooShort
		} else if l > dataMaxLen {
			return errorDataTooLong
		}
	}

	// make a test call to discord
	if sendTest {
		if _, err := w.Send(map[string]string{}); err != nil {
			return err
		}
	}

	return nil
}
