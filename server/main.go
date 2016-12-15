package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var c config

func init() {
	c.load("./config.json")
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", c.DBUser, c.DBPass, c.DBName, c.DBHost, c.SSLMode))
	if err = db.Ping(); err != nil {
		fmt.Println("an error occurred in init")
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rover/{rover}/page/{page}", getRoverPhotos)
	server := &http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: r,
	}
	fmt.Println("serving on port 3000")
	log.Fatal(server.ListenAndServe())
}

type Photo struct {
	Id        int    `json:"id"`
	Sol       int    `json:"sol"`
	Rover     string `json:"rover"`
	Camera    string `json:"camera"`
	EarthDate string `json:"earth_date"`
	S3ImgSrc  string `json:"img_src"`
}

func getRoverPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []Photo
	rover := mux.Vars(r)["rover"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	page = page * 10
	limit := 10
	rows, err := db.Query("SELECT id, sol, rover, camera, earthdate, s3imgsrc FROM photos WHERE rover=$1 order by sol desc limit $2 offset $3", rover, limit, page)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Photo
		err := rows.Scan(&p.Id, &p.Sol, &p.Rover, &p.Camera, &p.EarthDate, &p.S3ImgSrc)
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

type config struct {
	APIKey  string
	DBName  string
	DBUser  string
	DBPass  string
	SSLMode string
	DBHost  string
}

func (c *config) load(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("missing config file: ", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("could not parse config file: ", err)
	}
}
