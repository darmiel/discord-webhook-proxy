package main

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load config
	config.Load()

	// load database
	database, err := db.NewDatabase(config.ConfigMongo)
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
	webConfig := config.ConfigWeb
	ws := http.NewWebServer(webConfig, database)

	// auth
	auth.InitOAuth2(config.ConfigOAuth, ws.Router)

	if err := ws.Run(); err != nil {
		panic(err)
	}

	log.Println("✅  whgoxy is now (hopefully) running on " + webConfig.Addr)
	log.Println("Press CTRL-C to exit gracefully.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
