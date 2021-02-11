package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const rootTmpl = `{{ define "root" }} {{ template "base" . }} {{ end }}`

type TemplateParser struct {
}

func NewTemplateParser() (parser *TemplateParser) {
	return &TemplateParser{}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

var funcs = map[string]interface{}{
	"Avatar": func(u *discord.DiscordUser) string {
		return u.GetAvatarUrl()
	},
	"FullName": func(u *discord.DiscordUser) string {
		return u.GetFullName()
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
	"StrAgo": func(sec int64) string {
		if sec == 0 {
			return "/"
		}
		return fmtDuration(time.Since(time.Unix(sec, 0)))
	},
	"GetUserByID": func(userID string) *discord.DiscordUser {
		if userID == "0" {
			return &discord.DiscordUser{
				UserID:        "0",
				Username:      "whgoxy-System",
				Discriminator: "0000",
			}
		}

		user, err := db.GlobalDatabase.FindDiscordUser(userID)
		if err != nil {
			return nil
		}
		return user
	},
	"CMSGetUpdateInfo": func(cms *cms.CMSPage) *cms.CMSUpdateInfo {
		update := cms.GetLastUpdate()
		if update == nil {
			return nil
		}
		return update.GetInfo()
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
