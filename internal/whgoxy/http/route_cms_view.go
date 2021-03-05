package http

import (
	"log"
	"net/http"
	"strings"
)

func (ws *WebServer) error404(writer http.ResponseWriter, request *http.Request) {

	// requested path
	reqPage := request.URL.Path

	log.Println("Searching for CMS Sites, because user landed on 404")

	// Check if page is cms
	pages, err := ws.Database.FindAllCMSPages()
	log.Println("pages, err :=", pages, err)

	if pages != nil && err == nil {
		for _, cms := range pages {
			log.Println("[CMS] Checking", cms.URL, "(", cms.Meta.Title, ")", "with", reqPage)

			var matches = false
			if !cms.Preferences.URLCaseSensitive {
				matches = strings.ToLower(cms.URL) == strings.ToLower(reqPage)
			} else {
				matches = cms.URL == reqPage
			}

			if !matches {
				continue
			}

			log.Println("  Executing CMS", cms.URL, "(", cms.Meta.Title, ")")

			writer.WriteHeader(200)
			ws.MustExec("cms_template", writer, request, map[string]interface{}{
				"CMS": cms,
			})

			return
		}
	}

	ws.MustExec("err_404", writer, request, nil)
}
