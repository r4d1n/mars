package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// structs and functions for organizing and storing data from the nasa rover API

type Photo struct {
	Id        int
	Sol       int
	Rover     string
	Camera    string `json:"camera.full_name"`
	EarthDate string `json:"earth_date"`
	ImgSrc    string `json:"img_src"`
}

func (p Photo) save() error {
	result, err := db.Exec("INSERT INTO photos VALUES($1, $2, $3, $4, $5)", p.Id, p.Sol, p.Camera, p.EarthDate, p.ImgSrc)
	if err != nil {
		return err
	} else {
		fmt.Println("saved image", p.Id, result)
	}
	return nil
}

func (p Photo) seed(u string) error {
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&p)
	if err != nil {
		fmt.Println("second error")
		return err
	}
	fmt.Println("seeded", p, &p)
	return nil
}
