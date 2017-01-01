package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	rover := "curiosity"
	limit := 10
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc limit $2 offset $3", rover, limit, 10)
	var data []photo
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p photo
		err = rows.Scan(&p.Id, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, p)
	}
	t.Execute(w, data)
}

func getRoverPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []photo
	rover := mux.Vars(r)["rover"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	limit := 10
	page = page * limit
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc limit $2 offset $3", rover, limit, page)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p photo
		err = rows.Scan(&p.Id, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d, %s, %d, %s, %s\n", p.Id, p.Rover, p.Sol, p.Camera, p.S3ImgSrc)
		photos = append(photos, p)
	}
	j, err := json.Marshal(photos)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}
