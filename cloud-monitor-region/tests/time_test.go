package tests

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	timestr := "2020-12-30T05:04:49.455639638Z"
	parse, _ := time.Parse("2006-01-02T15:04:05Z", timestr)
	fmt.Println(parse)
	format := time.Unix(parse.Unix(), 0).Format("2006-01-02 15:04:05")

	fmt.Println(format)
}
