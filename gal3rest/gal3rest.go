//Copyright (c) 2012 Tim Shannon
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

//gal3rest packages handles all access to the gallery3 rest api
// and is meant to abstract accessing it
package gal3rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Client struct {
	Url    string
	APIKey string
}

type RestData struct {
	Url           string
	Entity        Entity
	Members       []string
	Relationships map[string]interface{} // may Type this later
}

type RestCreate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Title string `json:"title"`
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
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedValue)))

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

func (gClient *Client) UploadImage(title string, name string, imagePath string, parentUrl string) {
	gClient.checkClient()
	hClient := new(http.Client)

	c := RestCreate{Name: name, Title: title, Type: PHOTO}
	b, jErr := json.Marshal(c)
	if jErr != nil {
		log.Panicln("Error marshalling Rest create: ", jErr)
	}

	entity := "entity=" + url.QueryEscape(string(b))

	file, fErr := ioutil.ReadFile(imagePath)
	if fErr != nil {
		log.Panic("Error reading the image file: ", fErr)
	}

	var dataParts = make([]string, 12)
	boundry := "roPK9J3DoG4ZWP6etiDuJ97h-zeNAph"

	//build multipart request
	dataParts[0] = "--" + boundry
	dataParts[1] = `Content-Disposition: form-data; name="entity"`
	dataParts[2] = "Content-Type: text/plain; charset=UTF-8"
	dataParts[3] = "Content-Transfer-Encoding: 8bit"
	//space in 4
	dataParts[5] = entity
	dataParts[6] = "--" + boundry
	dataParts[7] = `Content-Disposition: form-data; name="file";` +
		`filename="` + path.Base(imagePath) + `"`
	dataParts[8] = "Content-Type: " + getContentType(imagePath)
	dataParts[9] = "Content-Transfer-Encoding: binary"
	//space in 10
	dataParts[11] = string(file)
	dataParts[12] = "--" + boundry

	data := strings.Join(dataParts, "\n")
	buffer := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest("POST", parentUrl, buffer)
	req.Header.Set("X-Gallery-Request-Method", "POST")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundry)
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

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

func getContentType(file string) string {
	ext := file[strings.LastIndex(file, "."):]
	mType := mime.TypeByExtension(ext)
	if mType == "" {
		return "application/octet-stream"
	}
	return mType

}
