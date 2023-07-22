package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	//"github.com/mailru/easyjson/buffer"
)

type MetaData struct {
	Assetid      string    `json:"ASSET-ID"`
	Aid          int       `json:"ASSET-NUMBER"`
	Etag         string    `json:"ETAG"`
	Lastmodified time.Time `json:"LAST-MODIFIED"`
	Maxage       int       `json:"MAX-AGE"`
	Age          int       `json:"AGE"`
}

var database []MetaData

func main() {

	meta := MetaData{}

	file, _ := os.Open("client.go")
	defer file.Close()

	filebuff := &bytes.Buffer{}

	_, err := io.Copy(filebuff, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("PUT", "http://10.10.10.10/Users/vb/go/db/dir007", filebuff)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Name", "file1.txt")

	resp, _ := meta.Do(req)
	log.Println("\n\nResponse :", resp)
}
