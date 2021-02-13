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
		{"/dashboard/create", ws.createWebhookFrontendRouteHandler},
		{"/dashboard/edit/{uid}", ws.editWebhookHandler},

		// cms
		{"/cms/create", ws.createCMSPageHandler},
		{"/cms/create/req", ws.createCMSPageBackendHandler},
	}
}
