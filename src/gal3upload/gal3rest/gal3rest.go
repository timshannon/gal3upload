package gal3rest

import (
	"fmt"
	"http"
)

func TestConnection() {
	client := new(http.Client)
	response, err := client.Get("http://www.timandmeg.net/gallery3/index.php/")
	if err != nil {
		fmt.Println("error: ", err.String())
	}
	fmt.Print("Status: ", response.Status)
}
