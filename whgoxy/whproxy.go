package whgoxy

import (
	"github.com/darmiel/whgoxy/db"
	"github.com/darmiel/whgoxy/router"
	"log"
	"net/http"
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

	log.Println("✅  whgoxy is now (hopefully) running on " + options.ApiBind)
	log.Println("Press CTRL-C to exit gracefully.")

	// Start http
	r := router.New(database)
	if err := http.ListenAndServe(options.ApiBind, r); err != nil {
		log.Fatalln("Error listening and serving:", err.Error())
		return
	}

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
