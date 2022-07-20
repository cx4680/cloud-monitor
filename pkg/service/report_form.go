package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"fmt"
	"strconv"
	"strings"
)

type ReportFormService struct {
}

func NewReportFormService() *ReportFormService {
	return &ReportFormService{}
}

func (s *ReportFormService) GetMonitorData(param form.ReportFormParam) ([]form.ReportForm, error) {
	var reportForm []form.ReportForm
	var instanceList []string
	instanceMap := make(map[string]form.InstanceForm)
	for _, v := range param.InstanceList {
		instanceList = append(instanceList, v.InstanceId)
		instanceMap[v.InstanceId] = v
	}
	instances := strings.Join(instanceList, "|")
	for _, v := range param.ItemList {
		item := dao.MonitorItem.GetMonitorItemCacheByName(v)
		pql := strings.ReplaceAll(item.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
		result := QueryRange(pql, strconv.Itoa(param.Start), strconv.Itoa(param.End), strconv.Itoa(param.Step)).Data.Result
		if len(result) == 0 {
			continue
		}
		labels := strings.Split(item.Labels, ",")
		for _, v1 := range labels {
			if v1 != "instance" {

			}
		}
		if len(param.Statistics) == 0 {
			for _, v1 := range result {
				for _, v2 := range v1.Values {
					f := form.ReportForm{
						Region:       param.RegionCode,
						InstanceName: instanceMap[v1.Metric["instance"]].InstanceName,
						InstanceId:   v1.Metric["instance"],
						Status:       instanceMap[v1.Metric["instance"]].Status,
						ItemName:     item.Name,
						Time:         util.TimestampToFullTimeFmtStr(int64(v2[0].(float64))),
						Timestamp:    int64(v2[0].(float64)),
						Value:        changeDecimal(v2[1].(string)),
					}
					for _, v3 := range labels {
						if v3 != "instance" && strutil.IsNotBlank(v1.Metric[v3]) {
							f.InstanceId = f.InstanceId + " - " + v1.Metric[v3]
						}
					}
					reportForm = append(reportForm, f)
				}
			}
		} else {
			maxResult := make(map[string]form.PrometheusResult)
			minResult := make(map[string]form.PrometheusResult)
			avgResult := make(map[string]form.PrometheusResult)
			start := param.Start + (86400 - (param.Start-57600)%86400)
			end := param.End + (86400 - (param.End-57600)%86400)
			for _, v1 := range param.Statistics {
				expr := fmt.Sprintf("%s_over_time((%s)[1d:1m])", v1, pql)
				result = QueryRange(expr, strconv.Itoa(start), strconv.Itoa(end), "86400").Data.Result
				if v1 == "max" {
					for i2, v2 := range result {
						maxResult[v2.Metric["instance"]+strconv.Itoa(i2)] = v2
					}
				}
				if v1 == "min" {
					for i2, v2 := range result {
						minResult[v2.Metric["instance"]+strconv.Itoa(i2)] = v2
					}
				}
				if v1 == "avg" {
					for i2, v2 := range result {
						avgResult[v2.Metric["instance"]+strconv.Itoa(i2)] = v2
					}
				}
			}
			for i1, v1 := range result {
				for i2, v2 := range v1.Values {
					f := form.ReportForm{
						Region:       param.RegionCode,
						InstanceName: instanceMap[v1.Metric["instance"]].InstanceName,
						InstanceId:   v1.Metric["instance"],
						Status:       instanceMap[v1.Metric["instance"]].Status,
						ItemName:     item.Name,
						Time:         util.TimestampToDayTimeFmtStr(int64(v2[0].(float64)) - 1),
						Timestamp:    int64(v2[0].(float64) - 1),
					}
					if len(maxResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values) != 0 {
						f.MaxValue = changeDecimal(maxResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values[i2][1].(string))
					}
					if len(minResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values) != 0 {
						f.MinValue = changeDecimal(minResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values[i2][1].(string))
					}
					if len(avgResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values) != 0 {
						f.AvgValue = changeDecimal(avgResult[v1.Metric["instance"]+strconv.Itoa(i1)].Values[i2][1].(string))
					}
					for _, v3 := range labels {
						if v3 != "instance" && strutil.IsNotBlank(v1.Metric[v3]) {
							f.InstanceId = f.InstanceId + " - " + v1.Metric[v3]
						}
					}
					reportForm = append(reportForm, f)
				}
			}
		}
	}
	return reportForm, nil
}

func (s *ReportFormService) ExportMonitorData(param form.ReportFormParam, userInfo string) error {
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.Export
	asyncParams := []form.AsyncExportParam{
		{
			SheetSeq:   0,
			SheetName:  "云监控",
			SheetParam: jsonutil.ToString(param),
		},
	}
	asyncRequest := form.AsyncExportRequest{
		TemplateId: "cloud_monitor",
		Params:     asyncParams,
	}
	header := map[string]string{"user-info": userInfo}
	result, err := httputil.HttpPostJson(url, asyncRequest, header)
	logger.Logger().Infof("AsyncExport：%v", result)
	if err != nil {
		logger.Logger().Infof("AsyncExportError：%v", err)
		return errors.NewBusinessError("异步导出API调用失败")
	}
	return nil
}

func (s *ReportFormService) GetAlarmRecord(param form.AlarmRecordPageQueryForm) ([]form.AlarmRecord, error) {
	param.PageNum = 1
	param.PageSize = 10000
	page := dao2.AlarmRecord.GetPageList(global.DB, param.TenantID, param)
	var list []form.AlarmRecord
	if page.Records != nil {
		for _, v := range page.Records.([]vo.AlarmRecordPageVO) {
			list = append(list, form.AlarmRecord{
				AlarmId:     v.BizId,
				AlarmTime:   v.CreateTime.Format(util.FullTimeFmt),
				MonitorType: v.MonitorType,
				SourceType:  v.SourceType,
				SourceId:    v.SourceId,
				RuleName:    v.RuleName,
				Expression:  v.Expression,
				Status:      statusMap[v.Status],
				Level:       levelMap[v.Level],
			})
		}
	}
	return list, nil
}

func (s *ReportFormService) ExportAlarmRecord(param form.AlarmRecordPageQueryForm, userInfo string) error {
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.Export
	asyncParams := []form.AsyncExportParam{
		{
			SheetSeq:   0,
			SheetName:  "云监控",
			SheetParam: jsonutil.ToString(param),
		},
	}
	asyncRequest := form.AsyncExportRequest{
		TemplateId: "cloud_monitor_alarm_record",
		Params:     asyncParams,
	}
	header := map[string]string{"user-info": userInfo}
	result, err := httputil.HttpPostJson(url, asyncRequest, header)
	logger.Logger().Infof("AsyncExport：%v", result)
	if err != nil {
		logger.Logger().Infof("AsyncExportError：%v", err)
		return errors.NewBusinessError("异步导出API调用失败")
	}
	return nil
}

var statusMap = map[string]string{"firing": "告警触发", "resolved": "告警恢复"}
var levelMap = map[int]string{1: "紧急", 2: "重要", 3: "次要", 4: "提醒"}
