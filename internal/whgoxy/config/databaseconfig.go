package config

var ConfigMongo MongoConfig

type MongoConfig struct {
	MongoDatabase         string `json:"mongo_database"`
	MongoHost             string `json:"mongo_host"`
	MongoAuthUser         string `json:"mongo_auth_user"`
	MongoAuthPass         string `json:"mongo_auth_pass"`
	MongoConnectionString string `json:"mongo_connection_string"`
}
