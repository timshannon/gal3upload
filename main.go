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
	"strconv"
	"strings"
)

//cmd line flags
var url string
var apiKey string
var list bool
var parentName string
var parentUrl string
var parentId int
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
	flag.StringVar(&parentName, "p", "", "Set the parent gallery")
	flag.IntVar(&parentId, "pid", 1, "Set the parent gallery id")
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
	SetParent()
	if parentUrl == "" {
		if parentId == 0 {
			//parent couldn't be found by name rebuild cache and try again
			fmt.Println("Parent couldn't be found in local cache.")
			BuildCache()
			SetParent()
		}
		// if parent still isn't found, or pid was used and album not found
		// exit
		if parentUrl == "" {
			fmt.Println("Parent name or ID doesn't exist in the Gallery")
			return
		}
	}
	switch {
	case list:
		List()
		return
	case create != "":
		Create()
		return
	case folder:
		CreateFolders()
		return
	default:
		Upload()
	}
}

//List an album's contents in the following format
// Album Name [rest id]
//  children will be one tab in from their parents
func List() {
	fmt.Println("Listing Albums as Album Name [id]")
	fmt.Println()
	fmt.Println(parentName + " [" + GetId(parentUrl) + "]")

	PrintAlbum(parentUrl, "", gRecurse)
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

func CheckStatus(status int) bool {
	switch status {
	case 200:
		return true
	case 403:
		fmt.Println("Authorization Failure: Check your api key and url")
		return false
	default:
		fmt.Println("An error occurred accessing the REST resource")
		return false
	}
	return false

}

//Create an album
func Create() {
	fmt.Println("create")
	newUrl, status, err := client.CreateAlbum(create, create, parentUrl)
	if err != nil {
		panic(err.Error())
	}
	if CheckStatus(status) != true {
		return
	}

	cachedData = append(cachedData, &CacheData{newUrl, create, parentUrl})
	WriteCache()

}

//Upload and image
func Upload() {
	fmt.Println("upload")
}

//SetParent replaces the passed in parent id or name with the parent url
// returns the full name
func SetParent() (name string) {
	if parentName != "" {
		//lookup parent url by name
		// simple loop for now, may sort data later if need be
		// Try to match exactly first
		// Also check for names where spaces are replaced with -
		altName := strings.Replace(parentName, " ", "-", -1)
		parentId = 0
		parentUrl = ""
		for i := range cachedData {
			if cachedData[i].Name == parentName ||
				cachedData[i].Name == altName {
				parentUrl = cachedData[i].Url
				parentName = cachedData[i].Name
				parentId, _ = strconv.Atoi(GetId(parentUrl))
				break
			}
			if strings.Contains(cachedData[i].Name, parentName) {
				parentUrl = cachedData[i].Url
				parentName = cachedData[i].Name
				parentId, _ = strconv.Atoi(GetId(parentUrl))
				break
			}
		}
	} else {
		//pid used
		parentUrl = client.GetUrlFromId(parentId)
		data, good := GetAlbum(parentUrl)
		if good {
			parentName = data.Entity.Name
		} else {
			parentName = ""
			parentUrl = ""
		}
	}
	return
}

// BuildCache builds a local file that caches an album's name
// id and parent id 
func BuildCache() {
	fmt.Println("Building cache from REST data")
	cachedData = RecurseAlbums(client.GetUrlFromId(1), "")

	WriteCache()
}

func WriteCache() {
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

func GetAlbum(url string) (data *gal3rest.RestData, good bool) {
	params := map[string]string{
		"type": "album",
	}
	data, status, err := client.GetRESTItem(url, params)
	if err != nil {
		panic(err.Error())
	}
	good = CheckStatus(status)
	return
}

//Recursively gathers all albums from the given URL
func RecurseAlbums(url string, purl string) (returnData []*CacheData) {
	data, _ := GetAlbum(url)
	//fmt.Println("Retrieved album: ", data.Entity.Name)
	returnData = append(returnData, &CacheData{url, data.Entity.Name, purl})

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

//CreateFolder creates a local folder based on the passed in parent
// can recursively create a folder structure that represents the gallery structure
func CreateFolders() {
}
