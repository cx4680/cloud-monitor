package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"fmt"
	"strconv"
	"strings"
)

type MonitorReportFormService struct {
}

func NewMonitorReportFormService() *MonitorReportFormService {
	return &MonitorReportFormService{}
}

func (s *MonitorReportFormService) GetData(request form.PrometheusRequest) (*form.PrometheusValue, error) {
	if strutil.IsBlank(request.Instance) {
		return nil, errors.NewBusinessError("instance为空")
	}
	if !checkUserInstanceIdentity(request.TenantId, request.ProductId, request.Instance) {
		return nil, errors.NewBusinessError("该租户无此实例")
	}
	pql := getPql(request)
	prometheusResponse := Query(pql, request.Time)
	if len(prometheusResponse.Data.Result) == 0 {
		return nil, nil
	}
	result := prometheusResponse.Data.Result[0].Value
	prometheusValue := &form.PrometheusValue{
		Time:  strconv.Itoa(int(result[0].(float64))),
		Value: changeDecimal(result[1].(string)),
	}
	return prometheusValue, nil
}

func (s *MonitorReportFormService) GetTop(request form.PrometheusRequest) ([]form.PrometheusInstance, error) {
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
	var instanceList []form.PrometheusInstance
	for i := range result {
		instanceDTO := form.PrometheusInstance{
			Instance: result[i].Metric[constant.INSTANCE],
			Value:    changeDecimal(result[i].Value[1].(string)),
		}
		instanceList = append(instanceList, instanceDTO)
	}
	return instanceList, nil
}

func (s *MonitorReportFormService) GetAxisData(request form.PrometheusRequest) (*form.PrometheusAxis, error) {
	if strutil.IsBlank(request.Instance) {
		return nil, errors.NewBusinessError("instance为空")
	}
	if !checkUserInstanceIdentity(request.TenantId, request.ProductId, request.Instance) {
		return nil, errors.NewBusinessError("该租户无此实例")
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

	prometheusAxis := &form.PrometheusAxis{
		XAxis: timeList,
		YAxis: yAxisFillEmptyData(result, timeList, label, request.Name),
	}
	return prometheusAxis, nil
}

func yAxisFillEmptyData(result []form.PrometheusResult, timeList []string, label string, metricName string) map[string][]string {
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
		if strutil.IsBlank(result[i].Metric[label]) {
			key = metricName
		} else {
			key = result[i].Metric[label]
		}
		resultMap[key] = arr
	}
	return resultMap
}

func getPql(request form.PrometheusRequest) string {
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
	if strutil.IsBlank(value) {
		return ""
	}
	v, _ := strconv.ParseFloat(value, 64)
	return fmt.Sprintf("%.2f", v)
}

//根据监控名查询监控项
func getMonitorItemByName(name string) model.MonitorItem {
	monitorItemDao := dao.MonitorItem
	return monitorItemDao.GetMonitorItemByName(name)
}

//查询租户的ECS实例列表
func getEcsInstances(tenantId string) (string, error) {
	list, err := GetInstanceList("1", tenantId)
	if err != nil {
		return "", err
	}
	return strings.Join(list, "|"), nil
}

//校验该租户下是否拥有该实例
func checkUserInstanceIdentity(tenantId, productId, instanceId string) bool {
	list, err := GetInstanceList(productId, tenantId)
	if err != nil {
		logger.Logger().Error("获取实例列表失败")
		return false
	}
	for _, v := range list {
		if instanceId == v {
			return true
		}
	}
	return false
}

//获取实例ID列表
func GetInstanceList(productId string, tenantId string) ([]string, error) {
	f := commonService.InstancePageForm{
		TenantId: tenantId,
		Product:  constant.ProductMap[productId],
		Current:  1,
		PageSize: 10000,
	}
	instanceService := external.ProductInstanceServiceMap[f.Product]
	page, err := instanceService.GetPage(f, instanceService.(commonService.InstanceStage))
	if err != nil {
		return nil, errors.NewBusinessError(err.Error())
	}
	var instanceList []string
	for _, v := range page.Records.([]commonService.InstanceCommonVO) {
		instanceList = append(instanceList, v.InstanceId)
	}
	return instanceList, nil
}
