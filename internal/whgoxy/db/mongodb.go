package db

import (
	"context"
	"errors"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const CollectionName string = "whgoxy"

type Database interface {
	// SaveWebhook inserts the specified webhook into the database or updates it if the _id is already present
	// returns an error if anything went wrong.
	SaveWebhook(w *discord.Webhook) (err error)

	// FindWebhook searches for a webhook by the given id (uid)
	// returns the webhook if found, otherwise an error if anything went wrong.
	FindWebhook(uid string, userID string) (w *discord.Webhook, err error)

	// FindWebhook searches for a webhook by the given id (uid) AND the matching secret
	// returns the webhook if found, otherwise an error if anything went wrong.
	FindWebhookWithSecret(uid string, userID string, secret string) (w *discord.Webhook, err error)

	// FindWebhooks returns all webhooks created by the user with the ID {userID}
	FindWebhooks(userID string) (w []*discord.Webhook, err error)

	// DeleteWebhook deletes the specified webhook if it was found.
	// Does not check further if the webhook existed before!
	DeleteWebhook(uid string, userID string) (err error)

	Disconnect() (err error)
}

/// Mongo Functions

type MongoDatabase struct {
	client   *mongo.Client
	context  context.Context
	database string
}

func (mdb *MongoDatabase) collection() (collection *mongo.Collection) {
	return mdb.client.Database(mdb.database).Collection(CollectionName)
}

func (mdb *MongoDatabase) Disconnect() (err error) {
	if err := mdb.client.Disconnect(mdb.context); err != nil {
		log.Fatalln("Error while disconnecting:", err.Error())
	}
	return nil
}

func NewMongoDatabase(client *mongo.Client, context context.Context, database string) (db Database) {
	return &MongoDatabase{
		client:   client,
		context:  context,
		database: database,
	}
}

func (mdb *MongoDatabase) findWebhookWithFilter(filter bson.M) (w *discord.Webhook, err error) {
	res := mdb.collection().FindOne(mdb.context, filter)
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

///

/// Webhook Query Functions

// FindWebhook ...
func (mdb *MongoDatabase) FindWebhook(uid string, userID string) (w *discord.Webhook, err error) {
	// check cache
	if w, found := webhookCache.Get(getCacheKey(w)); found {
		return w.(*discord.Webhook), nil
	}

	filter := (&discord.Webhook{UserID: userID, UID: uid}).CreateFilter(false)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

// FindWebhookWithSecret ...
func (mdb *MongoDatabase) FindWebhookWithSecret(uid string, userID string, secret string) (w *discord.Webhook, err error) {
	// check cache
	if w, found := webhookCache.Get(getCacheKey(w)); found {
		webhook := w.(*discord.Webhook)
		if webhook.Secret == secret {
			return webhook, nil
		} else {
			return nil, errors.New("invalid secret")
		}
	}

	filter := (&discord.Webhook{UserID: userID, UID: uid, Secret: secret}).CreateFilter(true)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

// FindWebhooks ...
func (mdb *MongoDatabase) FindWebhooks(userID string) (w []*discord.Webhook, err error) {
	// check cache
	if w, found := userWebhookCache.Get(userID); found {
		return w.([]*discord.Webhook), nil
	}

	filter := bson.M{
		"user_id": userID,
	}

	res, err := mdb.collection().Find(mdb.context, filter, options.Find())
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
	userWebhookCache.Set(userID, w, cache.DefaultExpiration)

	return w, nil
}

///

/// Webhook Update Functions

// SaveWebhook ...
func (mdb *MongoDatabase) SaveWebhook(w *discord.Webhook) (err error) {
	filter := w.CreateFilter(false)
	update := bson.M{"$set": w}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.collection().UpdateOne(mdb.context, filter, update, updateOpts)

	// save to cache
	if err == nil {
		webhookCache.Set(w.UserID+":"+w.UID, w, cache.DefaultExpiration)
		// invalidate user cache
		userWebhookCache.Delete(w.UserID)
	}

	return err
}

// DeleteWebhook ...
func (mdb *MongoDatabase) DeleteWebhook(uid string, userID string) (err error) {
	// invalidate user cache
	userWebhookCache.Delete(userID)

	filter := bson.M{
		"user_id": userID,
		"uid":     uid,
	}

	_, err = mdb.collection().DeleteOne(mdb.context, filter)

	// delete from cache
	if err == nil {
		webhookCache.Delete(getCacheKeyManual(userID, uid))
	}

	return
}

///
