package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/loader"
	"os"
)

func main() {
	l := &loader.SysLoaderImpl{}
	if err := l.StartServe(l); err != nil {
		os.Exit(-1)
	}

}
