package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysDb"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonTask "code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/consumer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"context"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"strings"
)

type TransactionActuatorStage struct {
}

func (s *TransactionActuatorStage) Exec(c *context.Context) error {
	return translate.InitTrans("zh")
}

type TaskActuatorStage struct {
}

func (ta *TaskActuatorStage) Exec(c *context.Context) error {
	bt := commonTask.NewBusinessTaskImpl()
	if err := task.AddSyncJobs(bt); err != nil {
		return err
	}
	bt.Start()
	return nil
}

type MQActuatorStage struct {
}

func (ma *MQActuatorStage) Exec(c *context.Context) error {
	return sysRocketMq.StartConsumersScribe(sysRocketMq.Group(config.GetCommonConfig().RegionName), []*sysRocketMq.Consumer{{
		Topic:   sysRocketMq.AlertContactTopic,
		Handler: consumer.AlertContactHandler,
	}, {
		Topic:   sysRocketMq.AlertContactGroupTopic,
		Handler: consumer.AlertContactGroupHandler,
	}, {
		Topic:   sysRocketMq.RuleTopic,
		Handler: consumer.AlarmRuleHandler,
	}})
}

type K8sActuatorStage struct{}

func (k *K8sActuatorStage) Exec(c *context.Context) error {
	return k8s.InitK8s()
}

type WebActuatorStage struct {
}

func (wa *WebActuatorStage) Exec(c *context.Context) error {
	return web.Start(config.GetServeConfig())
}

type ProjectInitializerFetch struct {
}

func (p *ProjectInitializerFetch) Fetch(db *gorm.DB) ([]interface{}, []string, error) {
	var tables []interface{}
	var sqls []string

	tables = append(tables, &models.UserPrometheusID{})

	//加载SQL
	sqlBytes, err := ioutil.ReadFile("script/region.sql")
	if err != nil {
		log.Println("load sql file error", err)
		return nil, nil, err
	}
	sql := string(sqlBytes)
	if tools.IsNotBlank(sql) {
		sqls = append(sqls, strings.Split(sql, ";")...)
	}

	return tables, sqls, nil
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

type IamActuatorStage struct {
}

func (i *IamActuatorStage) Exec(c *context.Context) error {
	cfg := config.GetIamConfig()
	middleware.InitIamConfig(cfg.Site, cfg.Region, cfg.Log)
	return nil
}
