package gal3rest

import (
	"log"
	"http"
	"strings"
	"io/ioutil"
	"json"
)

func TestConnection(url string) string {
	client := new(http.Client)
	reader := new(strings.Reader) //TODO: build actuall content
	if url[:1] != "/" {
		url += "/"
	}
	req, _ := http.NewRequest("GET", url+"rest/item/1", reader)
	req.Header.Set("X-Gallery-Request-Method", "GET")
	req.Header.Set("X-Gallery-Request-Key", "79daf60695177e16ff2480f8338b5fcc")
	response, err := client.Do(req)

	if err != nil {
		log.Panic("Error connecting to: "+url+" Error: ", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Close()
	if err != nil {
		log.Panic("Error reading response: ", err)
	}

	json.Unmarshal(body, result)
	return string(result)

}

type RestData struct {
	url           string
	entity        map[string]string
	relationships map[string]string
	members       map[string]string
}
