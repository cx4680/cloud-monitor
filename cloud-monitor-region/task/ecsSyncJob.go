package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/ecs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"log"
)

type EcsSyncJob struct {
}

func NewEcsJob() *EcsSyncJob {
	return &EcsSyncJob{}
}
func (ecsSyncJob *EcsSyncJob) SyncJob() {
	log.Println("ecs instanceJob start")
	SyncUpdate(ecsSyncJob, "1", true)
	log.Println("ecs instanceJob end")
}
func (ecsSyncJob *EcsSyncJob) GetInstanceList(tenantId string) ([]*models.AlarmInstance, error) {
	pageForm := &forms.EcsQueryPageForm{
		TenantId: tenantId,
		Current:  pageIndex,
		PageSize: pageSize,
	}
	pageVO, err := ecs.PageList(pageForm)
	if err != nil {
		return nil, err
	}
	var infos = make([]*models.AlarmInstance, pageVO.Total)
	for i, row := range pageVO.Rows {
		infos[i] = &models.AlarmInstance{InstanceID: row.InstanceId, InstanceName: row.InstanceName}
	}
	return infos, nil
}
