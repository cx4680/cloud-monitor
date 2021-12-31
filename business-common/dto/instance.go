package dto

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"

type Instance struct {
	TenantId string
	List     []*model.AlarmInstance
}
