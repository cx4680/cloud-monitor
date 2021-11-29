package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/ecs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"github.com/robfig/cron"
	"log"
)

func CronEcsInstanceJob() {
	c := cron.New()
	ecsInstanceJob := newEcsJob()
	err := c.AddFunc("0 0 0/1 * * ?", ecsInstanceJob.syncJob)
	if err != nil {
		log.Println("clearAlertRecordJob error", err)
	}
	c.Start()
	defer c.Stop()
	select {}
}

type EcsSyncJob struct {
}

func newEcsJob() *EcsSyncJob {
	return &EcsSyncJob{}
}
func (ecsSyncJob *EcsSyncJob) syncJob() {
	log.Println("ecs instanceJob start")
	SyncUpdate(ecsSyncJob, "1", true)
	log.Println("ecs instanceJob end")
}
func (ecsSyncJob *EcsSyncJob) GetInstanceList(tenantId string) (interface{}, error) {
	pageForm := &forms.EcsQueryPageForm{
		TenantId: tenantId,
		Current:  pageIndex,
		PageSize: pageSize,
	}
	pageVO, err := ecs.PageList(pageForm)
	if err != nil {
		return nil, err
	}
	return pageVO.Rows, nil
}
