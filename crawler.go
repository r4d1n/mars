package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Crawler interface {
	loadConfig() Nasa
	manifest() string
	images() string
	savePhoto() string
}

type Nasa struct {
	APIKey string
}

func (n Nasa) crawl(s string) error {
	url1 := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", s, "?api_key=", n.APIKey)
	res, err := http.Get(url1)
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
		for i, m := range r.Manifest.Photos {
			url2 := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", s, "/photos?sol=", r.Manifest.Photos[i].Sol, "&api_key=", n.APIKey)
			photos := parsePhotos(url2)
			for _, ph := range photos {
				ph.Rover = s
				ph.save()
			}
		}
	}
	return nil
}

func parsePhotos(u string) []Photo {
	var pr PhotoResponse
	res, err := http.Get(u)
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
