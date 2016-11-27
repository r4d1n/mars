package main

import "fmt"

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

type PhotoResponse struct {
	Photos []Photo
}

type Photo struct {
	Id        int
	Sol       int
	Rover     string `json:"rover.name"`
	Camera    Camera
	EarthDate string `json:"earth_date"`
	ImgSrc    string `json:"img_src"`
}

type Camera struct {
	Name string
}

func (p *Photo) save() (err error) {
	fmt.Println("saving photo", p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.ImgSrc)
	statement := "INSERT INTO photos (id, sol, rover, camera, earthdate, imgsrc) VALUES($1, $2, $3, $4, $5, $6) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.ImgSrc).Scan(&p.Id)
	return
}
