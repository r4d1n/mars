package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Scraper struct {
	APIKey    string
	AWSRegion string
	S3Bucket  string
}

func (s Scraper) crawl(name string) error {
	murl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", name, "?api_key=", s.APIKey)
	res, err := http.Get(murl)
	if err != nil {
		log.Fatal(err)
	} else {
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()
		var r Rover
		err := decoder.Decode(&r)
		if err != nil {
			log.Fatal(err)
		}
		for i, _ := range r.Manifest.Photos {
			purl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", name, "/photos?sol=", r.Manifest.Photos[i].Sol, "&api_key=", s.APIKey)
			photos := parsePhotos(purl)
			for _, ph := range photos {
				ph.Rover = name
				ph.copyToS3(s.AWSRegion, s.S3Bucket)
				ph.save()
			}
		}
	}
	return nil
}

func parsePhotos(url string) []Photo {
	var pr PhotoResponse
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&pr)
		if err != nil {
			log.Fatal(err)
		}
	}
	return pr.Photos
}
