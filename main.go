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

package main

import (
	"code.google.com/p/gal3upload/gal3rest"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//cmd line flags
var url string
var apiKey string
var list bool
var parent string
var pid int
var gRecurse bool
var create string
var folder bool
var rebuildCache bool

//globals
var client gal3rest.Client
var cachedData []*CacheData

func init() {
	//setup command line flags
	flag.StringVar(&url, "u", "", "url of the gallery")
	flag.StringVar(&apiKey, "a", "", "API Key of the user of the gallery")
	flag.BoolVar(&list, "l", false, "List the contents of the gallery")
	flag.StringVar(&parent, "p", "", "Set the parent gallery")
	flag.IntVar(&pid, "pid", 1, "Set the parent gallery id")
	flag.BoolVar(&gRecurse, "r", false, "Recurse the gallery or file system's sub folders")
	flag.StringVar(&create, "c", "", "Creates a gallery with the given name")
	flag.BoolVar(&folder, "f", false, "Creates a local folder structure based on the gallery")
	flag.BoolVar(&rebuildCache, "rebuild", false, "Forces a rebuild of the local cache file")

	flag.Parse()
}

type CacheData struct {
	Url       string
	Name      string
	ParentUrl string
}

func main() {
	//check required flags
	if url == "" {
		fmt.Println("No URL specified with -u")
		return
	}
	if apiKey == "" {
		fmt.Println("No API key specified with -a")
		return
	}
	client = gal3rest.NewClient(url, apiKey)
	if rebuildCache {
		BuildCache()
	}

	LoadCache()
	switch {
	case list:
		List()
		return
	case create != "":
		Create()
		return
	default:
		Upload()
	}
}

//List an album's contents in the following format
// Album Name [rest id]
//  children will be one tab in from their parents
func List() {
	parentName := SetParent()
	fmt.Println("Listing Albums as Album Name [id]")
	fmt.Println()
	fmt.Println(parentName + " [" + GetId(parent) + "]")

	PrintAlbum(parent, "", gRecurse)
}

//GetId returns the REST id for the given url
func GetId(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1:][0]
}

func PrintAlbum(url string, tabs string, recurse bool) {
	tabs += "    "

	for i := range cachedData {
		if cachedData[i].ParentUrl == url {
			fmt.Println(tabs + cachedData[i].Name + " [" + GetId(cachedData[i].Url) + "]")
			if recurse {
				PrintAlbum(cachedData[i].Url, tabs, recurse)
			}
		}
	}

}

//Create an album
func Create() {
	fmt.Println("create")
}

//Upload and image
func Upload() {
	fmt.Println("upload")
}

//SetParent replaces the passed in parent id or name with the parent url
// returns the full name
func SetParent() (name string) {
	if parent != "" {
		//lookup parent url by name
		// simple loop for now, may sort data later if need be
		for i := range cachedData {
			if strings.Contains(cachedData[i].Name, parent) {
				parent = cachedData[i].Url
				name = cachedData[i].Name
			}
		}
	} else {
		parent = client.GetUrlFromId(pid)
		data := GetAlbum(parent)
		name = data.Entity.Name
	}
	return
}

// BuildCache builds a local file that caches an album's name
// id and parent id 
func BuildCache() {
	fmt.Println("Building cache from REST data")
	cachedData = RecurseAlbums(client.GetUrlFromId(1), "")

	//write out cachedData
	data, err := json.Marshal(cachedData)
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(".cache", data, 0644)
	if err != nil {
		panic(err.Error())
	}
}

func GetAlbum(url string) (data *gal3rest.RestData) {
	params := map[string]string{
		"type": "album",
	}
	data, status, err := client.GetRESTItem(url, params)
	if err != nil {
		panic(err.Error())
	}
	if status != 200 {
		fmt.Println("Gallery at "+url+" returned a status of: ", status)
	}
	return
}

//Recursively gathers all albums from the given URL
func RecurseAlbums(url string, parentUrl string) (returnData []*CacheData) {
	data := GetAlbum(url)
	//fmt.Println("Retrieved album: ", data.Entity.Name)
	returnData = append(returnData, &CacheData{url, data.Entity.Name, parentUrl})

	for m := range data.Members {
		returnData = append(returnData, RecurseAlbums(data.Members[m], url)...)
	}

	return
}

//LoadCache loads the local cache file into memory
func LoadCache() {
	data, err := ioutil.ReadFile(".cache")
	if !os.IsNotExist(err) {
		if err != nil {
			panic(string(err.Error()))
		}
		err = json.Unmarshal(data, &cachedData)
		if err != nil {
			panic(string(err.Error()))
		}
	} else {
		BuildCache()
	}

}
