package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// structs and functions for scraping the nasa rover API

type Crawler interface {
	loadConfig() Nasa
	manifest() string
	images() string
	savePhoto() string
}

type Nasa struct {
	APIKey string
}

type Rover struct {
	Manifest Manifest `json:"photo_manifest"`
}

type Manifest struct {
	Name        string
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string
	MaxSol      int    `json:"max_sol"`
	MaxDate     string `json:"max_date"`
	TotalPhotos int    `json:"total_photos"`
	Photos      []Sol
}

type Sol struct {
	Sol         int
	TotalPhotos int `json:"total_photos"`
	Cameras     []string
}

func (n Nasa) crawl(s string) error {
	u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", s, "?api_key=", n.APIKey)
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	} else {
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()
		var r Rover
		err := decoder.Decode(&r)
		// log.Print("###### ", r)
		log.Print("###### ", r.Manifest.Photos[666])
		if err != nil {
			log.Fatal(err)
		}
		for i, m := range r.Manifest.Photos {
			fmt.Println("in the range loop", i, m, s)
			u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", s, "/photos?sol=", m.Sol, "&api_key=", n.APIKey)
			p := new(Photo)
			if err := p.seed(u); err != nil {
				return err
			}
			// p.save()
			// a = append(a, p)
		}
	}
	return nil
}
