package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const rootTmpl = `{{ define "root" }} {{ template "base" . }} {{ end }}`

type TemplateParser struct {
}

func NewTemplateParser() (parser *TemplateParser) {
	return &TemplateParser{}
}

var funcs = map[string]interface{}{
	"Avatar": func(u *discord.DiscordUser) string {
		return u.GetAvatarUrl()
	},
	"WebhookCount": func(u *discord.DiscordUser) uint {
		count, err := db.GlobalDatabase.CountWebhooksForUser(u.UserID)
		if err != nil {
			log.Println("🐛 Error counting webhooks:", err)
			return 0
		}
		return count
	},
	"GetStats": func(w *discord.Webhook) *discord.WebhookStats {
		stats := w.GetStats()
		return stats
	},
	"Escape": func(s string) string {
		return template.HTMLEscaper(s)
	},
}

func (parser *TemplateParser) ParseTemplate(name string) (tpl *template.Template, err error) {
	root, err := template.New("root").Funcs(funcs).Parse(rootTmpl)
	if err != nil {
		return nil, err
	}

	tmplDir := fmt.Sprintf("%s/template", "web")
	componentsDir := tmplDir + "/components"

	basePath := fmt.Sprintf("%s/base.gohtml", tmplDir)
	tmplPath := fmt.Sprintf("%s/%s.gohtml", tmplDir, name)

	var files []string

	// read all components from "components" dir
	if err := filepath.Walk(componentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// add component to file if it has the ".gohtml" suffix
		if strings.HasSuffix(path, ".gohtml") {
			log.Println("[+]", name, "Added component", path)
			files = append(files, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// add base template and selected template
	files = append(files, basePath, tmplPath)

	// return parsed template
	return root.ParseFiles(files...)
}

// MustParseTemplate calls ParseTemplate(...) and panics on an error.
func (parser *TemplateParser) MustParseTemplate(name string) *template.Template {
	tmpl, err := parser.ParseTemplate(name)
	if err != nil {
		panic(err)
	}

	return tmpl
}
