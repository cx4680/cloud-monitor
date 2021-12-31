package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
)

type ConfigItemDao struct {
}

var ConfigItem = new(ConfigItemDao)

func (dao *ConfigItemDao) GetConfigItem(code interface{}, pid string, data string) *model.ConfigItem {
	item := model.ConfigItem{}
	db := global.DB
	if code != nil {
		db = db.Where("code", code)
	}
	if len(pid) > 0 {
		db = db.Where("pid", pid)
	}
	if len(data) > 0 {
		db = db.Where("data", data)
	}
	db.Find(&item)
	return &item
}

func (dao *ConfigItemDao) GetConfigItemList(pid string) []*model.ConfigItem {
	var list []*model.ConfigItem
	db := global.DB
	db = db.Where("pid", pid).Order("sort_id ASC")
	db.Find(&list)
	return list
}

const (
	StatisticalPeriodPid  = "1"  //统计周期
	ContinuousCyclePid    = "2"  //持续周期
	StatisticalMethodsPid = "3"  //统计方式
	ComparisonMethodPid   = "4"  //对比方式
	OverviewItemPid       = "21" //概览监控项
	RegionListPid         = "24" //region列表
	MonitorRange          = "5"  //监控周期
	NoticeChannel         = "33" //通知方式
	AlarmLevel            = "28" //告警级别
)
