package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/// Webhook Update Functions

// SaveWebhook ...
func (mdb *mongoDatabase) SaveWebhook(w *discord.Webhook) (err error) {
	filter := w.CreateFilter(false)
	update := bson.M{"$set": w}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.webhookCollection().UpdateOne(mdb.context, filter, update, updateOpts)

	// save to cache
	if err == nil {
		db.WebhookCache.Set(w.UserID+":"+w.UID, w, cache.DefaultExpiration)
		// invalidate user cache
		db.UserWebhookCache.Delete(w.UserID)
	}

	return
}

// DeleteWebhook ...
func (mdb *mongoDatabase) DeleteWebhook(uid string, userID string) (err error) {
	// invalidate user cache
	db.UserWebhookCache.Delete(userID)

	filter := bson.M{
		"user_id": userID,
		"uid":     uid,
	}

	_, err = mdb.webhookCollection().DeleteOne(mdb.context, filter)

	// delete from cache
	if err == nil {
		db.WebhookCache.Delete(db.GetCacheKeyManual(userID, uid))
	}

	return
}
