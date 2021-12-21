package dao

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"strconv"
	"text/template"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (d *MonitorItemDao) SelectMonitorItemsById(productId string, osType string) []models.MonitorItem {
	var productList []models.MonitorItem
	global.DB.Where("status = ?", "1").Where("is_display = ?", "1").Where("product_id = ?", productId).Find(&productList)
	if tools.IsBlank(osType) {
		return productList
	}
	var newProductList []models.MonitorItem
	for _, v := range productList {
		if tools.IsNotBlank(v.ShowExpression) && !isShow(v.ShowExpression, osType) {
			continue
		}
		newProductList = append(newProductList, v)
	}
	return newProductList
}

func (d *MonitorItemDao) GetMonitorItemByName(name string) models.MonitorItem {
	var model = models.MonitorItem{}
	global.DB.Where("metric_name = ?", name).First(&model)
	return model
}

func isShow(exp string, os string) bool {
	m := map[string]string{"OSTYPE": os}
	var buf bytes.Buffer
	temp, _ := template.New("exp").Parse(exp)
	if err := temp.Execute(&buf, m); err != nil {
		logger.Logger().Error(err)
	}
	isShowBool, err := strconv.ParseBool(buf.String())
	if err != nil {
		logger.Logger().Errorf("展示表达式解析失败：%v", err)
		return true
	}
	return isShowBool
}