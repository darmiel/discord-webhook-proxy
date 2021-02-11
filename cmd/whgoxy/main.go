package main

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbmongo"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbredis"
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
			log.Fatalln("(Database) Fatal:", err)
		}
	}()
	db.GlobalDatabase = database

	///// TODO: Remove this
	// page := cms.CMSPage{
	// 	URL: "/test",
	// 	Meta: cms.CMSPageMeta{
	// 		Title:         "Test-Page",
	// 		CreatorUserID: "710491120903127080",
	// 		CreatedAt:     time.Now(),
	// 	},
	// 	Updates: []cms.CMSPageUpdate{
	// 		{
	// 			UpdatedAt:     time.Now(),
	// 			UpdaterUserID: "431537207606640651",
	// 		},
	// 	},
	// 	Preferences: cms.CMSPagePreferences{
	// 		AuthorVisible:    true,
	// 		Dynamic:          false,
	// 		URLCaseSensitive: true,
	// 		UseMarkdown:      true,
	// 	},
	// 	Content: "Hallo, das ist ein Test!",
	// }
	// log.Println("Saving cmd page to database:", database.SaveCMSPage(page))
	//

	// load redis
	client := dbredis.NewClient(conf.Redis)
	if err := client.Set(context.TODO(), "heartbeat", 1, 0).Err(); err != nil {
		log.Fatalln("(Redis) Fatal:", err)
		return
	}
	dbredis.GlobalRedis = client

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
