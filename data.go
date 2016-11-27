package main

// structs and functions for organizing and storing data from the nasa rover API

type Sol struct {
	Sol         int
	TotalPhotos int `json:"total_photos"`
	Cameras     []string
}

type Manifest struct {
	PhotoManifest PhotoManifest `json:"photo_manifest"`
}

type PhotoManifest struct {
	Name        string
	LandingDate string `json:"landing_date"`
	LaunchDate  string `json:"launch_date"`
	Status      string
	MaxSol      int    `json:"max_sol"`
	MaxDate     string `json:"max_date"`
	TotalPhotos int    `json:"total_photos"`
	Photos      []Sol
}

type Photo struct {
	Id        int
	Sol       int
	Rover     string
	Camera    string `json:"camera.full_name"`
	EarthDate string `json:"earth_date"`
	ImgSrc    string `json:"img_src"`
}
