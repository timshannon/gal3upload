package gal3rest

import (
	"fmt"
	"testing"
)

func TestRest(t *testing.T) {
	client := new(Client)
	client.Url = "http://www.timandmeg.net/gallery3/index.php"
	client.APIKey = "79daf60695177e16ff2480f8338b5fcc"

	members := client.GetMembers(1)
	fmt.Println("Members: ", members)
}
