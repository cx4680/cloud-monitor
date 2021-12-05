package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"gorm.io/gorm"
)

type TransactionActuatorStage struct {
}

func (s *TransactionActuatorStage) Exec(c *context.Context) error {
	return translate.InitTrans("zh")
}

type TaskActuatorStage struct {
}

func (ta *TaskActuatorStage) Exec(c *context.Context) error {
	bt := task.NewBusinessTaskImpl()
	if err := bt.Add(task.BusinessTaskDTO{
		Cron: "0 0 0/1 * * ?",
		Name: "clearAlertRecordJob",
		Task: task.Clear,
	}); err != nil {
		return err
	}

	bt.Start()
	return nil
}

type MQActuatorStage struct {
}

func (ma *MQActuatorStage) Exec(c *context.Context) error {
	return sysRocketMq.StartConsumersScribe("cloud-monitor-center", []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.InstanceTopic,
		Handler: consumer.InstanceHandler,
	}, {
		Topic:   sysRocketMq.SmsMarginReminderTopic,
		Handler: consumer.SmsMarginReminderConsumer,
	}, {
		Topic:   sysRocketMq.DeleteInstanceTopic,
		Handler: consumer.DeleteInstanceHandler,
	}})
}

type WebActuatorStage struct {
}

func (wa *WebActuatorStage) Exec(c *context.Context) error {
	return web.Start(config.GetServeConfig())
}

type ProjectInitializerFetch struct {
}

func (p *ProjectInitializerFetch) Fetch(db *gorm.DB) ([]interface{}, []string) {
	var tables []interface{}
	var sqls []string

	tables = append(tables, &models.UserPrometheusID{})
	// TODO

	return tables, sqls
}

type DBInitActuatorStage struct {
}

func (d *DBInitActuatorStage) Exec(c *context.Context) error {
	initializer := sysDb.DBInitializer{
		DB:      global.DB,
		Fetches: []sysDb.InitializerFetch{new(sysDb.CommonInitializerFetch), new(ProjectInitializerFetch)},
	}

	if err := initializer.Initnitialization(); err != nil {
		return err
	}
	return nil
}
