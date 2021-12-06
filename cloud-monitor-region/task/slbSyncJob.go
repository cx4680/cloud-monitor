package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/slb"
	"log"
)

type SlbSyncJob struct {
}

func NewSlbJob() *SlbSyncJob {
	return &SlbSyncJob{}
}
func (slbJob *SlbSyncJob) SyncJob() {
	log.Println("slbInstanceJob start")
	SyncUpdate(slbJob, "3", true)
	log.Println("slbInstanceJob end")
}

func (slbJob *SlbSyncJob) GetInstanceList(tenantId string) ([]*models.AlarmInstance, error) {
	pageVO, err := slb.GetSlbInstancePage(nil, pageIndex, pageSize, tenantId)
	if err != nil {
		return nil, err
	}
	var infos = make([]*models.AlarmInstance, pageVO.Total)
	for i, row := range pageVO.Rows {
		infos[i] = &models.AlarmInstance{InstanceID: row.LbUid, InstanceName: row.Name}
	}
	return infos, nil
}
