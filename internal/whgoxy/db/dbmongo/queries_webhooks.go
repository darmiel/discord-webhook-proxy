package dbmongo

import (
	"errors"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errInvalidSecret = errors.New("invalid secret")
)

func (mdb *mongoDatabase) findWebhookWithFilter(filter bson.M) (w *discord.Webhook, err error) {
	res := mdb.webhookCollection().FindOne(mdb.context, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	w = &discord.Webhook{}
	if err = res.Decode(w); err != nil {
		return nil, err
	}
	w.ParseNewLine()

	return w, nil
}

// FindWebhook ...
func (mdb *mongoDatabase) FindWebhook(uid string, userID string) (w *discord.Webhook, err error) {
	// check cache
	if w, found := db.WebhookCache.Get(db.GetCacheKeyManual(userID, uid)); found {
		return w.(*discord.Webhook), nil
	}

	filter := (&discord.Webhook{UserID: userID, UID: uid}).CreateFilter(false)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

// FindWebhookWithSecret ...
func (mdb *mongoDatabase) FindWebhookWithSecret(uid string, userID string, secret string) (w *discord.Webhook, err error) {
	// check cache
	if w, found := db.WebhookCache.Get(db.GetCacheKey(w)); found {
		webhook := w.(*discord.Webhook)
		if webhook.Secret == secret {
			return webhook, nil
		} else {
			return nil, errInvalidSecret
		}
	}

	filter := (&discord.Webhook{UserID: userID, UID: uid, Secret: secret}).CreateFilter(true)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

// FindWebhooks ...
func (mdb *mongoDatabase) FindWebhooks(userID string) (w []*discord.Webhook, err error) {
	// check cache
	if w, found := db.UserWebhookCache.Get(userID); found {
		return w.([]*discord.Webhook), nil
	}

	filter := bson.M{
		"user_id": userID,
	}

	res, err := mdb.webhookCollection().Find(mdb.context, filter, options.Find())
	if err != nil {
		return nil, err
	}

	for res.Next(mdb.context) {
		var webhook *discord.Webhook
		if err := res.Decode(&webhook); err != nil {
			return nil, err
		}

		w = append(w, webhook)
	}

	// update cache
	db.UserWebhookCache.Set(userID, w, cache.DefaultExpiration)

	return w, nil
}


func (mdb *mongoDatabase) CountWebhooksForUser(userID string) (count uint, err error) {
	webhooks, err := mdb.FindWebhooks(userID)
	if err != nil {
		return 0, err
	}
	return uint(len(webhooks)), nil
}