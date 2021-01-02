package whgoxy

import (
	"fmt"
	"github.com/darmiel/whgoxy/db"
	"github.com/darmiel/whgoxy/discord"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var Opt *Options
var WHDatabase db.Database

func New(options *Options) {
	Opt = options

	if options.MongoUse {
		log.Println("👉 Using mongo as database!")

		uri := db.BuildApplyURI(options.MongoAuthUser, options.MongoAuthPass, options.MongoHost, options.MongoDatabase)
		database, err := db.ConnectMongoDatabase(uri, options.MongoDatabase)
		if err != nil {
			log.Fatalln("Fatal: ", err)
			return
		}
		log.Println("✅ All done!")

		webhook := discord.NewWebhook("https://ptb.discord.com/api/webhooks/794739759909830666/6wjgrcP8emYcE4-awAd8G83aiF8gkaOEF84kiSOfg8Giaife4kF", &discord.WebhookData{"content": "jo"})
		if err := webhook.CheckValidity(true); err != nil {
			fmt.Println("Warn: Created webhook is not valid:", err)
		}

		log.Println("Okidoki")

		WHDatabase = database
	} else {
		// No database selected
		log.Fatalln("Error: No database selected.")
		return
	}

	log.Println("whgoxy is now running. Press CTRL-C to exit gracefully.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	defer func() {
		log.Println("[Shutdown] Closing database connection")
		if err := WHDatabase.Disconnect(); err != nil {
			log.Fatalln("Fatal:", err)
		}
	}()
}
