package db

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"log"
)

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
