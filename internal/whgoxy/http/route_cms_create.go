package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"net/http"
)

func (ws *WebServer) createCMSPageHandler(writer http.ResponseWriter, request *http.Request) {
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}
	// check permissions
	if !user.DiscordUser.HasPermission(discord.PermissionCMSCreatePage) {
		writer.WriteHeader(403)
		_, _ = fmt.Fprintln(writer, "You don't have permissions to create a cms page!")
		return
	}

	ws.MustExec("cms_vieweditcreate", writer, request, nil)
}
