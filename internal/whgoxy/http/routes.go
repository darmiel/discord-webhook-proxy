package http

import "net/http"

// routes
type Route struct {
	Path string
	Func func(http.ResponseWriter, *http.Request)
}

func (ws *WebServer) getRoutes() []Route {
	return []Route{
		// "static"
		{"/", ws.homeRouteHandler},
		{"/examples", ws.exampleRouteHandler},

		// "api"
		{"/dashboard/create/req", ws.createWebhookRouteHandler},
		{"/dashboard/delete", ws.deleteWebhookRouteHandler},

		// "other"
		{"/call/json/{userid}/{uid}/{secret}", ws.jsonWebhookRouteHandler},

		// dashboard
		{"/dashboard", ws.dashboardHandler},
		{"/dashboard/create", ws.createRouteHandler},
		{"/dashboard/edit/{uid}", ws.editWebhookHandler},
	}
}
