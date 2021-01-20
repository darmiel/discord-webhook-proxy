package db

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
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

	FindWebhooks(userID string) (w []*discord.Webhook, err error)

	Disconnect() (err error)
}

type MongoDatabase struct {
	client   *mongo.Client
	context  context.Context
	database string
}

func (mdb *MongoDatabase) collection() (collection *mongo.Collection) {
	return mdb.client.Database(mdb.database).Collection(CollectionName)
}

// SaveWebhook ...
func (mdb *MongoDatabase) SaveWebhook(w *discord.Webhook) (err error) {
	filter := w.CreateFilter(false)
	update := bson.M{"$set": w}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.collection().UpdateOne(mdb.context, filter, update, updateOpts)
	return err
}

// FindWebhook ...
func (mdb *MongoDatabase) FindWebhook(uid string, userID string) (w *discord.Webhook, err error) {
	filter := (&discord.Webhook{UserID: userID, UID: uid}).CreateFilter(false)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

// FindWebhookWithSecret ...
func (mdb *MongoDatabase) FindWebhookWithSecret(uid string, userID string, secret string) (w *discord.Webhook, err error) {
	filter := (&discord.Webhook{UserID: userID, UID: uid, Secret: secret}).CreateFilter(true)
	w, err = mdb.findWebhookWithFilter(filter)
	return w, err
}

func (mdb *MongoDatabase) Disconnect() (err error) {
	if err := mdb.client.Disconnect(mdb.context); err != nil {
		log.Fatalln("Error while disconnecting:", err.Error())
	}
	return nil
}

func (mdb *MongoDatabase) FindWebhooks(userID string) (w []*discord.Webhook, err error) {
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

	return w, nil
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
