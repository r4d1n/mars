package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	data "github.com/r4d1n/nasa-mars-photos/roverdata"
)

type config struct {
	APIKey    string
	DBName    string
	DBUser    string
	DBPass    string
	AWSRegion string
	S3Bucket  string
	SSLMode   string
	DBHost    string
}

var db *sql.DB
var c config

func init() {
	c.load("./config.json")
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", c.DBUser, c.DBPass, c.DBName, c.DBHost, c.SSLMode))
	data.DB = db
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	rovers := []string{"curiosity", "opportunity", "spirit"}
	s := Scraper{APIKey: c.APIKey, AWSRegion: c.AWSRegion, S3Bucket: c.S3Bucket}
	for _, name := range rovers {
		err := s.crawl(name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Completed rover: %s", name)
	}
	fmt.Printf("All rovers complete!")
}

func (c *config) load(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("missing config file: ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("could not parse config file: ", err)
	}
}
