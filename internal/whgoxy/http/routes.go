package http

import (
	"net/http"
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
