package discord

import (
	"encoding/json"
	"errors"
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
	errorDataTooShort       = errors.New("data too short")
	errorDataTooLong        = errors.New("data too long")
)

func (w *SavedWebhook) CheckValidity(sendTest bool) (err error) {
	// check webhook length
	l := len(w.WebhookURL)
	if l < webhookURLMinLen {
		return errorWebhookURLTooShort
	} else if l > webhookURLMaxLen {
		return errorWebhookURLTooLong
	}

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
