package gal3rest

import (
	"fmt"
	"testing"
)

func TestRest(t *testing.T) {
	client := new(Client)
	client.Url = "http://www.timandmeg.net/gallery3/index.php"
	client.APIKey = "79daf60695177e16ff2480f8338b5fcc"
	album := client.GetAlbum(1)
	fmt.Println("Number of Photos: ", len(album.Photos))
	fmt.Println("Number of Albums: ", len(album.Albums))
	fmt.Println(album)
	//for a := range album.Albums {
	//fmt.Println(album.Albums[a])
	//}
}
