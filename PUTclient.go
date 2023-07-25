package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const delivery = "/Users/vb/gitstore/store/ContentStore"

func main() {

	meta := MetaData{}

	file, _ := os.Open("simpleCDN.png")
	defer file.Close()

	filebuff := &bytes.Buffer{}

	_, err := io.Copy(filebuff, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	uri := "http://<ip-addr>" + delivery + "/DeliveryService01"

	req, err := http.NewRequest("PUT", uri, filebuff)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Name", "cdn1.png")
	req.Header.Set("Etag", "HTTP-ETG")
	req.Header.Set("Max-Age", "30")
	req.Header.Set("Age", "10")
	resp, _ := meta.Do(req)

	log.Println("\n\nResponse :", resp.Header)

}
