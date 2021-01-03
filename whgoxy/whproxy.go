package whgoxy

import (
	"github.com/darmiel/whgoxy/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func New(options *Options) {
	var database db.Database
	var err error

	if options.MongoUse {
		log.Println("👉 Using mongo as database!")

		uri := db.BuildApplyURI(options.MongoAuthUser, options.MongoAuthPass, options.MongoHost, options.MongoDatabase)
		database, err = db.ConnectMongoDatabase(uri, options.MongoDatabase)
		if err != nil {
			log.Fatalln("Fatal: ", err)
			return
		}
		log.Println("✅  All done!")
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
		if err := database.Disconnect(); err != nil {
			log.Fatalln("Fatal:", err)
		}
	}()
}
