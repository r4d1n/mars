package main

import "fmt"

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
	fmt.Println("time to save", p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.ImgSrc)
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
