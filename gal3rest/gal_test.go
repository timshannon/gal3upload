package gal3rest

import (
	"fmt"
	"testing"
	"reflect"
)

func TestRest(t *testing.T) {
	client := new(Client)
	client.Url = "http://www.timandmeg.net/gallery3/index.php"
	client.APIKey = "79daf60695177e16ff2480f8338b5fcc"
	album := client.GetAlbum(client.GetUrlFromId(1))
	fmt.Println("Number of Photos: ", len(album.Photos))
	fmt.Println("Number of Albums: ", len(album.Albums))
	for i := range album.Albums {
		fmt.Println("Album: ", album.Albums[i])
	}
	PrintEntity(&album.Entity)
}

func PrintEntity(entity *Entity) {
	ref := reflect.ValueOf(entity).Elem()
	entityType := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		fmt.Printf("%s: %s  %v\n",
			entityType.Field(i).Name,
			field.Type(), field.Interface())
	}
}
