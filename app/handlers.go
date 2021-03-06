package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var start int

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	rover := "curiosity"
	start = 10
	limit := start
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc, id desc limit $2", rover, limit)
	var data []photo
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p photo
		err = rows.Scan(&p.ID, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, p)
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func getRoverPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []photo
	rover := mux.Vars(r)["rover"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	limit, err := strconv.Atoi(mux.Vars(r)["limit"])
	if err != nil {
		log.Fatal(err)
	}
	page = (page * limit) + start
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc, id desc limit $2 offset $3", rover, limit, page)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p photo
		err = rows.Scan(&p.ID, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
		if err != nil {
			log.Fatal(err)
		}
		photos = append(photos, p)
	}
	j, err := json.Marshal(photos)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(j)
}
