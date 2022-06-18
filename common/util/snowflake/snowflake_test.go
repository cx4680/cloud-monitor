package snowflake

import (
	"fmt"
	"os"
	"testing"
)

func TestWorkerId(t *testing.T) {
	os.Setenv("POD_NAME", "test1")
	id := GetWorker().NextId()
	fmt.Println(id)
	nextId := GetWorker().NextId()
	fmt.Println(nextId)
}
