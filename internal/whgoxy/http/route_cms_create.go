package http

import (
	"encoding/json"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type cmsPageCreatePayload struct {
	PageURL   string `json:"page_url"`
	PageTitle string `json:"page_title"`

	Payload string `json:"payload"`

	AuthorVisible    string `json:"author_visible"`
	UpdatersVisible  string `json:"updaters_visible"`
	URLCaseSensitive string `json:"url_case_sensitive"`
	UseMarkdown      string `json:"use_markdown"`
}

func (ws *WebServer) createCMSPageHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("cms_vieweditcreate", writer, request, nil)
}

func (ws *WebServer) createCMSPageBackendHandler(writer http.ResponseWriter, request *http.Request) {
	u, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	// check permissions
	if !u.DiscordUser.HasPermission(discord.PermissionCMSCreatePage) {
		writer.WriteHeader(403)
		_, _ = fmt.Fprintln(writer, "You don't have permissions to create a cms page")
		return
	}

	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprintln(writer, "Error reading your payload")
		return
	}

	var data cmsPageCreatePayload
	if err := json.Unmarshal(all, &data); err != nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprintln(writer, "Error decoding your payload")
		return
	}

	// find existing page
	db := ws.Database
	if oldPage, err := db.FindCMSPage(data.PageURL); oldPage != nil || err == nil {
		writer.WriteHeader(300)
		_, _ = fmt.Fprintln(writer, "A page with the same page url already exists")
		return
	}

	// check if url is from another (inbuilt) route
	for _, r := range ws.getRoutes() {
		if strings.ToLower(r.Path) == strings.ToLower(data.PageURL) {
			writer.WriteHeader(403)
			_, _ = fmt.Fprintf(writer, "Cannot overwrite inbuilt route (%s <-> %s)\n", r.Path, data.PageURL)
			return
		}
	}

	// create page
	page := cms.CMSPage{
		URL: data.PageURL,
		Meta: cms.CMSPageMeta{
			Title:         data.PageTitle,
			CreatorUserID: u.DiscordUser.UserID,
			CreatedAt:     time.Now(),
		},
		Updates: []cms.CMSPageUpdate{},
		Preferences: cms.CMSPagePreferences{
			AuthorVisible:    data.AuthorVisible == "on",
			UpdatersVisible:  data.UpdatersVisible == "on",
			Dynamic:          false,
			URLCaseSensitive: data.URLCaseSensitive == "on",
			UseMarkdown:      data.UseMarkdown == "on",
		},
		Content: data.Payload,
	}

	// save page
	if err := db.SaveCMSPage(page); err != nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprintf(writer, "Error saving page to database: %s", err.Error())
		return
	}

	writer.WriteHeader(200)
	_, _ = fmt.Fprintln(writer, "Success")
}
