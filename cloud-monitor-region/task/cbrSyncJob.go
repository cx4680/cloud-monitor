package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/cbr"
	"log"
	"strconv"
)

type CbrSyncJob struct {
}

func NewCbrJob() *EipSyncJob {
	return &EipSyncJob{}
}
func (cbrSyncJob *CbrSyncJob) SyncJob() {
	log.Println("cbr instanceJob start")
	SyncUpdate(cbrSyncJob, "5", true)
	log.Println("cbr instanceJob end")
}
func (cbrSyncJob *CbrSyncJob) GetInstanceList(tenantId string) ([]*models.AlarmInstance, error) {
	pageForm := &cbr.QueryParam{
		TenantId:   tenantId,
		PageNumber: strconv.Itoa(pageIndex),
		PageSize:   strconv.Itoa(pageSize),
	}
	pageVO, err := cbr.PageList(pageForm)
	if err != nil {
		return nil, err
	}
	var infos = make([]*models.AlarmInstance, pageVO.Total_count)
	for i, row := range pageVO.Data {
		infos[i] = &models.AlarmInstance{
			InstanceID:   row.VaultId,
			InstanceName: row.Name,
		}
	}
	return infos, nil
}
