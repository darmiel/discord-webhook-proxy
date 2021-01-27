package http

import (
	"encoding/json"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CreateWebhookPayload struct {
	UID        string `json:"uid"`
	Secret     string `json:"secret"`
	WebhookURL string `json:"webhook_url"`
	Payload    string `json:"payload"`
	Args       string `json:"args"`
	Force      bool   `json:"force"`
}

type CreateWebhookResponse struct {
	Webhook  *discord.Webhook `json:"webhook"`
	SentJson string           `json:"sent_json"`
}

func (ws *WebServer) createWebhookRouteHandler(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	user, die := auth.GetUserOrDie(r, w)
	if die {
		return
	}

	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Request) reading your webhook: %s", err.Error())
		return
	}

	// "validate" json and create webhook data
	var data CreateWebhookPayload
	if err := json.Unmarshal(all, &data); err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Request) decoding your webhook: %s", err.Error())
		return
	}

	// read args
	var args interface{}
	if data.Args != "" {
		log.Println("😨 Data has args!")
		if strings.HasPrefix(data.Args, "{") && strings.HasSuffix(data.Args, "}") {
			if err := json.Unmarshal([]byte(data.Args), &args); err != nil {
				w.WriteHeader(400)
				_, _ = fmt.Fprintf(w, "Error: Example data could not be decoded")
				return
			}
		} else {
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, "Error: Example data must be in JSON format")
			return
		}
	} else {
		args = nil
	}

	// get database connection
	db := ws.Database

	var webhook *discord.Webhook

	// check if webhook already exists
	if data.UID != "" {
		// check validity of uid
		if err := discord.CheckUIDValidity(data.UID); err != nil {
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, "The UID is invalid. Expression: %s", discord.UIDExpr)
			return
		}

		// check for duplicates
		if wh, err := db.FindWebhook(data.UID, user.DiscordUser.UserID); err == nil {

			// check if forced request
			if data.Force {
				webhook = wh

				// update webhook
				if data.Payload != string(wh.Data) {
					webhook.Data = discord.WebhookData(data.Payload)
				}
				if data.WebhookURL != wh.WebhookURL {
					webhook.WebhookURL = data.WebhookURL
				}
				if data.Secret != wh.Secret {
					webhook.Secret = data.Secret
				}
			} else {
				w.WriteHeader(300)
				_, _ = fmt.Fprintf(w, "A webhook with the same UID already exists: %s", data.UID)
				return
			}
		}
	}

	// check secret
	if data.Secret != "" {
		// check secret validity
		if err := discord.CheckSecretValidity(data.Secret); err != nil {
			w.WriteHeader(400)
			_, _ = fmt.Fprint(w, "Secret is not valid")
			return
		}
	} else {
		// Generate new secret
		if webhook != nil {
			webhook.Secret = discord.GenerateSecret()
		}
	}

	if webhook == nil {
		// create webhook
		webhook = discord.NewWebhook(
			user.DiscordUser.UserID,
			data.UID,
			data.WebhookURL,
			data.Secret,
			discord.WebhookData(data.Payload),
		)
	}

	// validate webhook
	req, err := webhook.CheckValidityWithSend(args)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Webhook is not valid: %s | Sent Json: %s", err.Error(), req)
		return
	}

	// webhook is valid
	// -> Save to database

	if err := db.SaveWebhook(webhook); err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Webhook is vaid, but could not be saved due to a database error: %s", err.Error())
		return
	}

	w.WriteHeader(200)

	// create response
	response := &CreateWebhookResponse{
		Webhook:  webhook,
		SentJson: req,
	}

	// marshall webook
	if js, err := json.Marshal(response); err != nil {
		_, _ = fmt.Fprintf(w, "{}")
	} else {
		_, _ = fmt.Fprintf(w, string(js))
	}
}
