package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var static string

func init() {
	flag.StringVar(&static, "static_dir", "static", "the directory from which to serve static content")
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", os.Getenv("NASA_DB_USER"), os.Getenv("NASA_DB_PASS"), os.Getenv("NASA_DB_NAME"), os.Getenv("NASA_DB_HOST"), os.Getenv("NASA_DB_SSL")))
	if err = db.Ping(); err != nil {
		fmt.Println("an error occurred in init")
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
	r.HandleFunc("/rover/{rover}/limit/{limit}/page/{page}", getRoverPhotos)
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
