package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/eip"
	"log"
)

type EipSyncJob struct {
}

func NewEipJob() *EipSyncJob {
	return &EipSyncJob{}
}
func (eipSyncJob *EipSyncJob) SyncJob() {
	log.Println("eip instanceJob start")
	SyncUpdate(eipSyncJob, "2", true)
	log.Println("eip instanceJob end")
}
func (eipSyncJob *EipSyncJob) GetInstanceList(tenantId string) ([]*models.AlarmInstance, error) {
	pageVO, err := eip.GetEipInstancePage(nil, pageIndex, pageSize, tenantId)
	if err != nil {
		return nil, err
	}
	var infos = make([]*models.AlarmInstance, pageVO.Total)
	for i, row := range pageVO.Rows {
		infos[i] = &models.AlarmInstance{InstanceID: row.InstanceUid, InstanceName: row.Name}
	}
	return infos, nil
}
