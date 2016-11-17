package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=rover_user password=notamartian dbname=nasa_rover_data sslmode=disable")
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

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
	savePhoto() string
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

type Photo struct {
	Id        int
	Sol       int
	Rover     string
	Camera    string `json:"camera.full_name"`
	EarthDate string `json:"earth_date"`
	ImgSrc    string `json:"img_src"`
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

func (nc NasaCrawler) getManifest(rover string) PhotoManifest {
	u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", rover, "?api_key=", nc.APIKey)
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
	}
	fmt.Println("got the manifest")
	return m.PhotoManifest
}

func (nc NasaCrawler) getPhotoData(pm PhotoManifest) []*Photo {
	slc := make([]*Photo, len(pm.Photos))
	client := &http.Client{}
	for i, s := range pm.Photos {
		fmt.Println("in the range loop", i, s)
		u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", pm.Name, "/photos?sol=", s.Sol, "&api_key=", nc.APIKey)
		fmt.Println("url", u)
		request, err := http.NewRequest("GET", u, nil)
		request.Close = true
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("first error")
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			p := new(Photo)
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&p)
			if err != nil {
				fmt.Println("second error")
				log.Fatal(err)
			}
			fmt.Println(p)
			// nc.savePhoto(*p)
			slc = append(slc, p)
		}
	}
	return slc
}

// func (nc NasaCrawler) savePhoto(p Photo) {
// 	result, err := db.Exec("INSERT INTO photos VALUES($1, $2, $3, $4, $5)", p.Id, p.Sol, p.Camera, p.EarthDate, p.ImgSrc)
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		fmt.Println("saved image", p.Id, result)
// 	}
// }
