package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"fmt"
	"strconv"
	"strings"
)

type MonitorChartService struct {
	prometheus *PrometheusService
}

func NewMonitorChartService() *MonitorChartService {
	return &MonitorChartService{
		prometheus: NewPrometheusService(),
	}
}

func (s *MonitorChartService) GetData(request form.PrometheusRequest) (*form.PrometheusValue, error) {
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
	prometheusResponse := s.prometheus.Query(pql, request.Time)
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

func (s *MonitorChartService) GetTop(request form.PrometheusRequest) ([]form.PrometheusInstance, error) {
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
	result := s.prometheus.Query(pql, request.Time).Data.Result
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

func (s *MonitorChartService) GetAxisData(request form.PrometheusRequest) (*form.PrometheusAxis, error) {
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
	if strutil.IsNotBlank(request.Pid) {
		pql = strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+fmt.Sprintf(constant.PId, request.Pid))
	}
	if strutil.IsNotBlank(request.Statistics) {
		pql = fmt.Sprintf("%s_over_time((%s)[%s:1m])", request.Statistics, pql, request.Scope)
	}
	prometheusResponse := s.prometheus.QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
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

func (s *MonitorChartService) GetNetworkData(request form.PrometheusRequest) (*form.NetworkData, error) {
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
	prometheusResponse := s.prometheus.QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
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

func (s *MonitorChartService) GetProcessData(request form.PrometheusRequest) ([]form.ProcessData, error) {
	if strutil.IsBlank(request.Instance) {
		return nil, errors.NewBusinessError("instance为空")
	}
	if request.Start == 0 || request.End == 0 || request.Start > request.End {
		return nil, errors.NewBusinessError("时间参数错误")
	}
	cpu := dao.MonitorItem.GetMonitorItemCacheByName("ecs_processes_top5Cpus")
	mem := dao.MonitorItem.GetMonitorItemCacheByName("ecs_processes_top5Mems")
	fd := dao.MonitorItem.GetMonitorItemCacheByName("ecs_processes_top5Fds")
	if !checkUserInstanceIdentity(request.TenantId, "1", request.Instance) {
		return nil, errors.NewBusinessError("该租户无此实例")
	}
	cpuPql := strings.ReplaceAll(cpu.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+constant.FILTER)
	memPql := strings.ReplaceAll(mem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+constant.FILTER)
	fdPql := strings.ReplaceAll(fd.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+request.Instance+"',"+constant.FILTER)
	cpuResponse := s.prometheus.QueryRange(cpuPql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	memResponse := s.prometheus.QueryRange(memPql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	fdResponse := s.prometheus.QueryRange(fdPql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step))
	memMap := make(map[string]*form.PrometheusResult)
	fdMap := make(map[string]*form.PrometheusResult)
	for _, v := range memResponse.Data.Result {
		memMap[v.Metric["pid"]] = v
	}
	for _, v := range fdResponse.Data.Result {
		fdMap[v.Metric["pid"]] = v
	}
	var processList []form.ProcessData
	for _, v := range cpuResponse.Data.Result {
		process := form.ProcessData{
			Pid:     v.Metric["pid"],
			CmdLine: v.Metric["cmd_line"],
			Name:    getProcessName(v.Metric["cmd_line"]),
		}
		if len(v.Values) != 0 {
			process.Time = util.TimestampToFullTimeFmtStr(int64(v.Values[len(v.Values)-1][0].(float64)))
			process.Cpu = changeDecimal(v.Values[len(v.Values)-1][1].(string))
			process.Memory = changeDecimal(memMap[process.Pid].Values[len(memMap[process.Pid].Values)-1][1].(string))
			process.Openfiles = fdMap[process.Pid].Values[len(fdMap[process.Pid].Values)-1][1].(string)
		}
		processList = append(processList, process)
	}
	return processList, nil
}

func (s *MonitorChartService) GetTopDataByIam(request form.PrometheusRequest) ([]form.PrometheusInstance, error) {
	resourcesIdList, err := GetIamResourcesIdList(request.IamUserId)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
	}
	monitorItem := dao.MonitorItem.GetMonitorItemCacheByName(request.Name)
	var pql string
	if len(strings.Split(monitorItem.Labels, ",")) > 1 {
		for i, v := range resourcesIdList {
			resourcesIdList[i] = fmt.Sprintf(constant.TopExpr, "1", strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"='"+v+"'"))
		}
		pql = fmt.Sprintf(constant.TopExpr, strconv.Itoa(request.TopNum), strings.Join(resourcesIdList, " or "))
	} else {
		instances := strings.Join(resourcesIdList, "|")
		if strutil.IsBlank(instances) {
			return nil, nil
		}
		pql = fmt.Sprintf(constant.TopExpr, strconv.Itoa(request.TopNum), strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'"))
	}
	result := s.prometheus.Query(pql, request.Time).Data.Result
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

func yAxisFillEmptyData(result []*form.PrometheusResult, timeList []string, labels []string, instanceId string) map[string][]string {
	resultMap := make(map[string][]string)
	for _, v := range result {
		timeMap := map[string]string{}
		for _, value := range v.Values {
			key := strconv.Itoa(int(value[0].(float64)))
			timeMap[key] = value[1].(string)
		}
		var key string
		var arr []string
		for _, time := range timeList {
			arr = append(arr, changeDecimal(timeMap[time]))
		}
		key = instanceId
		for _, label := range labels {
			if label != "instance" && strutil.IsNotBlank(v.Metric[label]) {
				key = key + " - " + v.Metric[label]
			}
		}
		resultMap[key] = arr
	}
	return resultMap
}

func getValueAxis(result []*form.PrometheusResult, timeList []string) []string {
	var valueAxis []string
	for _, v := range result {
		timeMap := map[string]string{}
		for _, value := range v.Values {
			key := strconv.Itoa(int(value[0].(float64)))
			timeMap[key] = value[1].(string)
		}
		for _, time := range timeList {
			valueAxis = append(valueAxis, changeDecimal(timeMap[time]))
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

func getProcessName(cmdLine string) string {
	if strutil.IsBlank(cmdLine) {
		return "unknown"
	}
	list := strings.Split(strings.Split(cmdLine, " ")[0], "/")
	return list[len(list)-1]
}
