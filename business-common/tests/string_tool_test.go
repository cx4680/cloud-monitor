package tests

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"fmt"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	var s string = " "
	fmt.Println(strutil.IsNotEmpty(s))
	fmt.Println(strutil.IsNotBlank(s))
}
