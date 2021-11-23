package tests

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestAlertContactInsert(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/hawkeye?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	global.DB = db

	sysRocketMq.InitProducer()

	groupService := service.NewAlertContactGroupService()
	informationService := service.NewAlertContactInformationService()
	contactService := service.NewAlertContactService(groupService, informationService)

	param := forms.AlertContactParam{
		TenantId:    "123",
		ContactId:   "",
		ContactName: "jim",
		Phone:       "18521084140",
		Email:       "jim@126.com",
		Lanxin:      "",
		CreateUser:  "1",
		Description: "aaaaaa",
		ActiveCode:  "",
		PageCurrent: 0,
		PageSize:    0,
		GroupIdList: []string{"1", "2"},
	}

	//本地持久化+发送远程
	contactService.Persistence(contactService, config.GetRocketmqConfig().AlertContactTopic, param)
}
