package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const aidvar = 0

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

	resp := &http.Response{
		StatusCode: http.StatusMethodNotAllowed,
		Status:     http.StatusText(http.StatusMethodNotAllowed),
	}
	return resp, http.ErrNotSupported
}

func handlePUTConnection(req *http.Request) (*http.Response, error) {
	content := req.Header.Get("Content-Name")

	log.Println("PUT request for Content : ", content)

	resp := &http.Response{StatusCode: http.StatusBadRequest, Header: make(http.Header)}

	if CheckObject(req.URL.Path + "/" + content) {
		log.Println("Data Already there in memory!!! Checking for Updated data....")
		resp.StatusCode = http.StatusNotModified
		resp.Status = http.StatusText(http.StatusNotModified)
		return resp, nil
	} else if CheckDir(req.URL.Path) {
		log.Println("Creating objects in Directory...")
		newfile, _ := os.Create(req.URL.Path + "/" + content)
		_, err := io.Copy(newfile, req.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("Object Saved succesfully as ", newfile.Name())

		var tmp MetaData
		tmp.Assetid = content
		tmp.Aid = aidvar + 1
		tmp.Lastmodified = time.Now()
		tmp.Etag = req.Header.Get("Etag")
		tmp.Age, _ = strconv.Atoi(req.Header.Get("Age"))
		tmp.Maxage, _ = strconv.Atoi(req.Header.Get("Max-Age"))
		err = saveToFile(tmp, req.URL.Path+"/"+"Delivermetadata.json")
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusAccepted)
		return resp, err

	} else {
		if CreateDir(req.URL.Path) {
			newfile, _ := os.Create(req.URL.Path + "/" + content)
			_, err := io.Copy(newfile, req.Body)
			if err != nil {
				log.Println(err)
			}
			log.Println("Object Saved succesfully as ", newfile.Name())
			var tmp MetaData
			tmp.Assetid = content
			tmp.Aid = aidvar + 1
			tmp.Lastmodified = time.Now()
			tmp.Etag = req.Header.Get("Etag")
			err = saveToFile(tmp, req.URL.Path+"/"+"Delivermetadata.json")
			if err != nil {
				log.Println(err)
				return nil, err
			}
			resp.StatusCode = http.StatusOK
			resp.Status = http.StatusText(http.StatusAccepted)
			return resp, err
		}
	}
	return resp, nil
}

func handleGETConnection(req *http.Request) (*http.Response, error) {
	content := req.Header.Get("Content-Name")
	log.Println("GET request for content : ", content)

	m, got := loadfromFile(req.Header.Get("Content-Name"), req.URL.Path+"/"+"Delivermetadata.json")
	if got {
		log.Println("Got the object from the db")
		log.Println(m)
		resp := &http.Response{
			StatusCode: http.StatusFound,
			Header:     make(http.Header),
		}
		resp.Header.Set("Content-Type", "Image")
		resp.Header.Set("ETAG", m.Etag)
		resp.Header.Set("Last-Modified", m.Lastmodified.String())
		resp.Header.Set("Age", strconv.Itoa(m.Age))
		resp.Header.Set("Max-Age", strconv.Itoa(m.Maxage))
		log.Println("Getting content data from : ", req.URL.Path+"/"+content)
		file, _ := os.Open(req.URL.Path + "/" + content)
		defer file.Close()

		filebuff := &bytes.Buffer{}
		_, err := io.Copy(filebuff, file)
		if err != nil {
			log.Println(err)
			return resp, err
		}
		resp.Body = ioutil.NopCloser(filebuff)
		return resp, nil
	} else {
		log.Println("Failed to fetch from the db")
		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
		resp.Header.Set("Content-Type", "Image")
		resp.Header.Set("ETAG", m.Etag)
		resp.Header.Set("Last-Modified", m.Lastmodified.String())
		resp.Header.Set("Age", strconv.Itoa(m.Age))
		resp.Header.Set("Max-Age", strconv.Itoa(m.Maxage))
		return resp, nil
	}

}

func handleDELETEConnection(req *http.Request) (*http.Response, error) {
	content := req.Header.Get("Content-Name")
	path := req.URL.Path
	log.Println("DELETE request", content)

	resp := &http.Response{Header: make(http.Header)}

	if CheckObject(path + "/" + content) {

		err := os.Remove(path + "/" + content)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Successfully deleted content file: ", content)
		}

		ret, err := deleteFromFile(content, req.URL.Path+"/"+"Delivermetadata.json")
		if err != nil {
			log.Println(err)
		}
		if ret {
			log.Println("Deleted Successfully!")
		} else {
			log.Println("Can't Delete..")
			resp.StatusCode = http.StatusNotFound
			return resp, err
		}

		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
		resp.Header.Set("Content-Name", content)
		return resp, err
	} else {
		resp.StatusCode = http.StatusNotFound
		return resp, nil
	}

	return resp, nil
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

func saveToFile(md MetaData, f string) error {
	fd, err := os.OpenFile(f, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fd.Close()

	data, err := json.Marshal(md)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = fd.Write(append(data, '\n'))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func deleteFromFile(cont string, f string) (bool, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return false, err
	}

	lines := bytesToLines(data)
	var newData []byte
	deleted := false

	for _, line := range lines {
		var tmp MetaData
		err := json.Unmarshal(line, &tmp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if tmp.Assetid == cont {
			deleted = true
			continue
		}
		newData = append(newData, line...)
		newData = append(newData, '\n')
	}
	err = ioutil.WriteFile(f, newData, 0644)
	if err != nil {
		return false, err
	}
	return deleted, nil
}

func loadfromFile(assetname string, filename string) (MetaData, bool) {
	var temp MetaData
	found := false
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return temp, found
	}

	lines := bytesToLines(data)
	for _, line := range lines {
		err := json.Unmarshal(line, &temp)
		if err != nil {
			fmt.Println("Error unmarshalling JSON data:", err)
			continue
		}

		if temp.Assetid == assetname {
			temp = temp
			found = true
			break
		}
	}

	return temp, found
}
func bytesToLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
