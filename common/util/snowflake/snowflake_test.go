package snowflake

import (
	"fmt"
	"os"
	"testing"
)

func TestWorkerId(t *testing.T) {
	os.Setenv("POD_NAME", " cloud-monitor-86584c66d4-44lr5")
	id := GetWorker().NextId()
	fmt.Println(id)
	nextId := GetWorker().NextId()
	fmt.Println(nextId)
}
