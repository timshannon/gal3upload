package gal3rest

import (
	"log"
	"http"
	"strings"
	"io/ioutil"
	"json"
	"fmt"
)

func TestConnection(url string) (RestResponse, interface{}) {
	client := new(http.Client)
	reader := new(strings.Reader) //TODO: build actual content
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
	var data RestResponse
	var unmarshalled interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		//don't care about some values?
		fmt.Println("Error unmarshelling json data: ", err)
	}

	err = json.Unmarshal(body, &unmarshalled)
	return data, unmarshalled
}

type RestResponse struct {
	Url           string
	Entity        map[string]string
	Relationships map[string]string
	Members       []string
}
