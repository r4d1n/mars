package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	APIKey    string
	DBName    string
	DBUser    string
	DBPass    string
	AWSRegion string
	S3Bucket  string
}

var db *sql.DB
var c Config

func init() {
	c.load("./config.json")
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", c.DBUser, c.DBPass, c.DBName))
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	s := Scraper{APIKey: c.APIKey, AWSRegion: c.AWSRegion, S3Bucket: c.S3Bucket}
	s.crawl("curiosity")
}

func (c *Config) load(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Missing config file: ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Could not parse config file: ", err)
	}
}
