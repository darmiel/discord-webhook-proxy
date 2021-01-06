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

func Load() {
	// oauth config
	load("auth_config.toml", &ConfigOAuth)

	// web config
	load("web_config.toml", &ConfigWeb)

	// db config
	load("db_config.toml", &ConfigMongo)
}
