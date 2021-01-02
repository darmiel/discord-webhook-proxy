package whgoxy

import "flag"

type Options struct {
	DatabaseFile string `json:"database_file"`

	// MongoDB
	MongoUse      bool   `json:"use_mongo"`
	MongoDatabase string `json:"mongo_database"`
	MongoHost     string `json:"mongo_host"`
	MongoAuthUser string `json:"mongo_auth_user"`
	MongoAuthPass string `json:"mongo_auth_pass"`
}

func Parse() (opt *Options, err error) {
	opt = &Options{}

	// flags here \/
	flag.StringVar(&opt.DatabaseFile, "f", "./data.sqlite3", "File for database")

	// Mongo
	flag.BoolVar(&opt.MongoUse, "mongo", true, "Use MongoDB?")
	flag.StringVar(&opt.MongoDatabase, "mongo-db", "whgoxy", "MongoDB Database")
	flag.StringVar(&opt.MongoHost, "mongo-host", "whgoxy.x10rd.mongodb.net", "MongoDB Host")
	flag.StringVar(&opt.MongoAuthUser, "mongo-user", "whgoxy", "MongoDB Auth User")
	flag.StringVar(&opt.MongoAuthPass, "mongo-pass", "", "MongoDB Auth Password")
	// flags here /\

	flag.Parse()
	return opt, nil
}
