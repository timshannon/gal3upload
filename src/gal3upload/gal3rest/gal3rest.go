package gal3rest

import (
	"log"
	"http"
	"os"
	"strings"
)

func TestConnection(url string) *http.Response {
	client := new(http.Client)
	reader := new(strings.Reader)
	req, _ := http.NewRequest("GET", url+"/rest/item/1", reader)
	req.Header.Set("X-Gallery-Request-Method", "GET")
	req.Header.Set("X-Gallery-Request-Key", "79daf60695177e16ff2480f8338b5fcc")
	response, err := client.Do(req)
	if err != nil {
		log.Panic("Error connecting to: "+url+" Error: ", err)
	}
	return response
}

type GalleryRequest struct{}

func (r GalleryRequest) Read(p []byte) (n int, err os.Error) {
	return 0, nil
}
