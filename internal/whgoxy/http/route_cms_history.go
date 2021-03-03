package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"github.com/gorilla/mux"
	"github.com/sergi/go-diff/diffmatchpatch"
	"net/http"
	"strconv"
)

func (ws *WebServer) historyeditCMSPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	b64 := vars["full_url"]

	/// check perm
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}
	if !user.DiscordUser.HasPermission(discord.PermissionCMSViewHistory) {
		_, _ = fmt.Fprintln(writer, "You don't have permissions to view the history of a cms page")
		return
	}
	///

	page := &cms.CMSPage{
		URL:         "/not/found",
		Meta:        cms.CMSPageMeta{},
		Updates:     nil,
		Preferences: cms.CMSPagePreferences{},
		Content:     "< page not found >",
	}

	errorMessage := ""
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		errorMessage = err.Error()
	} else {
		fullUrl := string(bytes)
		cmsPage, err := ws.Database.FindCMSPage(fullUrl)
		if err != nil {
			errorMessage = err.Error()
		} else {
			page = cmsPage
		}
	}

	ws.MustExec("cms_history", writer, request, map[string]interface{}{
		"CMS":   page,
		"Error": errorMessage,
	})
}

func (ws *WebServer) historyGetCMSPageHandler(writer http.ResponseWriter, request *http.Request) {

	/// check perm
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}
	if !user.DiscordUser.HasPermission(discord.PermissionCMSViewHistory) {
		_, _ = fmt.Fprintln(writer, "You don't have permissions to view the history of a cms page")
		return
	}
	///

	he := func(err error) {
		writer.WriteHeader(500)
		_, _ = fmt.Fprintln(writer, `ERROR: `+err.Error())
	}

	vars := mux.Vars(request)
	b64 := vars["full_url"]
	indexStr := vars["index"]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		he(err)
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		he(err)
		return
	}
	page, err := ws.Database.FindCMSPage(string(bytes))
	if err != nil {
		he(err)
		return
	}
	if page == nil {
		he(errors.New("page was nil"))
		return
	}

	for i, update := range page.Updates {
		if index != i {
			continue
		}

		// make diff
		dmp := diffmatchpatch.New()
		diff := dmp.DiffMain(update.Previous, page.Content, true)

		writer.WriteHeader(200)
		_, _ = fmt.Fprintln(writer, dmp.DiffPrettyHtml(diff))

		return
	}

	writer.WriteHeader(404)
	_, _ = fmt.Fprintln(writer, "update "+indexStr+" not found")
}
