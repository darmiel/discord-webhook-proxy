package http

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func (ws *WebServer) editWebhookHandler(writer http.ResponseWriter, request *http.Request) {
	// get login
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	vars := mux.Vars(request)
	uid := vars["uid"]

	userID := user.DiscordUser.UserID

	// TODO: Reimplement this later with admin role
	//if u := request.URL.Query().Get("userID"); u != "" {
	//	// TODO: Check admin
	//	userID = u
	//}

	data := make(map[string]interface{})

	// get webhook
	db := ws.Database
	if w, _ := db.FindWebhook(uid, userID); w != nil {
		data["Webhook"] = w
	}

	ws.MustExec("vieweditcreate", writer, request, data)
}
