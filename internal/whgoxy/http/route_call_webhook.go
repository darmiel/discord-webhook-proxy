package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// /call/json/{userid}/{uid}/{secret}
func (ws *WebServer) jsonWebhookRouteHandler(w http.ResponseWriter, r *http.Request) {
	// parse vars
	vars := mux.Vars(r)
	userID := vars["userid"]
	uid := vars["uid"]
	secret := vars["secret"]

	log.Println("🌱 Requested JSON webhook", uid, "by", userID, "with secret", secret)

	// get database connection
	db := ws.Database

	// try to get webhook with uid
	webhook, err := db.FindWebhook(uid, userID)
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
	var params interface{}

	// read body
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, "Error (Request) reading your webhook: %s", err.Error())
		return
	}
	log.Println("💿 POST Data:", string(all))

	// try to unmarshall
	if err := json.Unmarshal(all, &params); err != nil {
		w.WriteHeader(205)
		log.Println("🚨 Error unmarshalling:", err)
	} else {
		w.WriteHeader(200)
	}

	// send webhook
	sentJson, err := webhook.Send(params)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error (Discord): %s", err.Error())
		if e := webhook.AddError(err, sentJson); e != nil {
			log.Println("📈 Error saving stat (error):", err)
		}
		return
	}

	_, _ = fmt.Fprintf(w, "Success: %s", sentJson)
	if e := webhook.AddSuccess(); e != nil {
		log.Println("📈 Error saving stat (success):", err)
	}

	log.Printf("Stats: %+v\n", webhook.GetStats())
}
