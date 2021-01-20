package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WebServer struct {
	Parser    *TemplateParser
	Router    *mux.Router
	templates map[string]*template.Template
	Conf      config.WebConfig
	Database  db.Database
}

func NewWebServer(conf config.WebConfig, db db.Database) (ws *WebServer) {
	router := mux.NewRouter()
	parser := NewTemplateParser()

	return &WebServer{
		Router:    router,
		Parser:    parser,
		templates: make(map[string]*template.Template),
		Conf:      conf,
		Database:  db,
	}
}

// updateRoutes adds a static route for the /static/ folder,
// a 404 not found route and all other "own" routes
func (ws *WebServer) updateRoutes() {
	router := ws.Router

	// static dir
	staticDir := ws.Conf.WebDir + "/static/"
	prefix := http.StripPrefix("/static", http.FileServer(http.Dir(staticDir)))
	router.PathPrefix("/static/").Handler(prefix)

	// routes
	type Route struct {
		Path string
		Func func(http.ResponseWriter, *http.Request)
	}

	routes := []Route{
		// "static"
		{"/", ws.homeRouteHandler},
		{"/examples", ws.exampleRouteHandler},

		// "api"
		{"/dashboard/create/req", ws.createWebhookRouteHandler},
		// {"/call/{user_id}/{uid}/{secret}", ws.safeWebhookRouteHandler},

		// "other"
		{"/call/json/{userid}/{uid}/{secret}", ws.jsonWebhookRouteHandler},

		// dashboard
		{"/dashboard", ws.dashboardHandler},
		{"/dashboard/create", ws.createRouteHandler},
	}

	for _, r := range routes {
		log.Println("[Router] Registered route ", r.Path)
		router.HandleFunc(r.Path, r.Func)
	}

	// 404
	router.NotFoundHandler = http.HandlerFunc(ws.error404)
}

// readTemplates reads all *.gohtml files from the web root and parses them
func (ws *WebServer) readTemplates() {
	// ws.templates["home"] = ws.Parser.MustParseTemplate("home")
	// ws.templates["err_404"] = ws.Parser.MustParseTemplate("err_404")

	templateDir := ws.Conf.WebDir + "/template/"
	dir, err := ioutil.ReadDir(templateDir)
	log.Println("Reading dir:", templateDir, "with result:", dir)
	if err != nil {
		panic(err)
	}

	// clear old templates
	ws.templates = make(map[string]*template.Template)

	// find all templates in folder
	for _, file := range dir {
		name := file.Name()
		log.Println("File/Dir:", name)

		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(name, ".gohtml") {
			continue
		}
		// remove extension
		name = name[:len(name)-7]

		log.Println("[Web] Parsing template:", name)
		ws.templates[name] = ws.Parser.MustParseTemplate(name)
	}

	log.Println("[Web] 👉 Parsed", len(ws.templates), "templates")
}

func (ws *WebServer) Run() (err error) {
	ws.readTemplates()
	ws.updateRoutes()

	// http
	return http.ListenAndServe(ws.Conf.Addr, ws.Router)
}

func (ws *WebServer) Exec(name string, r *http.Request, w http.ResponseWriter, data map[string]interface{}) (err error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	//// Global data
	// CurrentURL
	data["CurrentURL"] = r.URL.String()

	// User
	if user, ok := auth.GetUser(r); ok && user != nil {
		data["User"] = user.DiscordUser
		log.Println("OK user found:", user, ok)
	} else {
		log.Println("ERR user not found:", user, ok)
	}

	// get template
	tpl, ok := ws.templates[name]
	if !ok {
		w.WriteHeader(404)
		_, _ = fmt.Fprint(w, "Template "+name+" not found.")
		return
	}
	return tpl.Execute(w, data)
}

func (ws *WebServer) MustExec(name string, w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	log.Println(r.RemoteAddr, "requested uri:", r.RequestURI)

	if err := ws.Exec(name, r, w, data); err != nil {
		log.Println("[WARNING] Error occurred on rendering template:", err)
	}
}
