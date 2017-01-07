package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", os.Getenv("NASA_DB_USER"), os.Getenv("NASA_DB_PASS"), os.Getenv("NASA_DB_NAME"), os.Getenv("NASA_DB_HOST"), os.Getenv("NASA_DB_SSL")))
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup
	rovers := []string{"curiosity", "opportunity", "spirit"}
	worklist := make(chan []string)
	done := make(map[string]bool)
	// load rover names into worklist channel
	go func() { worklist <- rovers }()
	s := Scraper{APIKey: os.Getenv("NASA_API_KEY"), AWSRegion: os.Getenv("NASA_AWS_REGION"), S3Bucket: os.Getenv("NASA_S3_BUCKET")}
	for list := range worklist {
		for _, name := range list {
			if !done[name] {
				done[name] = true
				wg.Add(1)
				go func(nm string) {
					defer wg.Done()
					err := s.crawl(nm)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Completed rover: %s \n", nm)
					return
				}(name)
			}
		}
		wg.Wait()
		os.Exit(0)
	}
}
