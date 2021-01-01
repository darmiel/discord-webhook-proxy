package whgoxy

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)

var (
	db *sql.DB
)

func InitDatabase(file string) (err error) {
	database, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		log.Fatalln("Error opening database:", err.Error())
		return err
	}
	db = database

	// execute setup
	data, err := ioutil.ReadFile("./setup.sql")
	if err != nil {
		log.Fatalln("Error accessing setup.sql:", err.Error())
		return err
	}

	// no data
	if len(data) == 0 {
		log.Println("[WARN] No data in setup.sql")
	} else {
		statement, err := db.Prepare(string(data))
		if err != nil {
			log.Fatalln("Error preparing setup statement:", err.Error())
			return err
		}

		if _, err := statement.Exec(); err != nil {
			log.Fatalln("Error executing setup statement:", err.Error())
			return err
		}
	}

	return nil
}
