package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type Scraper struct {
	APIKey    string
	AWSRegion string
	S3Bucket  string
}

func (s Scraper) crawl(name string) error {
	last, err := checkLastInsert(name)
	fmt.Println(last)
	if err != nil {
		return err
	}
	murl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", name, "?api_key=", s.APIKey)
	fmt.Println(murl)
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
		for i := last.Sol; i < len(r.Manifest.Photos); i++ {
			purl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", name, "/photos?sol=", r.Manifest.Photos[i].Sol, "&api_key=", s.APIKey)
			fmt.Println(purl)
			photos, err := parsePhotos(purl)
			if err != nil {
				return err
			}
			sort.Sort(photos)
			for j := last.Id; j < len(photos); j++ {
				ph := photos[j]
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

// find the last saved photo from this rover
func checkLastInsert(rover string) (Photo, error) {
	var p Photo
	err := db.QueryRow("select id, sol from photos order by id desc limit 1").Scan(&p.Id, &p.Sol)
	if err != nil {
		return p, fmt.Errorf("retrieving last photo from rover %s: %v", rover, err)
	}
	return p, nil
}

type photoResponse struct {
	Photos []Photo
}

func parsePhotos(url string) (Photos, error) {
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
