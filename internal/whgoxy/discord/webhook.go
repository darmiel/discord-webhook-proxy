package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type WebhookData string

type WebhookStats struct {
	SuccessfulRequests uint64 `bson:"successful_requests"`
	ErroredRequests    uint64 `bson:"errored_requests"`
}

type Webhook struct {
	UID        string       `bson:"uid" json:"uid"`
	UserID     string       `bson:"user_id" json:"user_id"`
	Secret     string       `bson:"secret" json:"secret"`
	WebhookURL string       `bson:"webhook_url" json:"webhook_url"`
	Data       WebhookData  `bson:"data" json:"data"`
	Stats      WebhookStats `bson:"stats" json:"stats"`
}

// NewWebhook creates a new webhook and generates a secret and a UID
func NewWebhook(userID string, uid string, webhookURL string, secret string, data WebhookData) (w *Webhook) {
	if uid == "" {
		uid = uuid.New().String()
	}

	if secret == "" {
		secret = generateSecret()
	}

	return &Webhook{
		UID:        uid,
		UserID:     userID,
		Secret:     secret,
		WebhookURL: webhookURL,
		Data:       data,
		Stats: WebhookStats{
			SuccessfulRequests: 0,
			ErroredRequests:    0,
		},
	}
}

func (w *Webhook) ParseNewLine() {
	w.Data = WebhookData(strings.ReplaceAll(string(w.Data), "\n", ""))
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
	// replace params in data
	var parse *template.Template
	parse, err = template.New("").Parse(string(w.Data))
	if err != nil {
		return
	}

	// parse data
	var data interface{}
	if param != nil && len(param) >= 1 {
		data = param[0]
		if j, err := json.Marshal(data); err != nil {
			log.Println("🌚 Webhook got data (as raw):", data)
		} else {
			log.Println("🌚 Webhook got data (as json):", string(j))
		}
	} else {
		log.Println("🌚 Webhook got empty data.")
	}

	// execute template
	var buffer bytes.Buffer
	if err = parse.Execute(&buffer, data); err != nil {
		return
	}

	// read string from buffer
	jsd := buffer.String()
	log.Println("👉 Sending data to webhook:", jsd)

	// Send to discord
	return w.sendJson(jsd)
}

func (w *Webhook) sendJson(jsd string) (sentJson string, err error) {
	sentJson = jsd
	reader := bytes.NewReader([]byte(jsd))

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
		return jsd, err
	}

	s := math.Floor(float64(resp.StatusCode) / 100)
	if s != float64(2) {
		return "", errors.New("status was not 2xx but " + strconv.Itoa(resp.StatusCode))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return jsd, nil
}
