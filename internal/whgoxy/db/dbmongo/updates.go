package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

func (mdb *mongoDatabase) SaveDiscordUser(dcu *discord.DiscordUser) (err error) {
	log.Println("💾 Saving user", dcu.GetFullName(), "...")

	filter := bson.M{"user_id": dcu.UserID}
	update := bson.M{"$set": dcu}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.userCollection().UpdateOne(mdb.context, filter, update, updateOpts)
	return
}
