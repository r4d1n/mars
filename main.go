package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := new(Crawler)
	c.LoadConfig("./config.json")
	c.GetManifest("curiosity")
}

type Fetcher interface {
	LoadConfig() Crawler
	GetManifest() string
	GetPhotoData() string
}

type Crawler struct {
	APIKey string
}

type Sol struct {
	sol          int
	total_photos int
	cameras      []string
}

func (c *Crawler) LoadConfig(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
}

func (c Crawler) GetManifest(rover string) {
	fmt.Sprint("c", c)
	u := fmt.Sprint("https://api.nasa.gov/mars-photos/api/v1/manifests/", rover, "?api_key=", c.APIKey)
	response, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err.Error())
		}
		var f interface{}
		err := json.NewDecoder(response.Body).Decode(&f)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(f)
		m := f.(map[string]interface{})
		for k, v := range m {
			fmt.Println(k, v)
			// switch vv := v.(type) {
			// case string:
			// 	fmt.Println(k, "is string", vv)
			// case int:
			// 	fmt.Println(k, "is int", vv)
			// case []interface{}:
			// 	fmt.Println(k, "is an array:")
			// 	for i, u := range vv {
			// 		fmt.Println(i, u)
			// 	}
			// default:
			// 	fmt.Println(k, "is of a type I don't know how to handle")
			// }
		}
	}
}

func (c Crawler) GetPhotoData() {

}
