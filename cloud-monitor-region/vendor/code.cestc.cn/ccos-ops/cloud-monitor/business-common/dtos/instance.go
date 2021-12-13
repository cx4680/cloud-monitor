package dtos

import "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"

type Instance struct {
	TenantId string
	List     []*models.AlarmInstance
}
