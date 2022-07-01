package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
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
	monitorItem := dao.MonitorItem.GetMonitorItemCacheByName(request.Name)
	if strutil.IsBlank(monitorItem.MetricsLinux) {
		return nil, errors.NewBusinessError("指标不存在")
	}
	if !checkUserInstanceIdentity(request.TenantId, monitorItem.ProductBizId, request.Instance) {
		return nil, errors.NewBusinessError("该租户无此实例")
	}
	pql := strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+constant.FILTER)
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
	monitorItem := dao.MonitorItem.GetMonitorItemCacheByName(request.Name)
	var pql string
	list, err := getInstanceList("1", request.TenantId, "")
	if err != nil {
		return nil, err
	}
	if len(strings.Split(monitorItem.Labels, ",")) > 1 {
		if len(list) == 0 {
			return nil, nil
		}
		for i, v := range list {
			list[i] = fmt.Sprintf(constant.TopExpr, "1", strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+v+"'"))
		}

		pql = fmt.Sprintf(constant.TopExpr, strconv.Itoa(request.TopNum), strings.Join(list, " or "))
	} else {
		instances := strings.Join(list, "|")
		if strutil.IsBlank(instances) {
			return nil, nil
		}
		pql = fmt.Sprintf(constant.TopExpr, strconv.Itoa(request.TopNum), strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'"))
	}
	result := Query(pql, request.Time).Data.Result
	var instanceList []form.PrometheusInstance
	for _, v := range result {
		instanceDTO := form.PrometheusInstance{
			Instance: v.Metric[constant.INSTANCE],
			Value:    changeDecimal(v.Value[1].(string)),
		}
		instanceList = append(instanceList, instanceDTO)
	}
	return instanceList, nil
}

func (s *MonitorReportFormService) GetAxisData(request form.PrometheusRequest) (*form.PrometheusAxis, error) {
	if strutil.IsBlank(request.Instance) {
		return nil, errors.NewBusinessError("instance为空")
	}
	if request.Start == 0 || request.End == 0 || request.Start > request.End {
		return nil, errors.NewBusinessError("时间参数错误")
	}
	monitorItem := dao.MonitorItem.GetMonitorItemCacheByName(request.Name)
	if strutil.IsBlank(monitorItem.MetricsLinux) {
		return nil, errors.NewBusinessError("指标不存在")
	}
	if !checkUserInstanceIdentity(request.TenantId, monitorItem.ProductBizId, request.Instance) {
		return nil, errors.NewBusinessError("该租户无此实例")
	}
	pql := strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+constant.FILTER)
	prometheusResponse := QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	result := prometheusResponse.Data.Result

	labels := strings.Split(monitorItem.Labels, ",")
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
		YAxis: yAxisFillEmptyData(result, timeList, labels, request.Instance),
	}
	return prometheusAxis, nil
}

func (s *MonitorReportFormService) GetNetworkData(request form.PrometheusRequest) (*form.NetworkData, error) {
	if request.Start == 0 || request.End == 0 || request.Start > request.End {
		return nil, errors.NewBusinessError("时间参数错误")
	}
	monitorItem := dao.MonitorItem.GetMonitorItemCacheByName(request.Name)
	if strutil.IsBlank(monitorItem.MetricsLinux) {
		return nil, errors.NewBusinessError("指标不存在")
	}
	var instance string
	if strutil.IsNotBlank(request.Instance) {
		instance = constant.INSTANCE + "='" + request.Instance + "'"
	} else if strutil.IsNotBlank(request.TenantId) {
		list, err := getInstanceList("2", request.TenantId, "")
		if err != nil {
			return nil, err
		}
		instance = constant.INSTANCE + "=~'" + strings.Join(list, "|") + "'"
	}
	pql := strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, instance)
	prometheusResponse := QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	result := prometheusResponse.Data.Result
	start := request.Start
	end := request.End
	step := request.Step
	var timeList []string
	if len(result) == 0 {
		timeList = getTimeList(start, end, step, start)
	} else {
		timeList = getTimeList(start, end, step, int(result[0].Values[0][0].(float64)))
	}
	networkData := &form.NetworkData{
		TimeAxis:  timeList,
		ValueAxis: getValueAxis(result, timeList),
	}
	return networkData, nil
}

func yAxisFillEmptyData(result []form.PrometheusResult, timeList []string, labels []string, instanceId string) map[string][]string {
	resultMap := make(map[string][]string)
	for _, v1 := range result {
		timeMap := map[string]string{}
		for _, v2 := range v1.Values {
			key := strconv.Itoa(int(v2[0].(float64)))
			timeMap[key] = v2[1].(string)
		}
		var labelList []string
		var key string
		var arr []string
		for _, v3 := range timeList {
			arr = append(arr, changeDecimal(timeMap[v3]))
		}
		for _, v4 := range labels {
			if v4 != "instance" && strutil.IsNotBlank(v1.Metric[v4]) {
				labelList = append(labelList, v1.Metric[v4])
			}
		}
		if len(labelList) == 0 {
			key = instanceId
		} else {
			key = strings.Join(labelList, "-")
		}
		resultMap[key] = arr
	}
	return resultMap
}

func getValueAxis(result []form.PrometheusResult, timeList []string) []string {
	var valueAxis []string
	for _, v1 := range result {
		timeMap := map[string]string{}
		for _, v2 := range v1.Values {
			key := strconv.Itoa(int(v2[0].(float64)))
			timeMap[key] = v2[1].(string)
		}
		for _, v3 := range timeList {
			valueAxis = append(valueAxis, changeDecimal(timeMap[v3]))
		}
	}
	return valueAxis
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

//校验该租户下是否拥有该实例
func checkUserInstanceIdentity(tenantId, productBizId, instanceId string) bool {
	if strutil.IsBlank(tenantId) {
		return true
	}
	list, err := getInstanceList(productBizId, tenantId, instanceId)
	if err != nil {
		logger.Logger().Error("获取实例列表失败")
		return false
	}
	if len(list) == 0 {
		return false
	}
	return true
}

// GetInstanceList 获取实例ID列表
func getInstanceList(productBizId, tenantId, instanceId string) ([]string, error) {
	f := commonService.InstancePageForm{
		TenantId: tenantId,
		Product:  dao.MonitorProduct.GetMonitorProductByBizId(productBizId).Abbreviation,
		Current:  1,
		PageSize: 10000,
	}
	if strutil.IsNotBlank(instanceId) {
		f.InstanceId = instanceId
		f.PageSize = 1
	}
	instanceService := external.ProductInstanceServiceMap[f.Product]
	stage, ok := instanceService.(commonService.InstanceStage)
	if !ok {
		return nil, errors.NewBusinessError("获取监控产品服务失败")
	}
	page, err := instanceService.GetPage(f, stage)
	if err != nil {
		return nil, errors.NewBusinessError(err.Error())
	}
	var instanceList []string
	for _, v := range page.Records.([]commonService.InstanceCommonVO) {
		instanceList = append(instanceList, v.InstanceId)
	}
	return instanceList, nil
}
