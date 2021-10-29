package tests

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"fmt"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	var s string = " "
	fmt.Println(tools.IsNotEmpty(s))
	fmt.Println(tools.IsNotBlank(s))
}
