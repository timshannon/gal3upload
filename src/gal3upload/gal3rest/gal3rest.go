package gal3rest

import (
	"log"
	"http"
	"strings"
	"io/ioutil"
	"json"
)

func TestConnection(url string) interface{} {
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
	response.Body.Close()
	if err != nil {
		log.Panic("Error reading response: ", err)
	}
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Panic("Error unmarshelling json data: ", err)
	}

	return data
}

type RestData struct {
	Url           string
	Entity        map[string]string
	Relationships map[string]string
	Members       map[string]string
}
