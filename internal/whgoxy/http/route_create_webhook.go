package http

import (
	"encoding/json"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"io/ioutil"
	"log"
	"net/http"
)

type CreateWebhookPayload struct {
	UID        string `json:"uid"`
	UserID     string `json:"user_id"`
	Secret     string `json:"secret"`
	WebhookURL string `json:"webhook_url"`
	Payload    string `json:"payload"`
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

	// check if the users are the same
	if data.UserID != user.DiscordUser.UserID {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Request): UserID mismatch (%v <-> %v)",
			data.UserID,
			user.DiscordUser.UserID,
		)
		return
	}

	// get database connection
	db := ws.Database

	// check if webhook already exists
	if _, err := db.FindWebhook(data.UID, user.DiscordUser.UserID); err == nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Webhook): A webhook with the same UID already exists: %s", data.UID)
		return
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
	err = webhook.CheckValidity(true)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Webhook): Webhook is not vaid: %s", err.Error())
		return
	}

	// webhook is valid
	// -> Save to database

	if err := db.SaveWebhook(webhook); err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Webhook): Webhook is vaid, but could not store in database: %s", err.Error())
		return
	}

	w.WriteHeader(200)

	// marshall webook
	if js, err := json.Marshal(webhook); err != nil {
		_, _ = fmt.Fprintf(w, "{}")
	} else {
		_, _ = fmt.Fprintf(w, string(js))
	}
}
