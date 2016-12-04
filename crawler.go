package main

import (
	"encoding/json"
	"fmt"
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
			photos, err := parsePhotos(purl)
			if err != nil {
				return err
			}
			for _, ph := range photos {
				ph.Rover = name
				err := ph.copyToS3(s.AWSRegion, s.S3Bucket)
				if err != nil {
					return err
				}
				err = ph.save()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type photoResponse struct {
	Photos []Photo
}

func parsePhotos(url string) ([]Photo, error) {
	var pr photoResponse
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("retrieving url %s: %v", url, err)
	} else {
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&pr)
		if err != nil {
			return nil, fmt.Errorf("parsing photo json: %v", err)
		}
	}
	return pr.Photos, nil
}
