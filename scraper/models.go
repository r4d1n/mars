/*
Data structures for scraping the nasa mars rover photo api
*/

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

// the photos are organized by martian Sol
type Sol struct {
	Sol         int
	TotalPhotos int `json:"total_photos"`
}

// Sols implements the sort.Sort interface
type Sols []*Sol

func (slice Sols) Len() int {
	return len(slice)
}

func (slice Sols) Less(i, j int) bool {
	return slice[i].Sol < slice[j].Sol
}

func (slice Sols) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (c Sols) IndexOf(s Sol) int {
	for i, val := range c {
		if val.Sol == s.Sol {
			return i
		}
	}
	return -1
}

// Photo represents a single image and associated metadata
type Photo struct {
	Id         int
	Sol        int
	Rover      string `json:"rover.name"`
	Camera     Camera
	EarthDate  string `json:"earth_date"`
	NasaImgSrc string `json:"img_src"`
	S3ImgSrc   string
}

// Photos implements the sort.Sort interface
type Photos []*Photo

func (slice Photos) Len() int {
	return len(slice)
}

func (slice Photos) Less(i, j int) bool {
	return slice[i].Id < slice[j].Id
}

func (slice Photos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (c Photos) IndexOf(p Photo) int {
	for i, val := range c {
		if val.Id == p.Id {
			return i
		}
	}
	return -1
}

// Each photo was taken by an associated Camera
type Camera struct {
	Name string
}

func (p *Photo) Save() (err error) {
	statement := "INSERT INTO photos (id, sol, rover, camera, earthdate, nasaimgsrc, s3imgsrc) VALUES($1, $2, $3, $4, $5, $6, $7) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return fmt.Errorf("saving image %d to db: %v", p.Id, err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.Id, p.Sol, p.Rover, p.Camera.Name, p.EarthDate, p.NasaImgSrc, p.S3ImgSrc).Scan(&p.Id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("saving image %d to db: %v", p.Id, err)
	} else {
		log.Printf("successfully saved data for image %d \n", p.Id)
	}
	return
}

func (p *Photo) CopyToS3(region string, bucket string) (err error) {
	res, err := http.Get(p.NasaImgSrc)
	if err != nil {
		return fmt.Errorf("retrieving image %d from nasa: %v", p.Id, err)
	} else {
		defer res.Body.Close()
		reader := bufio.NewReader(res.Body)
		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(region)}))
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   reader,
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/%d.jpg", p.Rover, p.Id)),
		})
		if err != nil {
			return fmt.Errorf("uploading image %d to s3: %v", p.Id, err)
		}
		log.Printf("completed upload to s3 url: %s \n", result.Location)
		p.S3ImgSrc = result.Location
	}
	return
}
