package http

import "net/http"

func (ws *WebServer) homeRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("home", writer, request, nil)
}
