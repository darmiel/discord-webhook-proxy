package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"io/ioutil"
	"log"
	"net/http"
)

type DeleteWebhookPayload struct {
	UID    string `json:"uid"`
	UserID string `json:"user_id"`
}

type DeleteWebhookResponse struct {
	Succes bool   `json:"succes"`
	Error  string `json:"error"`
}

func killReq(writer http.ResponseWriter, data interface{}) {
	var header = 400 // error

	switch t := data.(type) {
	case error:
		data = DeleteWebhookResponse{
			Succes: false,
			Error:  t.Error(),
		}
		break
	case DeleteWebhookResponse:
		if t.Succes {
			header = 200
		}
		break
	case string:
		data = DeleteWebhookResponse{
			Succes: false,
			Error:  t,
		}
	}

	// write (error) header
	writer.WriteHeader(header)

	// write error
	if data, err := json.Marshal(data); err != nil {
		// write "hard"-json
		_, _ = fmt.Fprintln(writer, `{ "success": false, "error": "error parsing response" }`)
	} else {
		_, _ = fmt.Fprintln(writer, string(data))
	}
}

func (ws *WebServer) deleteWebhookRouteHandler(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in
	user, die := auth.GetUserOrDie(r, w)
	if die {
		return
	}

	if !user.DiscordUser.HasPermission(discord.PermissionWebhookDelete) {
		killReq(w, "You don't have permissions to delete a webhook.")
		return
	}

	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		killReq(w, DeleteWebhookResponse{
			Succes: false,
			Error:  err.Error(),
		})
		return
	}

	// print data
	log.Println("Data:", string(all))

	// "validate" json and create webhook data
	var data DeleteWebhookPayload
	if err := json.Unmarshal(all, &data); err != nil {
		killReq(w, err)
		return
	}

	// check if the users are the same
	if data.UserID != user.DiscordUser.UserID {
		killReq(w, errors.New("user id mismatch"))
		return
	}

	if err := discord.CheckUserIDValidity(data.UserID); err != nil {
		killReq(w, errors.New("invalid user id"))
		return
	}

	if err := discord.CheckUIDValidity(data.UID); err != nil {
		killReq(w, errors.New("invalid webhook id"))
		return
	}

	// get database connection
	db := ws.Database

	// check if webhook exists
	webhook, err := db.FindWebhook(data.UID, data.UserID)
	if err != nil || webhook == nil {
		killReq(w, err)
		return
	}

	// delete webhook
	err = db.DeleteWebhook(webhook.UID, webhook.UserID)
	if err != nil {
		killReq(w, err)
	} else {
		killReq(w, DeleteWebhookResponse{
			Succes: true,
			Error:  "",
		})
	}
}
