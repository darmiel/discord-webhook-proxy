package discord

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Webhook struct {
	UID        string      `bson:"uid" json:"uid"`
	UserID     string      `bson:"user_id" json:"user_id"`
	Secret     string      `bson:"secret" json:"secret"`
	WebhookURL string      `bson:"webhook_url" json:"webhook_url"`
	Data       WebhookData `bson:"data" json:"data"`
}

// NewWebhook creates a new webhook and generates a secret and a UID
func NewWebhook(userID string, uid string, webhookURL string, secret string, data WebhookData) (w *Webhook) {
	if uid == "" {
		uid = uuid.New().String()
	}

	if secret == "" {
		secret = GenerateSecret()
	}

	return &Webhook{
		UID:        uid,
		UserID:     userID,
		Secret:     secret,
		WebhookURL: webhookURL,
		Data:       data,
	}
}

func (w *Webhook) ParseNewLine() {
	w.Data = WebhookData(strings.ReplaceAll(string(w.Data), "\n", ""))
}

func (w *Webhook) GetID() string {
	return w.UserID + ":" + w.UID
}

func (w *Webhook) CreateFilter(includeSecret bool) (filter bson.M) {
	params := []bson.M{
		{"uid": w.UID},
		{"user_id": w.UserID},
	}
	if includeSecret {
		params = append(params, bson.M{"secret": w.Secret})
	}
	filter = bson.M{
		"$and": params,
	}
	return
}

// Send sends the webhook directly to discord without any further validation checks
// so be sure to check the Webhook before calling Send
func (w *Webhook) Send(param ...interface{}) (sentJson string, err error) {
	// parse data
	data, err := w.Data.Exec(param...)
	if err != nil {
		return "", err
	}
	// send data
	sentJson, err = w.SendJson(data)
	return
}

func (w *Webhook) SendJson(json string) (sent string, err error) {
	sent = json

	// Redis Stats
	go w.AddCallStats()

	reader := bytes.NewReader([]byte(json))
	var req *http.Request
	req, err = http.NewRequest("POST", w.WebhookURL, reader)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return sent, err
	}

	s := math.Floor(float64(resp.StatusCode) / 100)
	if s != float64(2) {
		return sent, errors.New("status was not 2xx but " + strconv.Itoa(resp.StatusCode))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return
}
