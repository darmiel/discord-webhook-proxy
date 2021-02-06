package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"log"
)

func load(fpath string, v interface{}) {
	if _, err := toml.DecodeFile(fpath, v); err != nil {
		panic(err)
	}

	if marshal, err := json.Marshal(v); err == nil {
		log.Println(string(marshal))
	}
}

type Config struct {
	Web   WebConfig
	Auth  OAuthConfig
	Mongo MongoConfig
	Redis RedisConfig
}

type MongoConfig struct {
	MongoDatabase         string `json:"mongo_database"`
	MongoHost             string `json:"mongo_host"`
	MongoAuthUser         string `json:"mongo_auth_user"`
	MongoAuthPass         string `json:"mongo_auth_pass"`
	MongoConnectionString string `json:"mongo_connection_string"`
}

type OAuthConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string

	// Endpoint
	AuthURL  string
	TokenURL string

	// Cookie
	CookieHost   string
	CookieSecret string
	CookieName   string
}

type WebConfig struct {
	WebDir string
	Addr   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Load() (conf Config) {
	// oauth config
	load("config.toml", &conf)
	return
}
