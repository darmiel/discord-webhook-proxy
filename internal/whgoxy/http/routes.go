package http

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (ws *WebServer) homeRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("home", writer, request, nil)
}

func (ws *WebServer) createRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("create", writer, request, nil)
}

func (ws *WebServer) exampleRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("examples", writer, request, nil)
}

func (ws *WebServer) error404(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("err_404", writer, request, nil)
}

func (ws *WebServer) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	data := bson.M{}

	// find webhooks
	webhooks, err := ws.Database.FindWebhooks(user.DiscordUser.UserID)
	if err != nil {
		data["Error"] = err.Error()
	} else {
		data["Webhooks"] = webhooks
	}

	ws.MustExec("dashboard", writer, request, data)
}
