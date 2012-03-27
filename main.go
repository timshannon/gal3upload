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
)

//cmd line flags
var url string
var apiKey string
var list bool
var parent string
var recurse bool
var create string
var folder bool
var rebuildCache bool
var client gal3rest.Client
var cachedData []CacheData

func init() {
	//setup command line flags
	flag.StringVar(&url, "u", "", "url of the gallery")
	flag.StringVar(&apiKey, "a", "", "API Key of the user of the gallery")
	flag.BoolVar(&list, "l", false, "List the contents of the gallery")
	flag.StringVar(&parent, "p", "1", "Set the parent gallery name or ID")
	flag.BoolVar(&recurse, "r", false, "Recurse the gallery or file system's sub folders")
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
	LoadCache()
	switch {
	case list:
		List()
		return
	case create != "":
		Create()
		return
	case rebuildCache:
		BuildCache()
		return
	default:
		Upload()
	}
}

//List an album's contents
func List() {
	fmt.Println("list")
	BuildCache()
}

//Create an album
func Create() {
	fmt.Println("create")
}

//Upload and image
func Upload() {
	fmt.Println("upload")
}

// BuildCache builds a local file that caches an album's name
// id and parent id 
func BuildCache() {
	cachedData = RecurseAlbums(client.GetUrlFromId(1), "")
	//write out cachedData
}

func RecurseAlbums(url string, parentUrl string) (cacheData []CacheData) {
	params := map[string]string{
		"type": "album",
	}
	data, status, err := client.GetRESTItem(parentUrl, params)
	if err != nil {
		panic(err.Error())
	}
	if status != 200 {
		fmt.Println("Gallery returned a status of: ", status)
		return
	}
	cachedData = append(cachedData, CacheData{url, data.Entity.Name, parentUrl})

	for m := range data.Members {
		cachedData = append(cachedData, RecurseAlbums(data.Members[m], url)...)
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
		//TODO: May not be able to go from interface to slice
		err = json.Unmarshal(data, &cachedData)
		if err != nil {
			panic(string(err.Error()))
		}
	}

}
