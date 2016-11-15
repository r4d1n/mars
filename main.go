package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := new(NasaCrawler)
	c.loadConfig("./config.json")
	pm := c.getManifest("curiosity")
	c.getPhotoData(pm)
}

type Fetcher interface {
	loadConfig() NasaCrawler
	getManifest() string
	getPhotoData() string
}

type NasaCrawler struct {
	APIKey string
}

type Sol struct {
	Sol         int
	TotalPhotos int `json:"total_photos"`
	Cameras     []string
}

type Manifest struct {
	PhotoManifest PhotoManifest `json:"photo_manifest"`
}

type PhotoManifest struct {
	Name        string
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string
	MaxSol      int    `json:"max_sol"`
	MaxDate     string `json:"max_date"`
	TotalPhotos int    `json:"total_photos"`
	Photos      []Sol
}

func (c *NasaCrawler) loadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}

func (c NasaCrawler) getManifest(rover string) PhotoManifest {
	u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", rover, "?api_key=", c.APIKey)
	m := new(Manifest)
	response, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(m.PhotoManifest.Photos[4].Cameras)
	}
	return m.PhotoManifest
}

func (c NasaCrawler) getPhotoData(pm PhotoManifest) {
	for _, entry := range pm.Photos {
		fmt.Println(entry)
	}
}
