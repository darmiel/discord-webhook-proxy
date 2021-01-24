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
	errorInvalidUID         = errors.New("invalid uid")
	errorInvalidUserID      = errors.New("invalid userid")
	errorInvalidSecret      = errors.New("invalid secret")
)

const (
	DiscordURLExpr = "https://((ptb|canary)\\.)?discord(app)?\\.com/api/webhooks/[0-9]+/[A-Za-z0-9-]+"
	UIDExpr        = "^[a-zA-Z0-9_-]{1,36}$"
	UserIDExpr     = "^[0-9]{18}$"
	SecretExpr     = "^[A-Za-z0-9-_.]{6,64}$"
)

var (
	DiscordURLRegex *regexp.Regexp
	UIDRegex        *regexp.Regexp
	UserIDRegex     *regexp.Regexp
	SecretRegex     *regexp.Regexp
)

func init() {
	var err error
	DiscordURLRegex, err = regexp.Compile(DiscordURLExpr)
	if err != nil {
		log.Fatalln("Error compiling regex expression:", err)
		return
	}
	UIDRegex, err = regexp.Compile(UIDExpr)
	if err != nil {
		log.Fatalln("Error compiling UID expression:", err)
		return
	}
	UserIDRegex, err = regexp.Compile(UserIDExpr)
	if err != nil {
		log.Fatalln("Error compiling UserID expression:", err)
		return
	}
	SecretRegex, err = regexp.Compile(SecretExpr)
	if err != nil {
		log.Fatalln("Error compiling Secret expression:", err)
		return
	}
	log.Println("Compiled regex:", DiscordURLRegex)
}

func (w *Webhook) CheckValidity() (err error) {
	// check webhook length
	l := len(w.WebhookURL)
	if l < webhookURLMinLen {
		return errorWebhookURLTooShort
	} else if l > webhookURLMaxLen {
		return errorWebhookURLTooLong
	}

	// check webhook url
	if !DiscordURLRegex.Match([]byte(w.WebhookURL)) {
		return errorUnknownWebhookURL
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

	return nil
}

func (w *Webhook) CheckValidityWithSend(testData interface{}) (req string, err error) {
	// check validity
	if err := w.CheckValidity(); err != nil {
		return "", err
	}

	if testData == nil {
		testData = make(map[string]interface{})
	}

	req, err = w.Send(testData)

	return
}

func CheckUIDValidity(uid string) (err error) {
	if !UIDRegex.Match([]byte(uid)) {
		return errorInvalidUID
	}
	return nil
}

func CheckUserIDValidity(userID string) (err error) {
	if !UserIDRegex.Match([]byte(userID)) {
		return errorInvalidUserID
	}
	return nil
}

func CheckSecretValidity(secret string) (err error) {
	if !SecretRegex.Match([]byte(secret)) {
		return errorInvalidSecret
	}
	return nil
}
