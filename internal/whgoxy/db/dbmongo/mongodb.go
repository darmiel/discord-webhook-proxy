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

const WebhookCollectionName string = "whgoxy-webhooks"
const UserCollectionName string = "whgoxy-users"
const CMSCollectionName string = "whgoxy-cms"

/// Mongo Functions

type mongoDatabase struct {
	client   *mongo.Client
	context  context.Context
	database string
}

// webhookCollection returns the webhookCollection used for whgoxy
func (mdb *mongoDatabase) webhookCollection() (collection *mongo.Collection) {
	return mdb.client.Database(mdb.database).Collection(WebhookCollectionName)
}

func (mdb *mongoDatabase) userCollection() (collection *mongo.Collection) {
	return mdb.client.Database(mdb.database).Collection(UserCollectionName)
}

func (mdb *mongoDatabase) cmsCollection() (collection *mongo.Collection) {
	return mdb.client.Database(mdb.database).Collection(CMSCollectionName)
}

// Disconnect disconnects the mongo client and should be run at the end of the program
func (mdb *mongoDatabase) Disconnect() (err error) {
	if err := mdb.client.Disconnect(mdb.context); err != nil {
		log.Fatalln("Error while disconnecting:", err.Error())
	}
	return nil
}

func connect(applyURI string, database string) (mdb db.Database, err error) {
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

	return &mongoDatabase{client, ctx, database}, nil
}

func NewDatabase(options config.MongoConfig) (db db.Database, err error) {
	log.Println("👉 Using mongo as database!")

	var uri string
	if options.MongoConnectionString != "" {
		uri = options.MongoConnectionString
	} else {
		uri = BuildApplyURI(options.MongoAuthUser, options.MongoAuthPass, options.MongoHost, options.MongoDatabase)
	}

	db, err = connect(uri, options.MongoDatabase)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BuildApplyURI(authUser string, authPass string, host string, database string) (res string) {
	res = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", authUser, authPass, host, database)
	return res
}
