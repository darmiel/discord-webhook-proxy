package http

import (
	"net/http"
)

func (ws *WebServer) createCMSPageHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("cms_vieweditcreate", writer, request, nil)
}
