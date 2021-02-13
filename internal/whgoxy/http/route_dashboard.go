package http

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func (ws *WebServer) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	data := bson.M{}

	// find webhooks
	/* benchmark */
	start := time.Now()
	webhooks, err := ws.Database.FindWebhooks(user.DiscordUser.UserID)
	log.Println("🐛 Debug: Fetched webhooks in:", time.Since(start).Milliseconds())
	if err != nil {
		data["Error"] = err.Error()
	} else {
		data["Webhooks"] = webhooks
	}

	ws.MustExec("dashboard", writer, request, data)
}
