package main

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums/handlerType"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"testing"
)

func TestDD(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}

	rule := []*models.AlarmInstance{{
		InstanceID:   "2",
		InstanceName: "ecs-1",
		RegionName:   "A",
	},
		{
			InstanceID:   "2",
			InstanceName: "ecs-3",
			RegionName:   "B",
		},
	}

	db.Clauses(clause.OnConflict{DoNothing: false}).Create(&rule)
}

func TestRuleAdd(t *testing.T) {
	str := "\t{\n\t\t\"alarmLevel\": 1,\n\t\t\"groupList\": [\n\t\"1\"\n\t],\n\t\"instanceList\": [\n\t{\n\t\"instanceId\": \"123\",\n\t\"zoneCode\": \"11\",\n\t\"regionCode\": \"wh\",\n\t\"regionName\": \"wuhan\",\n\t\"zoneName\": \"a\",\n\t\"ip\": \"192.10.11.123\",\n\t\"instanceName\": \"实例名称111\"\n\t}\n\t],\n\t\"monitorType\": \"云产品监控\",\n\t\"noticeChannel\": \"all\",\n\t\"productType\": \"云服务器ECS\",\n\t\"ruleName\": \"string\",\n\t\"scope\": \"ALL\",\n\t\"silencesTime\": \"3小时\",\n\t\"triggerCondition\": {\n\t\"comparisonOperator\": \"greater\",\n\t\"metricName\": \"ecs_cpu_base_usage\",\n\t\"period\": 10,\n\t\"statistics\": \"Maximum\",\n\t\"threshold\": 10,\n\t\"times\": 10\n\t}\n}"
	dto := &forms.AlarmRuleAddReqDTO{}
	tools.ToObject(str, dto)
	dto.AlarmHandlerList = []*forms.Handler{
		{
			HandleType: handlerType.Email,
		}, {
			HandleType:   handlerType.Http,
			HandleParams: "http://127.0.0.1:9876/1",
		},
	}
	dto.ResourceGroupList = []*forms.ResGroupInfo{
		{
			ResGroupName: "fex-1",
			ResourceList: []*forms.InstanceInfo{
				{
					InstanceId:   "ecs-1",
					InstanceName: "CES-1",
				},
			},
		},
	}
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}
	dto.TenantId = "1"
	dao.AlarmRule.SaveRule(db, dto)
}

func TestRuleGet(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}
	detail := dao.AlarmRule.GetDetail(db, "920729272954388480", "1")
	fmt.Printf("%+v", detail)
}

func TestRuleUpdate(t *testing.T) {
	str := "\t{\n\t\t\"alarmLevel\": 2,\n\t\t\"groupList\": [\n\t\"1\"\n\t],\n\t\"instanceList\": [\n\t{\n\t\"instanceId\": \"2222222\",\n\t\"zoneCode\": \"11\",\n\t\"regionCode\": \"wh\",\n\t\"regionName\": \"wuhan\",\n\t\"zoneName\": \"a\",\n\t\"ip\": \"192.10.11.123\",\n\t\"instanceName\": \"实例名称111\"\n\t}\n\t],\n\t\"monitorType\": \"云产品监控\",\n\t\"noticeChannel\": \"all\",\n\t\"productType\": \"云服务器ECS\",\n\t\"ruleName\": \"string\",\n\t\"scope\": \"ALL\",\n\t\"silencesTime\": \"3小时\",\n\t\"triggerCondition\": {\n\t\"comparisonOperator\": \"greater\",\n\t\"metricName\": \"ecs_cpu_base_usage\",\n\t\"period\": 10,\n\t\"statistics\": \"Maximum\",\n\t\"threshold\": 10,\n\t\"times\": 10\n\t}\n}"
	dto := &forms.AlarmRuleAddReqDTO{}
	tools.ToObject(str, dto)
	dto.AlarmHandlerList = []*forms.Handler{
		{
			HandleType: handlerType.Email,
		}, {
			HandleType:   handlerType.Http,
			HandleParams: "http://127.0.0.1:9876/1",
		},
	}
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}
	dto.Id = "920729272954388480"
	dto.TenantId = "1"
	dao.AlarmRule.UpdateRule(db, dto)
}

func TestDeleteInstances(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}
	var i int
	db.Raw("delete FROM t_alarm_rule_resource_rel where tenant_id= ? and resource_id in (?)", "xx", []string{"1", "2", "3"}).Find(&i)

}
