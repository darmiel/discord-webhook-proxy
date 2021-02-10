package db

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"time"
)

// Cache
var WebhookCache *cache.Cache
var UserWebhookCache *cache.Cache
var CMSCache *cache.Cache

func init() {
	WebhookCache = cache.New(5*time.Minute, 10*time.Minute)
	UserWebhookCache = cache.New(5*time.Minute, 10*time.Minute)
	CMSCache = cache.New(60*time.Minute, 60*time.Minute)
}

func GetCacheKeyManual(userID string, uid string) string {
	return userID + ":" + uid
}
func GetCacheKey(w *discord.Webhook) string {
	return GetCacheKeyManual(w.UserID, w.UID)
}
