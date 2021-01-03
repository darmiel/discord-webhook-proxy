package router

import (
	"fmt"
	"github.com/darmiel/whgoxy/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func New(db db.Database) (res *mux.Router) {
	router := mux.NewRouter()

	// call_s -> Call Safe (mainly for debugging purposes)
	// = send a seperate message if the secrets were different
	router.HandleFunc("/call_s/{uuid}/{secret}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["uuid"]
		secret := vars["secret"]

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
		for k, _ := range query {
			params[k] = query.Get(k)
			log.Println("[Debug] Param", k, "=", query.Get(k))
		}

		// send webhook
		if err := webhook.Send(params); err != nil {
			_, _ = fmt.Fprintf(w, "Error (Discord): %s", err.Error())
			return
		}

		_, _ = fmt.Fprintf(w, "Success")
	})

	return router
}
