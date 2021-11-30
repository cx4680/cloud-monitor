package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/loader"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description  This is a sample server Petstore server.
// @BasePath /
func main() {
	l := &loader.SysLoaderImpl{}
	if err := l.StartServe(l); err != nil {
		os.Exit(-1)
	}
}
