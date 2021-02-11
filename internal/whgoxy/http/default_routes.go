package http

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

func (ws *WebServer) homeRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("home", writer, request, nil)
}

func (ws *WebServer) createRouteHandler(writer http.ResponseWriter, request *http.Request) {
	// ws.MustExec("create", writer, request, nil)
	ws.MustExec("vieweditcreate", writer, request, map[string]interface{}{
		"ModeCreate": true,
	})
}

func (ws *WebServer) exampleRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("examples", writer, request, nil)
}

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

func (ws *WebServer) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	data := bson.M{}

	// find webhooks
	/* benchmark */
	start := time.Now()
	webhooks, err := ws.Database.FindWebhooks(user.DiscordUser.UserID)
	log.Println("🐛 Debug: Fetched webhooks in:", time.Since(start).Milliseconds())
	if err != nil {
		data["Error"] = err.Error()
	} else {
		data["Webhooks"] = webhooks
	}

	ws.MustExec("dashboard", writer, request, data)
}

func (ws *WebServer) editWebhookHandler(writer http.ResponseWriter, request *http.Request) {
	// get login
	user, die := auth.GetUserOrDie(request, writer)
	if die {
		return
	}

	vars := mux.Vars(request)
	uid := vars["uid"]

	userID := user.DiscordUser.UserID

	// TODO: Reimplement this later with admin role
	//if u := request.URL.Query().Get("userID"); u != "" {
	//	// TODO: Check admin
	//	userID = u
	//}

	data := make(map[string]interface{})

	// get webhook
	db := ws.Database
	if w, _ := db.FindWebhook(uid, userID); w != nil {
		data["Webhook"] = w
	}

	ws.MustExec("vieweditcreate", writer, request, data)
}
