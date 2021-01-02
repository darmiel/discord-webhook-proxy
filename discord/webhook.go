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

type SavedWebhook struct {
	UUID       string       `bson:"uuid"`
	Secret     string       `bson:"secret"`
	WebhookURL string       `bson:"webhook_url"`
	Data       *WebhookData `bson:"data"`
}

// NewWebhook creates a new webhook and generates a secret and a UUID
func NewWebhook(webhookURL string, data *WebhookData) (w *SavedWebhook) {
	u := uuid.New().String()
	s := generateSecret()

	return &SavedWebhook{
		u,
		s,
		webhookURL,
		data,
	}
}

// Send sends the webhook directly to discord without any further validation checks
// so be sure to check the SavedWebhook before calling Send
func (w *SavedWebhook) Send(param ...map[string]string) (err error) {
	// marshal data
	jsdb, err := json.Marshal(w.Data)
	if err != nil {
		return err
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
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	s := math.Floor(float64(resp.StatusCode) / 100)
	if s != float64(2) {
		return errors.New("status was not 2xx but " + strconv.Itoa(resp.StatusCode))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}
