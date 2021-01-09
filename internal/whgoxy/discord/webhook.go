package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type WebhookData bson.M

type WebhookStats struct {
	SuccessfulRequests uint64 `bson:"successful_requests"`
	ErroredRequests    uint64 `bson:"errored_requests"`
}

type Webhook struct {
	UID        string        `bson:"uid"`
	UserID     string        `bson:"user_id"`
	Secret     string        `bson:"secret"`
	WebhookURL string        `bson:"webhook_url"`
	Data       *WebhookData  `bson:"data"`
	Stats      *WebhookStats `bson:"stats"`
}

// NewWebhook creates a new webhook and generates a secret and a UID
func NewWebhook(userID string, uid string, webhookURL string, secret string, data *WebhookData) (w *Webhook) {
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
		Stats: &WebhookStats{
			SuccessfulRequests: 0,
			ErroredRequests:    0,
		},
	}
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
func (w *Webhook) Send(param ...map[string]string) (sentJson string, err error) {
	// marshal data
	jsdb, err := json.Marshal(w.Data)
	if err != nil {
		return "", err
	}

	jsd := string(jsdb)

	// replace params in data
	if param != nil && len(param) >= 1 && len(param[0]) > 0 {
		for key, value := range param[0] {
			re := strings.NewReplacer(
				fmt.Sprintf("{{%s}}", key), value,
				fmt.Sprintf("{{ %s }}", key), value,

				fmt.Sprintf("{{ %s}}", key), value, // also \       / "faulty" \            /
				fmt.Sprintf("{{%s }}", key), value, //       replace            placeholders
			)

			jsd = re.Replace(jsd)
		}
	}

	// Send to discord
	reader := bytes.NewReader([]byte(jsd))
	req, err := http.NewRequest("POST", w.WebhookURL, reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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
