package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"net/http"
)

func (ws *WebServer) adminDashboardPageHandler(writer http.ResponseWriter, request *http.Request) {
	// check perm
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}
	// check permissions
	if !user.DiscordUser.HasPermission(discord.PermissionAdminDashboardView) {
		_, _ = fmt.Fprintln(writer, "You don't have permission to acces the admin dashboard.")
		return
	}

	errorMessage := ""

	// load pages
	pages, err := ws.Database.FindAllCMSPages()
	if err != nil {
		errorMessage = err.Error()
	}

	if pages == nil {
		pages = []*cms.CMSPage{}
	}

	ws.MustExec("admin", writer, request, map[string]interface{}{
		"Error": errorMessage,
		"Pages": pages,
	})
}
