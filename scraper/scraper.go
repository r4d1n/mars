package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	data "github.com/r4d1n/nasa-mars-photos/roverdata"
)

type Scraper struct {
	APIKey    string
	AWSRegion string
	S3Bucket  string
}

type manifestResponse struct {
	Manifest Manifest `json:"photo_manifest"`
}

type Manifest struct {
	Name        string
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string
	MaxSol      int       `json:"max_sol"`
	MaxDate     string    `json:"max_date"`
	TotalPhotos int       `json:"total_photos"`
	Sols        data.Sols `json:"photos"`
}

func (s Scraper) crawl(name string) error {
	last, err := checkLastInsert(name)
	fmt.Printf("rover %s: last saved image %d of sol: %d \n", name, last.Id, last.Sol)
	if err != nil {
		return err
	}
	client := &http.Client{}
	murl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", name, "?api_key=", s.APIKey)
	req, err := http.NewRequest("GET", murl, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	} else {
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()
		var r manifestResponse
		err := decoder.Decode(&r)
		sort.Sort(r.Manifest.Sols)
		if err != nil {
			return err
		}
		// make a Sol for most recent photo Sol and get index in manifest sols to find initial loop position
		d := data.Sol{Sol: last.Sol}
		i := r.Manifest.Sols.IndexOf(d)
		// need to advance if nothing has been saved or if all photos have been found
		count, err := checkTotalSaved(name, d.Sol)
		if err != nil {
			return err
		}
		if i == -1 || r.Manifest.Sols[i].TotalPhotos <= count {
			i++ // start at 0 or next sol if necessary
		}
		sollist := make(chan data.Sols)
		photolist := make(chan data.Photos)
		errorch := make(chan error)                      // for handling errorch in concurrent crawling
		done := make(map[int]bool)                       // will check Sol.Id
		go func() { sollist <- r.Manifest.Sols[i+1:] }() // initial position is index of last photo saved + 1
		for list := range sollist {
			for _, sol := range list {
				fmt.Println("!!!!!!! ranging over list", sol.Sol)
				if !done[sol.Sol] {
					done[sol.Sol] = true
					purl := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/rovers/", name, "/photos?sol=", sol.Sol, "&api_key=", s.APIKey)
					go func(uri string) {
						var photos data.Photos
						photos, err = getPhotos(uri)
						sort.Sort(photos)
						if err != nil {
							fmt.Println("error blockS")
							errorch <- err
						}
						photolist <- photos
					}(purl)
					for phos := range photolist {
						for _, ph := range phos {
							ph.Rover = name
							fmt.Printf("ph id: %d / sol: %d / rover: %s \n", ph.Id, ph.Sol, ph.Rover)
							// err := ph.CopyToS3(s.AWSRegion, s.S3Bucket)
							if err != nil {
								errorch <- err
							}
							// err = ph.Save()
							if err != nil {
								errorch <- err
							}
						}
					}
				}
				e := <-errorch
				if e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// find the last saved photo from this rover
func checkLastInsert(rover string) (data.Photo, error) {
	var p data.Photo
	err := db.QueryRow("select id, sol from photos where rover=$1 order by sol desc, id desc limit 1", rover).Scan(&p.Id, &p.Sol)
	if err == sql.ErrNoRows {
		return p, nil
	} else if err != nil {
		return p, fmt.Errorf("retrieving last photo from rover %s: %v", rover, err)
	}
	return p, nil
}

// check how many photos have been saved for a rover on a given sol
func checkTotalSaved(rover string, sol int) (int, error) {
	var count int
	err := db.QueryRow("select count(*) from photos where rover=$1", rover).Scan(&count)
	if err == sql.ErrNoRows {
		return count, nil
	} else if err != nil {
		return count, fmt.Errorf("retrieving count from rover %s for sol %d: %v", rover, sol, err)
	}
	return count, nil
}

type photoResponse struct {
	Photos data.Photos
}

// fetch and parse photos for a given sol
func getPhotos(url string) (data.Photos, error) {
	fmt.Println("fetching url:", url)
	var pr photoResponse
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
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
	sort.Sort(pr.Photos)
	return pr.Photos, nil
}
