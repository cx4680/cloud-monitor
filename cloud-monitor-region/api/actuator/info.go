package actuator

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

type info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateTime  string `json:"createTime"`
	Author      string `json:"author"`
	Git         string `json:"git"`
}

func Info() string {
	return jsonutil.ToString(info{
		Name:        "云监控区域化服务",
		Description: "区域化的业务程序",
		CreateTime:  "2021-11-11",
		Author:      "Jim",
		Git:         "",
	})
}
