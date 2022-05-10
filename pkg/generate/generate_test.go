package generate

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	generate("t_alarm_item", "alarm_item.go")
}
