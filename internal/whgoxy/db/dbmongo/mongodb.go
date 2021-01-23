package dbmongo

import (
	"context"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const CollectionName string = "whgoxy"

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

func NewMongoDatabase(client *mongo.Client, context context.Context, database string) (db db.Database) {
	return &MongoDatabase{
		client:   client,
		context:  context,
		database: database,
	}
}

///

func NewDatabase(options config.MongoConfig) (db db.Database, err error) {
	log.Println("👉 Using mongo as database!")

	var uri string
	if options.MongoConnectionString != "" {
		uri = options.MongoConnectionString
	} else {
		uri = BuildApplyURI(options.MongoAuthUser, options.MongoAuthPass, options.MongoHost, options.MongoDatabase)
	}

	log.Println("🤫", uri)

	db, err = ConnectMongoDatabase(uri, options.MongoDatabase)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectMongoDatabase(applyURI string, database string) (mdb db.Database, err error) {
	uri := options.Client().ApplyURI(applyURI)
	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return NewMongoDatabase(client, ctx, database), nil
}

func BuildApplyURI(authUser string, authPass string, host string, database string) (res string) {
	res = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", authUser, authPass, host, database)
	return res
}
