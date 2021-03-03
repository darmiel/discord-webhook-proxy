package http

import (
	"encoding/base64"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"github.com/gorilla/mux"
	"net/http"
)

func (ws *WebServer) editCMSPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	b64 := vars["full_url"]

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

	ws.MustExec("cms_vieweditcreate", writer, request, map[string]interface{}{
		"CMS":   page,
		"Error": errorMessage,
	})
}
