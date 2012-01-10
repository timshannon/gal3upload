package gal3rest

import (
	"fmt"
	"testing"
)

func TestRest(t *testing.T) {
	response, unmarshalled := TestConnection("http://www.timandmeg.net/gallery3/index.php")
	fmt.Println("Members:")
	fmt.Println(response.Members)
	fmt.Println("Unmarshalled:")
	fmt.Println(unmarshalled)
}
