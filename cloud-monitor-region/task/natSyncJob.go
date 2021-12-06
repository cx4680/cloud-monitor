package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/nat"
	"log"
)

type NatSyncJob struct {
}

func NewNatJob() *EcsSyncJob {
	return &EcsSyncJob{}
}
func (natSyncJob *NatSyncJob) SyncJob() {
	log.Println("nat instanceJob start")
	SyncUpdate(natSyncJob, "6", true)
	log.Println("nat instanceJob end")
}
func (natSyncJob *NatSyncJob) GetInstanceList(tenantId string) ([]*models.AlarmInstance, error) {
	pageForm := &nat.QueryParam{}
	pageVO, err := nat.GetNatInstancePage(pageForm,pageIndex, pageSize, tenantId)
	if err != nil {
		return nil, err
	}
	var infos = make([]*models.AlarmInstance, pageVO.Total)
	for i, row := range pageVO.Rows {
		infos[i] = &models.AlarmInstance{InstanceID: row.Uid, InstanceName: row.Name}
	}
	return infos, nil
}

