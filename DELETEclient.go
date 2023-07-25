package main

import (
	"log"
	"net/http"
)

const delivery = "/Users/vb/gitstore/store/ContentStore"

func main() {

	meta := MetaData{}

	uri := "http://<ip-addr>" + delivery + "/DeliveryService01"

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Name", "cdn1.png")
	resp, _ := meta.Do(req)

	log.Println("\n\nResponse :", resp.Header)

}
