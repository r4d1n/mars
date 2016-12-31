package main

type photo struct {
	Id        int    `json:"id"`
	Sol       int    `json:"sol"`
	Rover     string `json:"rover"`
	Camera    string `json:"camera"`
	EarthDate string `json:"earth_date"`
	S3ImgSrc  string `json:"img_src"`
}
