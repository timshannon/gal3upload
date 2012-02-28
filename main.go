package main

import (
	"code.google.com/p/gal3upload/gal3rest"
	"flag"
	"fmt"
	"reflect"
)

//cmd line flags
var url string
var apiKey string
var test bool

func init() {
	//setup command line flags
	flag.StringVar(&url, "u", "", "url of the gallery")
	flag.StringVar(&apiKey, "a", "", "API Key of the user of the gallery")
	flag.BoolVar(&test, "t", false, "Test the connection with the given credentials")
	flag.Parse()
}
func main() {
	//check flags
	if url == "" {
		fmt.Println("No URL specified with -u")
		return
	}
	if apiKey == "" {
		fmt.Println("No API key specified with -a")
		return
	}
	if test {
		TestClient()
	}
}

func TestClient() {
	client := gal3rest.NewClient(url,
		apiKey)

	album := client.GetAlbum(client.GetUrlFromId(1))
	fmt.Println("Number of Photos: ", len(album.Photos))
	fmt.Println("Number of Albums: ", len(album.Albums))
	for i := range album.Albums {
		fmt.Println("Album: ", album.Albums[i])
	}

	fmt.Println("Creating a new album")
	client.CreateAlbum("This is my Sample Album", "Sample Album",
		client.GetUrlFromId(1))
}

func PrintEntity(entity *gal3rest.Entity) {
	ref := reflect.ValueOf(entity).Elem()
	entityType := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		fmt.Printf("%s: %s  %v\n",
			entityType.Field(i).Name,
			field.Type(), field.Interface())
	}
}
