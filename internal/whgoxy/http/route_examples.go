package http

import "net/http"

func (ws *WebServer) exampleRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("examples", writer, request, nil)
}
