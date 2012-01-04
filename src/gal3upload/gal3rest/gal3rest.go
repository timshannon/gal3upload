package gal3rest

import (
	"gorest.googlecode.com/hg/gorest"
	"fmt"
)

func TestConnection() {
	rb, _ := gorest.NewRequestBuilder()
	var u string
	rb.Get(u, "http://www.timandmeg.net/gallyer3/index.php/rest/item/1")
	fmt.Println(u)
}
