package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"fmt"
	"strconv"
	"strings"
)

type MonitorReportFormService struct {
}

func NewMonitorReportFormService() *MonitorReportFormService {
	return &MonitorReportFormService{}
}

func (s *MonitorReportFormService) GetData(request forms.PrometheusRequest) (*forms.PrometheusValue, error) {
	if request.Instance == "" {
		return nil, errors.NewBusinessError("instance为空")
	}
	pql := getPql(request)
	prometheusResponse := Query(pql, request.Time)
	if len(prometheusResponse.Data.Result) == 0 {
		return nil, nil
	}
	result := prometheusResponse.Data.Result[0].Value
	prometheusValue := &forms.PrometheusValue{
		Time:  strconv.Itoa(int(result[0].(float64))),
		Value: changeDecimal(result[1].(string)),
	}
	return prometheusValue, nil
}

func (s *MonitorReportFormService) GetTop(request forms.PrometheusRequest) ([]forms.PrometheusInstance, error) {
	var pql string
	instances, err := getEcsInstances(request.TenantId)
	if err != nil {
		return nil, err
	}
	if request.Name == constant.EcsCpuUsage {
		pql = strings.ReplaceAll(constant.EcsCpuUsageTopExpr, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
	} else {
		monitorItem := getMonitorItemByName(request.Name)
		pql = strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
	}
	prometheusResponse := Query(pql, request.Time)
	result := prometheusResponse.Data.Result
	var instanceList []forms.PrometheusInstance
	for i := range result {
		instanceDTO := forms.PrometheusInstance{
			Instance: result[i].Metric[constant.INSTANCE],
			Value:    changeDecimal(result[i].Value[1].(string)),
		}
		instanceList = append(instanceList, instanceDTO)
	}
	return instanceList, nil
}

func (s *MonitorReportFormService) GetAxisData(request forms.PrometheusRequest) (forms.PrometheusAxis, error) {
	if request.Instance == "" {
		return forms.PrometheusAxis{}, errors.NewBusinessError("instance为空")
	}
	pql := getPql(request)
	prometheusResponse := QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	result := prometheusResponse.Data.Result

	labels := strings.Split(getMonitorItemByName(request.Name).Labels, ",")
	var label string
	for i := range labels {
		if labels[i] != "instance" {
			label = labels[i]
		}
	}

	start := request.Start
	end := request.End
	step := request.Step
	var timeList []string
	if len(result) == 0 {
		timeList = getTimeList(start, end, step, start)
	} else {
		timeList = getTimeList(start, end, step, int(result[0].Values[0][0].(float64)))
	}

	prometheusAxis := forms.PrometheusAxis{
		XAxis: timeList,
		YAxis: yAxisFillEmptyData(result, timeList, label, request.Name),
	}
	return prometheusAxis, nil
}

func yAxisFillEmptyData(result []forms.PrometheusResult, timeList []string, label string, metricName string) map[string][]string {
	resultMap := make(map[string][]string)
	for i := range result {
		timeMap := map[string]string{}
		for j := range result[i].Values {
			key := strconv.Itoa(int(result[i].Values[j][0].(float64)))
			timeMap[key] = result[i].Values[j][1].(string)
		}
		var key string
		var arr []string
		for k := range timeList {
			arr = append(arr, changeDecimal(timeMap[timeList[k]]))
		}
		if result[i].Metric[label] == "" {
			key = metricName
		} else {
			key = result[i].Metric[label]
		}
		resultMap[key] = arr
	}
	return resultMap
}

func getPql(request forms.PrometheusRequest) string {
	monitorItem := getMonitorItemByName(request.Name)
	metricLabels := constant.INSTANCE + "='" + request.Instance + "'," + constant.FILTER
	return strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, metricLabels)
}

func getTimeList(start int, end int, step int, firstTime int) []string {
	var timeList []string
	if start > end {
		return timeList
	}
	for firstTime-step >= start {
		firstTime -= step
	}
	for firstTime <= end {
		timeList = append(timeList, strconv.Itoa(firstTime))
		firstTime += step
	}
	return timeList
}

func changeDecimal(value string) string {
	if tools.IsBlank(value) {
		return ""
	}
	v, _ := strconv.ParseFloat(value, 64)
	return fmt.Sprintf("%.2f", v)
}

//根据监控名查询监控项
func getMonitorItemByName(name string) models.MonitorItem {
	monitorItemDao := dao.MonitorItem
	return monitorItemDao.GetMonitorItemByName(name)
}

//查询租户的ECS实例列表
func getEcsInstances(tenantId string) (string, error) {
	list, err := GetInstanceList(constant.Ecs, tenantId)
	if err != nil {
		return "", err
	}
	return strings.Join(list, "|"), nil
}
