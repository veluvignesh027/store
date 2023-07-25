package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const delivery = "/Users/vb/gitstore/store/ContentStore"

func main() {
	meta := MetaData{}

	uri := "http://<ip-addr>" + delivery + "/DeliveryService01"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Name", "cdn1.png")
	resp, _ := meta.Do(req)

	log.Println("\n\nResponse :", resp.Header)

	newfile, _ := os.Create("resp.png")

	_, err = io.Copy(newfile, resp.Body)
	if err != nil {
		log.Println(err)
	}
}
