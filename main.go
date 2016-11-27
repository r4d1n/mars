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
	nc := new(NasaCrawler)
	nc.loadConfig("./config.json")
	pm := nc.getManifest("curiosity")
	nc.getPhotoData(pm)
}

func (nc *NasaCrawler) loadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
	err = json.Unmarshal(file, &nc)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}
