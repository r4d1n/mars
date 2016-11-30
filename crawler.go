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
		return err
	} else {
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()
		var r manifestResponse
		err := decoder.Decode(&r)
		if err != nil {
			return err
		}
		for _, sol := range r.Manifest.Photos {
			purl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", name, "/photos?sol=", sol.Sol, "&api_key=", s.APIKey)
			photos := parsePhotos(purl)
			for _, ph := range photos {
				ph.Rover = name
				ph.copyToS3(s.AWSRegion, s.S3Bucket)
				ph.save()
			}
			sol.save()
		}
	}
	return nil
}

type photoResponse struct {
	Photos []Photo
}

func parsePhotos(url string) []Photo {
	var pr photoResponse
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
