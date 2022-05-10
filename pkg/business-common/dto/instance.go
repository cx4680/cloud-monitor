package dto

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
)

type Instance struct {
	TenantId string
	List     []*model.AlarmInstance
}
