package gal3rest

import (
	"log"
	"http"
	"strings"
	"io/ioutil"
	"strconv"
	"json"
)

type Client struct {
	Url    string
	APIKey string
}

type RestData struct {
	Url           string
	Entity        Entity
	Members       []interface{}
	Relationships map[string]interface{}
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
	Albums []int
}

const (
	PHOTO = "photo"
	ALBUM = "album"
)

type Photo struct {
	Entity
	ImageData []byte
}

func (gClient *Client) GetRESTItem(itemUrl string) *RestData {
	hClient := new(http.Client)
	reader := new(strings.Reader)
	req, _ := http.NewRequest("GET", itemUrl, reader)
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

func (r *RestData) GetMembers() []string {
	members := make([]string, len(r.Members))
	for i := range r.Members {
		members[i] = r.Members[i].(string)
	}
	return members
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

func (gClient *Client) GetAlbum(id int) *Album {
	data := gClient.GetRESTItem(gClient.GetUrlFromId(id))
	album := new(Album)
	album.Entity = data.Entity

	members := data.GetMembers()

	for i := range members {
		data = gClient.GetRESTItem(members[i])
		if data.Entity.Type == PHOTO {
			photo := new(Photo)
			photo.Entity = data.Entity
			album.Photos = append(album.Photos, photo)
		} else if data.Entity.Type == ALBUM {
			album.Albums = append(album.Albums, data.Entity.Id)
		}

	}

	return album
}
