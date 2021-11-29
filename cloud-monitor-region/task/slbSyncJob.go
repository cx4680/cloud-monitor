package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/slb"
	"github.com/robfig/cron"
	"log"
)

func CronSlbInstanceJob() {
	c := cron.New()
	job := newSlbJob()
	err := c.AddFunc("0 0 0/1 * * ?", job.syncJob)
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

type SlbSyncJob struct {
}

func newSlbJob() *SlbSyncJob {
	return &SlbSyncJob{}
}
func (slbJob *SlbSyncJob) syncJob() {
	log.Println("slbInstanceJob start")
	SyncUpdate(slbJob, "3", true)
	log.Println("slbInstanceJob end")
}

func (slbJob *SlbSyncJob) GetInstanceList(tenantId string) (interface{}, error) {
	pageVO, err := slb.GetSlbInstancePage(nil, pageIndex, pageSize, tenantId)
	if err != nil {
		return nil, err
	}
	var infos []*ProjectInstanceInfo
	for i, row := range pageVO.Rows {
		info := ProjectInstanceInfo{InstanceId: row.LbUid, InstanceName: row.Name}
		infos[i] = &info
	}
	return infos, nil
}
