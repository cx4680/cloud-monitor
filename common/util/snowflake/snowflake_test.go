package snowflake

import (
	"os"
	"testing"
)

func TestWorkerId(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	workerId := getWorkerId()
	t.Logf("workerId: %d", workerId)

	os.Setenv("POD_NAME", "test1")
	workerId = getWorkerId()
	t.Logf("workerId: %d", workerId)
}
