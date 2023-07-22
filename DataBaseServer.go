package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type MetaData struct {
	Assetid      string    `json:"ASSET-ID"`
	Aid          int       `json:"ASSET-NUMBER"`
	Etag         string    `json:"ETAG"`
	Lastmodified time.Time `json:"LAST-MODIFIED"`
	Maxage       int       `json:"MAX-AGE"`
	Age          int       `json:"AGE"`
}

func (db *MetaData) Do(req *http.Request) (*http.Response, error) {
	if req.Method == "GET" {
		return handleGETConnection(req)
	} else if req.Method == "PUT" {
		return handlePUTConnection(req)
	} else if req.Method == "DELETE" {
		return handleDELETEConnection(req)
	} else if req.Method == "HEAD" {
		return handleHEADConnection(req)
	}
	return nil, nil
}

func handleGETConnection(req *http.Request) (*http.Response, error) {
	log.Println("GET request", req)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "Plain/Text")
	resp.Header.Set("ETAG", "YSJSS273")
	return resp, nil
}

func handlePUTConnection(req *http.Request) (*http.Response, error) {
	log.Println("PUT request : ", req.URL.Path, "For Content : ", req.Header.Get("Content-Name"))

	if CheckObject(req.URL.Path + "/" + req.Header.Get("Content-Name")) {
		log.Println("Data Already there in memory!!! Checking for Updated data....")

	} else if CheckDir(req.URL.Path) {
		log.Println("Creating objects in Directory...")
		newfile, _ := os.Create(req.URL.Path + "/" + req.Header.Get("Content-Name"))
		fileBuffer := &bytes.Buffer{}

		_, err := io.Copy(fileBuffer, req.Body)
		if err != nil {
			log.Println(err)
		}

		_, err = io.Copy(newfile, fileBuffer)
		if err != nil {
			log.Println(err)
		}
		log.Println("Object Saved succesfully as ", newfile.Name())
	} else {
		if CreateDir(req.URL.Path) {
			newfile, _ := os.Create(req.URL.Path + "/" + req.Header.Get("Content-Name"))
			fileBuffer := &bytes.Buffer{}
			_, err := io.Copy(fileBuffer, req.Body)
			if err != nil {
				log.Println(err)
			}
			_, err = io.Copy(newfile, fileBuffer)
			if err != nil {
				log.Println(err)
			}
			log.Println("Object Saved succesfully as ", newfile.Name())
		}
	}
	return nil, nil
}

func handleDELETEConnection(req *http.Request) (*http.Response, error) {
	log.Println("DELETE request", req)
	return nil, nil
}

func handleHEADConnection(req *http.Request) (*http.Response, error) {
	log.Println("HEAD request", req)

	return nil, nil
}

//Utils

func CheckDir(path string) bool {
	log.Println("Checking Directory availablity...")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Directory NotFound!")
		return false
	} else {
		log.Println("Directory Found...")
		return true
	}
}

func CreateDir(path string) bool {
	log.Println("Creating new directory...", path)
	err := os.Mkdir(path, 0777)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Directory Created sucessfully!!")
	return true
}

func CheckObject(path string) bool {
	log.Println("Checking the object present or not...?")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Object Not found! Need to create as new object..")
		return false
	} else {
		log.Println("Object Found! Check for Modified or not....")
		return true
	}
}
