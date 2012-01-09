package gal3rest

import (
	"fmt"
	"testing"
	"io/ioutil"
)

func TestRest(t *testing.T) {
	response := TestConnection("http://www.timandmeg.net/gallery3/index.php")
	fmt.Println("content Length: ", response.ContentLength)
	b, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println(string(b))

}
