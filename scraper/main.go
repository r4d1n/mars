package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
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
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	rovers := []string{"curiosity", "opportunity", "spirit"}
	worklist := make(chan []string)
	done := make(map[string]bool)
	// load rover names into worklist channel
	go func() { worklist <- rovers }()
	s := Scraper{APIKey: c.APIKey, AWSRegion: c.AWSRegion, S3Bucket: c.S3Bucket}
	for list := range worklist {
		for _, name := range list {
			if !done[name] {
				done[name] = true
				go func(nm string) {
					err := s.crawl(nm)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Completed rover: %s \n", nm)
				}(name)
			}
		}
	}
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
