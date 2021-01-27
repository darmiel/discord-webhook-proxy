package main

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbmongo"
	"github.com/darmiel/whgoxy/internal/whgoxy/http"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load config
	conf := config.Load()

	// load database
	database, err := dbmongo.NewDatabase(conf.Mongo)
	if err != nil {
		log.Fatalln("❌ Error connecting to database:", database)
		return
	}
	defer func() {
		log.Println("[Database] Closing database connection")
		if err := database.Disconnect(); err != nil {
			log.Fatalln("Fatal:", err)
		}
	}()

	// create web server
	ws := http.NewWebServer(conf.Web, database)

	// auth
	auth.InitOAuth2(conf.Auth, ws.Router, database)

	if err := ws.Run(); err != nil {
		panic(err)
	}

	log.Println("✅  whgoxy is now (hopefully) running on " + conf.Web.Addr)
	log.Println("Press CTRL-C to exit gracefully.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
