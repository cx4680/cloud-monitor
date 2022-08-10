package dao

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"text/template"
	"time"
)

type MonitorItemDao struct {
}

var MonitorItem = new(MonitorItemDao)

func (d *MonitorItemDao) GetMonitorItem(productBizId, osType, display string) []model.MonitorItem {
	var monitorItemList []model.MonitorItem
	if strutil.IsNotBlank(display) {
		global.DB.Where("status = ? AND is_display = ? AND product_biz_id = ? AND display LIKE ?", "1", "1", productBizId, "%"+display+"%").Find(&monitorItemList)
	} else {
		global.DB.Where("status = ? AND is_display = ? AND product_biz_id = ?", "1", "1", productBizId).Find(&monitorItemList)
	}
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

func (d *MonitorItemDao) ChangeDisplay(db *gorm.DB, productBizId, display string, bizIdList []string) {
	db.Model(&model.MonitorItem{}).Where("product_biz_id = ? AND biz_id IN (?)", productBizId, bizIdList).Update("display", display)
}

func (d *MonitorItemDao) GetMonitorItemCacheByName(name string) model.MonitorItem {
	value, err := sys_redis.Get(fmt.Sprintf(constant.MonitorItemKey, name))
	if err != nil {
		logger.Logger().Error("key=" + name + ", error:" + err.Error())
	}
	var monitorItemModel = model.MonitorItem{}
	if strutil.IsNotBlank(value) {
		jsonutil.ToObject(value, &monitorItemModel)
		return monitorItemModel
	}
	monitorItemModel = d.GetMonitorItemByName(name)
	if monitorItemModel == (model.MonitorItem{}) {
		logger.Logger().Info("获取监控项为空")
		return monitorItemModel
	}
	if e := sys_redis.SetByTimeOut(name, jsonutil.ToString(monitorItemModel), time.Hour); e != nil {
		logger.Logger().Error("设置监控项缓存错误, key=" + name)
	}
	return monitorItemModel
}

func (d *MonitorItemDao) GetMonitorItemByName(name string) model.MonitorItem {
	var monitorItemModel = model.MonitorItem{}
	global.DB.Where("metric_name = ?", name).First(&monitorItemModel)
	return monitorItemModel
}

func (d *MonitorItemDao) GetMonitorItemByMetricCode(metricCode string) form.MonitorItem {
	var monitorItem = form.MonitorItem{}
	global.DB.Raw(SelectMonitorItem, metricCode).Find(&monitorItem)
	if monitorItem == (form.MonitorItem{}) {
		logger.Logger().Info("获取监控项为空")
		return monitorItem
	}
	return monitorItem
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

var SelectMonitorItem = "SELECT mi.name AS item_name, mi.metrics_linux AS metric, mi.labels AS labels, mp.abbreviation AS product_abbreviation FROM t_monitor_item AS mi LEFT JOIN t_monitor_product mp ON mi.product_biz_id = mp.biz_id WHERE mi.metric_name = ?;"
