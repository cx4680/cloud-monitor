package dao

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
	"text/template"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (d *MonitorItemDao) SelectMonitorItemsById(productId string, osType string) []model.MonitorItem {
	var monitorItemList []model.MonitorItem
	global.DB.Where("status = ?", "1").Where("is_display = ?", "1").Where("product_id = ?", productId).Find(&monitorItemList)
	if strutil.IsBlank(osType) {
		return monitorItemList
	}
	var newMonitorItemList []model.MonitorItem
	for _, v := range monitorItemList {
		if strutil.IsNotBlank(v.ShowExpression) && !isShow(v.ShowExpression, osType) {
			continue
		}
		newMonitorItemList = append(newMonitorItemList, v)
	}
	return newMonitorItemList
}

func (d *MonitorItemDao) GetMonitorItemByName(name string) model.MonitorItem {
	var model = model.MonitorItem{}
	global.DB.Where("metric_name = ?", name).First(&model)
	return model
}

func isShow(exp string, os string) bool {
	m := map[string]string{"OSTYPE": os}
	var buf bytes.Buffer
	temp, _ := template.New("exp").Parse(exp)
	if err := temp.Execute(&buf, m); err != nil {
		logger.Logger().Errorf("展示表达式解析失败：%v", err)
		return true
	}
	isShowBool, err := strconv.ParseBool(buf.String())
	if err != nil {
		logger.Logger().Errorf("展示表达式解析失败：%v", err)
		return true
	}
	return isShowBool
}
