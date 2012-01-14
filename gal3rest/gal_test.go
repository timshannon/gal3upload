package gal3rest

import (
	"fmt"
	"testing"
)

func TestRest(t *testing.T) {
	client := new(Client)
	client.Url = "http://www.timandmeg.net/gallery3/index.php"
	client.APIKey = "79daf60695177e16ff2480f8338b5fcc"

	printMembers(client, 39380)
}

func printMembers(client *Client, id int) {
	data := client.GetRESTItem(id)
	members := data.GetMembers()
	fmt.Println(members)
	for m := range members {
		fmt.Println(members[m])
	}
}
