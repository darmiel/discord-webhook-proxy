package db

import (
	"context"
	"github.com/darmiel/whgoxy/discord"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const CollectionName string = "whgoxy"

type Database interface {
	// SaveWebhook inserts the specified webhook into the database or updates it if the _id is already present
	// returns an error if anything went wrong.
	SaveWebhook(w *discord.SavedWebhook) (err error)

	// FindWebhook searches for a webhook by the given id (object id)
	// returns the webhook if found, otherwise an error if anything went wrong.
	FindWebhook(uuid string) (w *discord.SavedWebhook, err error)

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
func (mdb *MongoDatabase) SaveWebhook(w *discord.SavedWebhook) (err error) {
	filter := bson.M{"uuid": w.UUID}
	update := bson.M{"$set": w}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.collection().UpdateOne(mdb.context, filter, update, updateOpts)
	return err
}

// FindWebhook ...
func (mdb *MongoDatabase) FindWebhook(uuid string) (w *discord.SavedWebhook, err error) {
	filter := bson.M{
		"uuid": uuid,
	}

	res := mdb.collection().FindOne(mdb.context, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	w = &discord.SavedWebhook{}
	if err = res.Decode(w); err != nil {
		return nil, err
	}

	return w, nil
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
