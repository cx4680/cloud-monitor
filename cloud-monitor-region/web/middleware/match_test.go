package middleware

import (
	"log"
	"path"
	"testing"
)

func Test_pathMather(t *testing.T) {
	pattern := "/hawkeye/inner/**"
	uri := "/hawkeye/inner/ab/aa"

	match, err := path.Match(pattern, uri)
	if err != nil {
		log.Println(err)
	}
	log.Println(match)
}
