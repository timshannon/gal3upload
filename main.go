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
	"bitbucket.org/tshannon/gal3upload/gal3rest"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
var workingDir string
var maxThreads = 2
var skipCache bool
var connectionFile string
var verbose bool
var deleteOnComplete bool

//globals
var client gal3rest.Client
var cachedData []*CacheData

const cRootName string = "Gallery3"

//valid upload file types
var fileTypes = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
}

type CacheData struct {
	Url       string
	Name      string
	ParentUrl string
}

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
	flag.BoolVar(&verbose, "v", false, "Verbose - Prints each upload and album creation")
	flag.BoolVar(&deleteOnComplete, "d", false, "Deletes images once they've been uploaded successfully")
	flag.BoolVar(&rebuildCache, "rebuild", false, "Forces a rebuild of the local cache file")
	flag.StringVar(&workingDir, "wd", "", "Sets the working directory of the uploader")
	flag.IntVar(&maxThreads, "t", 1, "Sets the number of threads to use for uploads.")
	flag.BoolVar(&skipCache, "skipCache", false, "Skips building the local cache file")
	flag.StringVar(&connectionFile, "connectFile", "", "File containing the API key and url of the gallery to connect to")

	flag.Parse()
}

func main() {
	//check required flags
	if connectionFile != "" {
		fileBytes, err := ioutil.ReadFile(connectionFile)
		if err != nil {
			panic(err)
		}

		buffer := bytes.NewBuffer(fileBytes)
		url, err = buffer.ReadString('\n')
		apiKey, err = buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		url = strings.TrimSpace(url)
		apiKey = strings.TrimSpace(apiKey)

	}
	if url == "" {
		fmt.Println("No URL specified with -u")

		return
	}
	if apiKey == "" {
		fmt.Println("No API key specified with -a")
		return
	}
	client = gal3rest.NewClient(url, apiKey)

	if maxThreads < 1 {
		maxThreads = 1
	}

	if rebuildCache {
		BuildCache()
	}
	if workingDir != "" {
		os.Chdir(workingDir)
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
		_ = Create(create, parentUrl)
		return
	case folder:
		CreateFolders(parentUrl, parentName, gRecurse)
		return
	default:
		wd, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}
		Upload(wd, parentUrl, gRecurse)
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
	case 200, 201:
		return true
	case 403:
		fmt.Println("Authorization Failure: Check your api key and url")
		return false
	default:
		fmt.Println("An error occurred accessing the REST resource")
		fmt.Println("Status returned was ", status)
		return false
	}
	return false

}

//Create an album
func Create(name string, albumParentUrl string) (newUrl string) {
	if verbose {
		fmt.Println("Creating album " + name)
	}
	newUrl, status, err := client.CreateAlbum(name, name, albumParentUrl)
	if err != nil {
		panic(err.Error())
	}
	if CheckStatus(status) != true {
		return
	}

	fmt.Println("Album "+name+" created with an id of ", GetId(newUrl))
	cachedData = append(cachedData, &CacheData{newUrl, name, albumParentUrl})
	WriteCache()
	return

}

//Upload an image
func Upload(dir string, dirParentUrl string, recurse bool) {
	//Get Url for current dir
	//fmt.Println("In directory: ", dir)
	var found bool
	var dirUrl string
	var subDirs []string
	var complete = make(chan *CacheData)
	var uploadList []*CacheData
	//get list of previously uploaded images
	dirCache := LoadUploadCache(dir)

	//get url for current directory
	if path.Base(dir) == cRootName {
		dirUrl = client.GetUrlFromId(1)
		found = true
	} else {
		for i := range cachedData {
			if cachedData[i].Name == path.Base(dir) &&
				cachedData[i].ParentUrl == dirParentUrl {
				dirUrl = cachedData[i].Url
				found = true
			}
		}
	}
	if !found {
		//Create new album for folder
		dirUrl = Create(path.Base(dir), dirParentUrl)
	}
	//Get list of png, jpg, jpeg, and gif files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err.Error())
	}

	for f := range files {
		//fmt.Println("file: ", files[f].Name())
		if files[f].IsDir() {
			subDirs = append(subDirs, path.Join(dir, files[f].Name()))
		} else {
			var isImage bool = false
			for t := range fileTypes {
				if strings.ToLower(path.Ext(files[f].Name())) == fileTypes[t] {
					isImage = true
					break
				}
			}
			if isImage {
				urlCache, ok := dirCache[files[f].Name()]
				if !ok {
					urlCache = &CacheData{}
				}
				uploadList = append(uploadList, &CacheData{urlCache.Url,
					path.Join(dir, files[f].Name()),
					dirUrl})
			}
		}
	}
	var i int
	for u := 0; u < len(uploadList); {
		i = 0
		for ; i < maxThreads && u < len(uploadList); i++ {
			go UploadImage(uploadList[u].Name,
				uploadList[u].ParentUrl,
				uploadList[u].Url,
				complete)
			u += 1
		}
		for j := 0; j < i; j++ {
			c := <-complete
			if c.Url != "" {
				dirCache[c.Name] = c
			}
		}
	}
	//If uploader isn't cleaning up after itself, then previous uploads should be
	// cached so they don't get uploaded again.
	if !deleteOnComplete {
		WriteUploadCache(dir, dirCache)
	}
	if recurse {
		for d := range subDirs {
			Upload(subDirs[d], dirUrl, recurse)
		}
	}

}

//load local list of previous images uploaded to make sure
// the same images don't get uploaded more than once
func LoadUploadCache(dir string) (uploadCache map[string]*CacheData) {
	data, err := ioutil.ReadFile(path.Join(dir, ".uploadcache"))
	if !os.IsNotExist(err) {
		if err != nil {
			panic(string(err.Error()))
		}
		err = json.Unmarshal(data, &uploadCache)
		if err != nil {
			panic(string(err.Error()))
		}
	}
	if uploadCache == nil {
		uploadCache = make(map[string]*CacheData)
	}
	return
}

func WriteUploadCache(dir string, uploadCache map[string]*CacheData) {
	data, err := json.Marshal(uploadCache)
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(path.Join(dir, ".uploadcache"), data, 0644)
	if err != nil {
		panic(err.Error())
	}
}

//UploadImage checks if the image was previously uploaded and still exists in REST
// and if not uploads the image
func UploadImage(imagePath string, uploadUrl string, imageUrl string, complete chan *CacheData) {
	if imageUrl != "" {
		if verbose {
			fmt.Println("Checking for previously uploaded image at " + imageUrl)
		}
		_, status, err := client.GetRESTItem(imageUrl, nil)
		if err != nil {
			panic(err.Error())
		}
		if verbose {
			fmt.Println("Finished checking for previously uploaded image at " + imageUrl)
		}
		if status == 200 {
			//image was previously uploaded
			//exit without uploading
			complete <- &CacheData{}
			return
		}
	}
	_, fileName := path.Split(imagePath)
	fmt.Println("Uploading image: ", imagePath)
	newUrl, status, err := client.UploadImage(fileName, imagePath, uploadUrl)
	if err != nil {
		fmt.Println("Error uploading image "+imagePath+": ", err.Error())
	}

	if verbose {
		fmt.Println("Finished uploading image: ", imagePath)
	}
	if !CheckStatus(status) {
		//Image didn't upload properly
		complete <- &CacheData{}
	}

	if deleteOnComplete {
		if os.Remove(fileName) != nil {
			fmt.Println("Error cleaning up successfully uploaded file: ", fileName)
		}
	}
	complete <- &CacheData{newUrl, fileName, uploadUrl}
}

//SetParent replaces the passed in parent id or name with the parent url
// returns the full name
func SetParent() {
	if parentName != "" && !skipCache {
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
			if parentId == 1 {
				parentName = cRootName
			}
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
	if !skipCache {
		cachedData = RecurseAlbums(client.GetUrlFromId(1), "")

		WriteCache()
	}
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
	if verbose {
		fmt.Println("Start Get Album " + url)
	}
	params := map[string]string{
		"type": "album",
	}
	data, status, err := client.GetRESTItem(url, params)
	if err != nil {
		panic(err.Error())
	}
	if verbose {
		fmt.Println("End Get Album " + url)
	}
	good = CheckStatus(status)
	return
}

//Recursively gathers all albums from the given URL
func RecurseAlbums(url string, purl string) (returnData []*CacheData) {
	data, _ := GetAlbum(url)
	//fmt.Println("Retrieved album: ", data.Entity.Name)
	name := data.Entity.Name
	if GetId(url) == "1" {
		name = cRootName
	}
	returnData = append(returnData, &CacheData{url, name, purl})

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
func CreateFolders(url string, dirPath string, recurse bool) {
	err := os.MkdirAll(dirPath, 0774)
	if err != nil {
		panic(err.Error())
	}
	if recurse {
		for p := range cachedData {
			if cachedData[p].ParentUrl == url {
				CreateFolders(cachedData[p].Url,
					path.Join(dirPath, cachedData[p].Name),
					recurse)
			}
		}
	}
}
