package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func (ws *WebServer) homeRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("home", writer, request, nil)
}

func (ws *WebServer) createRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("create", writer, request, nil)
}

func (ws *WebServer) error404(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("err_404", writer, request, nil)
}

func (ws *WebServer) safeWebhookRouteHandler(w http.ResponseWriter, r *http.Request) {
	// parse vars
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	secret := vars["secret"]

	// get database connection
	db := ws.Database

	// try to get webhook with uuid
	webhook, err := db.FindWebhook(uuid)
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
