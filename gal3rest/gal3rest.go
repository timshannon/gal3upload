package gal3rest

import (
	"log"
	"http"
	"io/ioutil"
	"strconv"
	"json"
	"bytes"
	"url"
)

type Client struct {
	Url    string
	APIKey string
}

type RestData struct {
	Url           string
	Entity        Entity
	Members       []string
	Relationships map[string]interface{} // may type this later
}

type RestCreate struct {
	Name  string
	Type  string
	Title string
}

type Entity struct {
	Id               int
	Description      string
	Name             string
	Web_url          string
	Mime_type        string
	Title            string
	Type             string
	Album_cover      string
	Thumb_url        string
	Thumb_url_public string
	Width            int
	Thumb_width      int
	Resize_width     int
	Resize_height    int
	View_count       int
	Sort_order       int
	Height           int
	Updated          int
	Captured         int
	View_1           string
	Can_edit         int
	View_2           string
	Thumb_size       int
	Level            int
	Created          int
	Sort_column      string
	Slug             string
	Rand_key         int
	Thumb_height     int
	Owner_id         int
}

type Album struct {
	Entity
	Photos []*Photo
	Albums []string
}

const (
	PHOTO = "photo"
	ALBUM = "album"
)

type Photo struct {
	Entity
	ImageData []byte
}

func NewClient(url string, apiKey string) Client {
	client := Client{Url: url, APIKey: apiKey}
	client.checkClient()

	return client
}
func (gClient *Client) GetRESTItem(itemUrl string) *RestData {
	hClient := new(http.Client)
	req, _ := http.NewRequest("GET", itemUrl, nil)
	req.Header.Set("X-Gallery-Request-Method", "GET")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)
	response, err := hClient.Do(req)
	if err != nil {
		log.Panic("Error connecting to: "+itemUrl+" Error: ", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Panic("Error reading response: ", err)
	}

	data := new(RestData)

	json.Unmarshal(body, &data)

	return data
}

func (gClient *Client) checkClient() {
	if gClient.Url == "" {
		log.Panicf("No URL specified in the client." +
			" Be sure to specify the REST url before making a request")
	} else {
		if gClient.Url[:1] != "/" {
			gClient.Url += "/"
		}
	}
	if gClient.APIKey == "" {
		log.Panicln("No API key specified in the client. " +
			"Be sure to specify the REST API key before making a request.")
	}
}

func (gClient *Client) GetUrlFromId(id int) string {
	gClient.checkClient()
	return gClient.Url + "rest/item/" + strconv.Itoa(id)
}

func (gClient *Client) GetAlbum(itemUrl string) *Album {
	data := gClient.GetRESTItem(itemUrl)
	album := new(Album)
	album.Entity = data.Entity
	for i := range data.Members {
		mData := gClient.GetRESTItem(data.Members[i])
		if mData.Entity.Type == PHOTO {
			photo := new(Photo)
			photo.Entity = mData.Entity
			album.Photos = append(album.Photos, photo)
		} else if mData.Entity.Type == ALBUM {
			album.Albums = append(album.Albums, data.Members[i])
		}
	}
	return album
}

func (gClient *Client) CreateAlbum(title string, name string, parentUrl string) {
	gClient.checkClient()
	hClient := new(http.Client)

	c := RestCreate{Name: name, Title: title, Type: ALBUM}
	b, jErr := json.Marshal(c)
	if jErr != nil {
		log.Panicln("Error marshalling Rest create: ", jErr)
	}

	//base64.URLEncoding.EncodeToString	
	encodedValue := "entity=" + url.QueryEscape(string(b))

	buffer := bytes.NewBuffer([]byte(encodedValue))

	req, _ := http.NewRequest("POST", parentUrl, buffer)
	req.Header.Set("X-Gallery-Request-Method", "POST")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))

	log.Println("request: ", req)
	response, err := hClient.Do(req)
	if err != nil {
		log.Panic("Error connecting to: "+parentUrl+" Error: ", err)
	}
	_, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Panic("Error reading response: ", err)
	}

}
