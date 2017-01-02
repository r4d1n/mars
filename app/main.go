package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/rover/{rover}/page/{page}", getRoverPhotos)
	r.HandleFunc("/", serveIndex)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("serving on port 8080")
	log.Fatal(server.ListenAndServe())
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
