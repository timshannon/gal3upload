package gal3rest

import (
	"fmt"
	"testing"
)

func TestRest(t *testing.T) {
	response := TestConnection("http://www.timandmeg.net/gallery3/index.php")
	fmt.Println(response)
}
