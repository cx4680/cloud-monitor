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
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ReportFormService struct {
	prometheus *PrometheusService
}

func NewReportFormService() *ReportFormService {
	return &ReportFormService{
		prometheus: NewPrometheusService(),
	}
}

func (s *ReportFormService) GetMonitorData(param form.ReportFormParam) ([]*form.ReportForm, error) {
	var instanceList []string
	instanceMap := make(map[string]*form.InstanceForm)
	for _, v := range param.InstanceList {
		instanceList = append(instanceList, v.InstanceId)
		instanceMap[v.InstanceId] = v
	}
	instances := strings.Join(instanceList, "|")
	item := dao.MonitorItem.GetMonitorItemCacheByName(param.ItemList[0])
	labels := strings.Split(item.Labels, ",")
	fmt.Println("导出dns：", item.MetricName)
	fmt.Println("param：", param.Name)
	pql := strings.ReplaceAll(item.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
	if item.MetricName == "private_dns_dns_requests_total" || item.MetricName == "private_dns_dns_requests_total_rate1m" {
		pql = strings.ReplaceAll(item.MetricsLinux, constant.MetricLabel, "instanceId=~'"+instances+"'")
	}
	fmt.Println("pql:", pql)
	//获取单个指标的所有实例数据
	reportFormList := s.getOneItemData(param, item, instanceMap, pql, labels)
	return reportFormList, nil
}

func (s *ReportFormService) getOneItemData(param form.ReportFormParam, item model.MonitorItem, instanceMap map[string]*form.InstanceForm, pql string, labels []string) []*form.ReportForm {
	if len(param.Statistics) == 0 {
		return s.getOriginData(param, item, instanceMap, pql, labels)
	}
	return s.getAggregationData(param, item, instanceMap, pql, labels)

}

func (s *ReportFormService) getOriginData(param form.ReportFormParam, item model.MonitorItem, instanceMap map[string]*form.InstanceForm, pql string, labels []string) []*form.ReportForm {
	result := s.prometheus.QueryRange(pql, strconv.Itoa(param.Start), strconv.Itoa(param.End), strconv.Itoa(param.Step)).Data.Result
	if len(result) == 0 {
		return nil
	}
	fmt.Println("orgin-pql:", pql)
	var list []*form.ReportForm
	for _, prometheusResult := range result {
		for _, prometheusValue := range prometheusResult.Values {
			if f := s.buildOriginReportForm(param, instanceMap, item, labels, prometheusResult, prometheusValue); f != nil {
				list = append(list, f)
			}
		}
	}
	return list
}

func (s *ReportFormService) getAggregationData(param form.ReportFormParam, item model.MonitorItem, instanceMap map[string]*form.InstanceForm, pql string, labels []string) []*form.ReportForm {
	//计算开始时间当天的23时59分59秒
	start := param.Start + (86400 - (param.Start-57600)%86400)
	//计算结束时间当天的23时59分59秒
	end := param.End + (86400 - (param.End-57600)%86400)
	var result map[string]*form.PrometheusResult
	var ret = make(map[string]map[string]*form.PrometheusResult)
	//开启协程
	group := &sync.WaitGroup{}
	group.Add(len(param.Statistics))
	for _, statistics := range param.Statistics {
		m := make(map[string]*form.PrometheusResult)
		go s.getStatisticsMap(statistics, pql, start, end, param.Statistics, labels, group, m)
		ret[statistics] = m
		if result == nil {
			result = m
		}
	}
	group.Wait()
	var list []*form.ReportForm
	for k, v := range result {
		for i := range v.Values {
			dataMap := make(map[string][]interface{})
			for calcStyle, d := range ret {
				dataMap[calcStyle] = d[k].Values[i]
			}
			if item.MetricName == "private_dns_dns_requests_total" || item.MetricName == "private_dns_dns_requests_total_rate1m" {
				if f := s.buildAggregationReportForm(v.Metric["instanceId"], k, item.Name, instanceMap, dataMap); f != nil {
					list = append(list, f)
				}
			} else {
				if f := s.buildAggregationReportForm(v.Metric["instance"], k, item.Name, instanceMap, dataMap); f != nil {
					list = append(list, f)
				}
			}
		}
	}
	return list
}

func (s *ReportFormService) buildOriginReportForm(param form.ReportFormParam, instanceMap map[string]*form.InstanceForm, item model.MonitorItem, labels []string, prometheusResult *form.PrometheusResult, prometheusValue []interface{}) (f *form.ReportForm) {
	defer func() {
		if e := recover(); e != nil {
			logger.Logger().Error(e)
		}
	}()
	fmt.Printf("prometheusResult:%v", prometheusResult.Metric)
	if item.MetricName == "private_dns_dns_requests_total" || item.MetricName == "private_dns_dns_requests_total_rate1m" {
		f = &form.ReportForm{
			Region:       param.RegionCode,
			InstanceName: instanceMap[prometheusResult.Metric["instanceId"]].InstanceName,
			InstanceId:   prometheusResult.Metric["instanceId"],
			Status:       instanceMap[prometheusResult.Metric["instanceId"]].Status,
			ItemName:     item.Name,
			Time:         util.TimestampToFullTimeFmtStr(int64(prometheusValue[0].(float64))),
			Timestamp:    int64(prometheusValue[0].(float64)),
			Value:        changeDecimal(prometheusValue[1].(string)),
		}
	} else {
		f = &form.ReportForm{
			Region:       param.RegionCode,
			InstanceName: instanceMap[prometheusResult.Metric["instance"]].InstanceName,
			InstanceId:   prometheusResult.Metric["instance"],
			Status:       instanceMap[prometheusResult.Metric["instance"]].Status,
			ItemName:     item.Name,
			Time:         util.TimestampToFullTimeFmtStr(int64(prometheusValue[0].(float64))),
			Timestamp:    int64(prometheusValue[0].(float64)),
			Value:        changeDecimal(prometheusValue[1].(string)),
		}
	}
	for _, label := range labels {
		fmt.Println("label:", label)
		if label != "instance" && strutil.IsNotBlank(prometheusResult.Metric[label]) {
			f.InstanceId = f.InstanceId + " - " + prometheusResult.Metric[label]
		}
	}
	fmt.Println(fmt.Printf("f:%v", f))
	return
}

func (s *ReportFormService) buildAggregationReportForm(instanceId, key, itemName string, instanceMap map[string]*form.InstanceForm, dataMap map[string][]interface{}) (f *form.ReportForm) {
	defer func() {
		if e := recover(); e != nil {
			logger.Logger().Error(e)
		}
	}()
	time := ""
	timestamp := int64(0)
	for _, v := range dataMap {
		time = util.TimestampToDayTimeFmtStr(int64(v[0].(float64)) - 1)
		timestamp = int64(v[0].(float64) - 1)
		break
	}
	f = &form.ReportForm{
		Region:       config.Cfg.Common.RegionName,
		InstanceName: instanceMap[instanceId].InstanceName,
		InstanceId:   key,
		Status:       instanceMap[instanceId].Status,
		ItemName:     itemName,
		Time:         time,
		Timestamp:    timestamp,
	}

	for calcStyle, d := range dataMap {
		rf := reflect.ValueOf(f)
		ff := rf.Elem().FieldByName(firstUpper(calcStyle) + "Value")
		ff.SetString(changeDecimal(d[1].(string)))
	}
	return
}

func (s *ReportFormService) getStatisticsMap(aggregation, pql string, start, end int, statistics, labels []string, group *sync.WaitGroup, resultMap map[string]*form.PrometheusResult) {
	defer group.Done()
	for _, v := range statistics {
		if aggregation == v {
			expr := fmt.Sprintf("%s_over_time((%s)[1d:1h])", aggregation, pql)
			result := s.prometheus.QueryRange(expr, strconv.Itoa(start), strconv.Itoa(end), "86400").Data.Result
			for _, prometheusResult := range result {
				key := prometheusResult.Metric["instance"]
				for _, label := range labels {
					if label != "instance" && strutil.IsNotBlank(prometheusResult.Metric[label]) {
						key = key + " - " + prometheusResult.Metric[label]
					}
				}
				resultMap[key] = prometheusResult
			}
		}
	}
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]

}

func (s *ReportFormService) ExportMonitorData(param form.ReportFormParam, userInfo string) error {
	var sheetParamList []string
	var newParam form.ReportFormParam
	for _, item := range param.ItemList {
		for _, instance := range param.InstanceList {
			newParam = param
			newParam.ItemList = []string{item}
			newParam.InstanceList = []*form.InstanceForm{instance}
			sheetParamList = append(sheetParamList, jsonutil.ToString(newParam))
		}
	}
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.Export
	asyncParams := []form.AsyncExportParam{
		{
			SheetSeq:       0,
			SheetName:      "云监控-云产品监控",
			SheetParamList: sheetParamList,
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
	page := dao2.AlarmRecord.GetPageList(global.DB, param.TenantId, param)
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
			SheetName:  "云监控-告警历史",
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
var ecsCpuBaseUsageDownSampling = "100 * avg by(instance,instanceType)(rate(ecs_base_vcpu_seconds{$INSTANCE}[3h]))"

func (s *ReportFormService) GetReportFormData(param form.ReportFormParam) ([]*form.ReportForm, error) {
	if len(param.InstanceList) == 0 {
		return nil, errors.NewBusinessError("实例不能为空")
	}
	if len(param.ItemList) == 0 {
		return nil, errors.NewBusinessError("指标不能为空")
	}
	var instanceList []string
	instanceMap := make(map[string]*form.InstanceForm)
	for _, v := range param.InstanceList {
		instanceList = append(instanceList, v.InstanceId)
		instanceMap[v.InstanceId] = v
	}
	instances := strings.Join(instanceList, "|")
	item := dao.MonitorItem.GetMonitorItemCacheByName(param.ItemList[0])
	labels := strings.Split(item.Labels, ",")
	var downSampling = false
	if int(time.Now().Unix())-param.Start >= 3024000 {
		downSampling = true
		if item.MetricName == "ecs_cpu_base_usage" {
			item.MetricsLinux = ecsCpuBaseUsageDownSampling
		}
	}

	pql := strings.ReplaceAll(item.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
	//获取单个指标的所有实例数据
	//计算开始时间当天的23时59分59秒
	start := param.Start + (86400 - (param.Start-57600)%86400)
	//计算结束时间当天的23时59分59秒
	end := param.End + (86400 - (param.End-57600)%86400)
	var result map[string]*form.PrometheusResult
	var ret = make(map[string]map[string]*form.PrometheusResult)
	//开启协程
	group := &sync.WaitGroup{}
	group.Add(len(param.Statistics))
	for _, statistics := range param.Statistics {
		m := make(map[string]*form.PrometheusResult)
		go func(aggregation, pql string, start, end int, statistics, labels []string, group *sync.WaitGroup, resultMap map[string]*form.PrometheusResult) {
			defer group.Done()
			for _, v := range statistics {
				if aggregation == v {
					expr := fmt.Sprintf("%s_over_time((%s)[1d:1h])", aggregation, pql)
					var results []*form.PrometheusResult
					if downSampling {
						results = s.prometheus.QueryRangeDownSampling(expr, strconv.Itoa(start), strconv.Itoa(end), "86400").Data.Result
					} else {
						results = s.prometheus.QueryRange(expr, strconv.Itoa(start), strconv.Itoa(end), "86400").Data.Result
					}
					for _, prometheusResult := range results {
						key := prometheusResult.Metric["instance"]
						for _, label := range labels {
							if label != "instance" && strutil.IsNotBlank(prometheusResult.Metric[label]) {
								key = key + " - " + prometheusResult.Metric[label]
							}
						}
						resultMap[key] = prometheusResult
					}
				}
			}
		}(statistics, pql, start, end, param.Statistics, labels, group, m)

		ret[statistics] = m
		if result == nil {
			result = m
		}
	}
	group.Wait()
	var list []*form.ReportForm
	for k, v := range result {
		for i := range v.Values {
			dataMap := make(map[string][]interface{})
			for calcStyle, d := range ret {
				dataMap[calcStyle] = d[k].Values[i]
			}
			if f := s.buildAggregationReportForm(v.Metric["instance"], k, item.Name, instanceMap, dataMap); f != nil {
				list = append(list, f)
			}
		}
	}
	return list, nil
}
