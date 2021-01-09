package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func (ws *WebServer) safeWebhookRouteHandler(w http.ResponseWriter, r *http.Request) {
	// parse vars
	vars := mux.Vars(r)
	userID := vars["user_id"]
	uid := vars["uid"]
	secret := vars["secret"]

	log.Println("Requested webhook", uid, "by", userID, "with secret", secret)

	// get database connection
	db := ws.Database

	// try to get webhook with uid
	webhook, err := db.FindWebhook(uid, "abc")
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error (Database): %s", err.Error())
		return
	}

	// check secret
	if webhook.Secret != secret {
		_, _ = fmt.Fprintf(w, "Error (Auth): Secret mismatch")
		return
	}

	// build params
	params := make(map[string]string)
	query := r.URL.Query()
	for k, v := range query {
		if strings.HasSuffix(k, "[]") {
			k = k[:len(k)-2]
		}

		params[k] = strings.Join(v, ",")
		log.Println("[Debug] Param", k, "=", params[k])
	}

	// send webhook
	sentJson, err := webhook.Send(params)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error (Discord): %s", err.Error())
		return
	}

	_, _ = fmt.Fprintf(w, "Success: %s", sentJson)
}
