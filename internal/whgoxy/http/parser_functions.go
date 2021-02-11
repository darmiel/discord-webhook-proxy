package http

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"html"
	"html/template"
	"log"
	"time"
)

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
	"Unescape": func(s string) string {
		return html.UnescapeString(s)
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
	"GetContent": func(cms *cms.CMSPage) string {
		return cms.GetContent()
	},

	// Permissions
	"HasPermissionCMSEditPage": func(u *discord.DiscordUser) bool {
		log.Println("User", u, "requested permission", discord.PermissionCMSEditPage)
		res := u.HasPermission(discord.PermissionCMSEditPage)
		log.Println("  -> Res:", res)
		return res
	},
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
