package form

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"

type MonitorItemParam struct {
	BizIdList    []string `form:"bizIdList"`
	ProductBizId string   `form:"productBizId"`
	OsType       string   `form:"osType"`
	Display      string   `form:"display"`
	EventEum     enum.EventEum
}
