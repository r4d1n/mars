package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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
	Id         int
	Sol        int
	Rover      string `json:"rover.name"`
	Camera     Camera
	EarthDate  string `json:"earth_date"`
	NasaImgSrc string `json:"img_src"`
	S3ImgSrc   string
}

type Camera struct {
	Name string
}

func (p *Photo) save() (err error) {
	fmt.Println("saving photo", p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.NasaImgSrc, p.S3ImgSrc)
	statement := "INSERT INTO photos (id, sol, rover, camera, earthdate, nasaimgsrc) VALUES($1, $2, $3, $4, $5, $6) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.NasaImgSrc, p.S3ImgSrc).Scan(&p.Id)
	return
}

func (p *Photo) copyToS3() (err error) {
	res, err := http.Get(p.NasaImgSrc)
	if err != nil {
		return err
	} else {
		defer res.Body.Close()
		reader := bufio.NewReader(res.Body)
		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   reader,
			Bucket: aws.String("nasa-rover-photos"),
			Key:    aws.String(fmt.Sprint(p.Id)),
		})
		if err != nil {
			log.Fatalln("Failed to upload", err)
		}
		log.Println("Successfully uploaded to", result.Location)
		p.S3ImgSrc = result.Location
	}
	return
}
