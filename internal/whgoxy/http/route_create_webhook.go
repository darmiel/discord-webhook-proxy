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
	UserID     string `json:"user_id"`
	Secret     string `json:"secret"`
	WebhookURL string `json:"webhook_url"`
	Payload    string `json:"payload"`
	Args       string `json:"args"`
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

	// print data
	log.Println("Data:", string(all))

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

	// check if the users are the same
	if data.UserID != user.DiscordUser.UserID {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "UserID mismatch (%v <-> %v)",
			data.UserID,
			user.DiscordUser.UserID,
		)
		return
	}

	// get database connection
	db := ws.Database

	// check if webhook already exists
	if data.UID != "" {
		// check validity of uid
		if err := discord.CheckUIDValidity(data.UID); err != nil {
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, "The UID is invalid. Expression: %s", discord.UIDExpr)
			return
		}

		// check for duplicates
		if _, err := db.FindWebhook(data.UID, user.DiscordUser.UserID); err == nil {
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, "A webhook with the same UID already exists: %s", data.UID)
			return
		}
	}

	// create webhook
	webhook := discord.NewWebhook(
		user.DiscordUser.UserID,
		data.UID,
		data.WebhookURL,
		data.Secret,
		discord.WebhookData(data.Payload),
	)

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
