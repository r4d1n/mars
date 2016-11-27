package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=rover_user password=notamartian dbname=nasa_rover_data sslmode=disable")
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	n := new(Nasa)
	n.loadConfig("./config.json")
	n.crawl("curiosity")
}

func (n *Nasa) loadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
	err = json.Unmarshal(file, &n)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}
