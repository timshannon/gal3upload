package gal3rest

import (
	"log"
	"http"
	"strings"
	"io/ioutil"
	"json"
)

type Client struct {
	Url    string
	APIKey string
}

type Album struct {
	AlbumUrl string
	Photos   []string
}

func (gClient *Client) GetMembers(itemId int) []string {
	hClient := new(http.Client)
	reader := new(strings.Reader) //TODO: build actual content
	if gClient.Url == "" {
		log.Panic("No url defined for this gallery Client. " +
			"Be sure to set .Url before connectiong.")
	}
	if gClient.Url[:1] != "/" {
		gClient.Url += "/"
	}
	req, _ := http.NewRequest("GET", gClient.Url+"rest/item/"+string(itemId), reader)
	req.Header.Set("X-Gallery-Request-Method", "GET")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)
	response, err := hClient.Do(req)
	if err != nil {
		log.Panic("Error connecting to: "+gClient.Url+" Error: ", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Panic("Error reading response: ", err)
	}

	var data interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		//don't care about some values?
		log.Println("Error unmarshalling json data: ", err)
	}
	allData := data.(map[string]interface{})
	m := allData["members"].([]interface{})

	var members []string
	for i := range m {
		members = append(members, m[i].(string))
	}
	return members

}
