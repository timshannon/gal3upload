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
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
)

//Client holds the url and API keys for the gallery being used
// so you don't have to specify the url and api key for each request
// all gallery access functions are attached to this type
type Client struct {
	Url    string
	APIKey string
}

//RestData contains the general format of data returned from REST
// API calls to Gallery3
type RestData struct {
	Url           string
	Entity        Entity
	Members       []string
	Relationships map[string]interface{} // may Type this later
}

//RestCreate holds the structure for REST API requests to create a new
// gallery
type RestCreate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

//Entity is a general json data structure that is attached to most items
// in the gallery3 rest api
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

//String function for pretty printing the entity REST  data
func (entity *Entity) String() string {
	var strValue string
	ref := reflect.ValueOf(entity).Elem()
	entityType := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		strValue += fmt.Sprintf("%s: %s  %v\n",
			entityType.Field(i).Name,
			field.Type(), field.Interface())
	}
	return strValue
}

//Album is the somewhat abstracted type used to hold the relationship
// between albums and photos, as well as the REST API data associated to them
type Album struct {
	Entity
	Photos []*Photo
	Albums []*Album
}

//entity types
const (
	PHOTO = "photo"
	ALBUM = "album"
)

//Photo is the abstraction of the REST entity
type Photo struct {
	Entity
}

//Response is the json response from creating an album or uploading a photo
// it simply contains the REST URL to the item added
type Response struct {
	Url string
}

//Returns a new client for accessing the Galllery3 REST API
func NewClient(url string, apiKey string) Client {
	client := Client{Url: url, APIKey: apiKey}
	client.checkClient()

	return client
}

//GetRESTItem retrieves an arbitrary REST item from a Gallery3
// Possible key value parameters:
// scope: direct - only specific item
//	  all - All descendants (recursive) doesn't seem to work
// name: value - only return items containing value 
// random: true - returns a single random item
// type: photo
//       movie
//       album
func (gClient *Client) GetRESTItem(itemUrl string,
	parameters map[string]string) (restData *RestData, status int, err error) {
	restData = new(RestData)
	hClient := new(http.Client)

	if parameters != nil {
		urlValues := url.Values{}
		for k, v := range parameters {
			urlValues.Set(k, v)
		}
		itemUrl += "?" + urlValues.Encode()
	}

	req, _ := http.NewRequest("GET", itemUrl, nil)
	req.Header.Set("X-Gallery-Request-Method", "GET")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)

	response, err := hClient.Do(req)
	if err != nil {
		return
		//	log.Panic("Error connecting to: "+itemUrl+" Error: ", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return
		//log.Panic("Error reading response: ", err)
	}

	json.Unmarshal(body, &restData)
	status = response.StatusCode
	return
}

//checkClient checks to make sure the URL and API is set and configured properly
func (gClient *Client) checkClient() {
	if gClient.Url == "" {
		log.Panicln("No URL specified in the client." +
			" Be sure to specify the REST url before making a request")
	} else {
		if gClient.Url[len(gClient.Url)-1:] != "/" {
			gClient.Url += "/"
		}
	}
	if gClient.APIKey == "" {
		log.Panicln("No API key specified in the client. " +
			"Be sure to specify the REST API key before making a request.")
	}
}

//GetUrlFromId simply builds the REST url from the passed in ID
func (gClient *Client) GetUrlFromId(id int) string {
	gClient.checkClient()

	return gClient.Url + "rest/item/" + strconv.Itoa(id)
}

func (gClient *Client) GetItemsUrl() string {
	gClient.checkClient()

	return gClient.Url + "rest/items?"
}

//CreateAlbum creates an album with the passed in name and title inside the album at the
// passed in url
func (gClient *Client) CreateAlbum(title string, name string, parentUrl string) (itemUrl string, status int, err error) {
	gClient.checkClient()
	hClient := new(http.Client)

	c := &RestCreate{Name: name, Title: title, Type: ALBUM}
	b, err := json.Marshal(c)
	if err != nil {
		return
		//		log.Panicln("Error marshalling Rest create: ", jErr)
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
		return
		//log.Panic("Error connecting to: "+parentUrl+" Error: ", err)
	}

	rspValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
		//log.Panic("Error reading response: ", err)
	}
	response.Body.Close()
	itemUrl = getUrl(rspValue)
	status = response.StatusCode
	return
}

//Uploads an image to the passed in album ur
// sections of this should probabaly be scrapped and replaced
// with the go mime/multipart package
// a lot of this is hardcoded which could cause issues in the future
func (gClient *Client) UploadImage(title string, imagePath string,
	parentUrl string) (url string, status int, err error) {
	gClient.checkClient()
	hClient := new(http.Client)

	_, name := path.Split(imagePath)
	c := &RestCreate{Name: name, Title: title, Type: PHOTO}
	entity, err := json.Marshal(c)
	if err != nil {
		return
		//log.Panicln("Error marshalling Rest create: ", jErr)
	}

	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return url, status, err
		//	log.Panic("Error reading the image file: ", fErr)
	}

	var dataParts = make([]string, 13)
	boundry := "roPK9J3DoG4ZWP6etiDuJ97h-zeNAph"

	//build multipart request
	dataParts[0] = "--" + boundry
	dataParts[1] = `Content-Disposition: form-data; name="entity"`
	dataParts[2] = "Content-Type: text/plain; charset=UTF-8"
	dataParts[3] = "Content-Transfer-Encoding: 8bit"
	//space in 4
	dataParts[5] = string(entity)
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
	req, err := http.NewRequest("POST", parentUrl, buffer)
	if err != nil {
		return
	}
	req.Header.Set("X-Gallery-Request-Method", "POST")
	req.Header.Set("X-Gallery-Request-Key", gClient.APIKey)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundry)
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	//fmt.Println("request: ", req)
	response, err := hClient.Do(req)
	if err != nil {
		return
		//	log.Panic("Error connecting to: "+parentUrl+" Error: ", err)
	}

	rspValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	response.Body.Close()
	url = getUrl(rspValue)
	status = response.StatusCode
	err = nil
	return
}

//Returns the proper mime type for the passed in file
func getContentType(file string) string {
	ext := path.Ext(file)
	mType := mime.TypeByExtension(ext)
	if mType == "" {
		return "application/octet-stream"
	}
	return mType
}

//Pull URL out of rest response
func getUrl(data []byte) string {
	response := new(Response)
	json.Unmarshal(data, &response)
	return response.Url
}
