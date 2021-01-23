package db

import (
	"context"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Cache
var webhookCache *cache.Cache
var userWebhookCache *cache.Cache

func init() {
	webhookCache = cache.New(5*time.Minute, 10*time.Minute)
	userWebhookCache = cache.New(5*time.Minute, 10*time.Minute)
}
func getCacheKeyManual(userID string, uid string) string {
	return userID + ":" + uid
}
func getCacheKey(w *discord.Webhook) string {
	return getCacheKeyManual(w.UserID, w.UID)
}

//

func NewDatabase(options config.MongoConfig) (db Database, err error) {
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

func ConnectMongoDatabase(applyURI string, database string) (mdb Database, err error) {
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
