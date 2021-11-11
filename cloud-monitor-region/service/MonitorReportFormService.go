package service

import (
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

func (mpd *MonitorReportFormService) GetData(request forms.PrometheusRequest) []forms.PrometheusValue {
	pql := getPql(request)
	prometheusResponse := Query(pql, request.Time, request.TenantId)
	var Values []forms.PrometheusValue
	if len(prometheusResponse.Data.Result) == 0 {
		return Values
	}
	result := prometheusResponse.Data.Result[0].Value
	prometheusValue := forms.PrometheusValue{
		Time:  strconv.Itoa(int(result[0].(float64))),
		Value: result[1].(string),
	}
	return []forms.PrometheusValue{prometheusValue}
}

func (mpd *MonitorReportFormService) GetTop(request forms.PrometheusRequest) []forms.PrometheusInstance {
	pql := fmt.Sprintf("topk(%s,%s{%s})", constant.TOP_NUM, request.Name, constant.INSTANCE+"=~'"+request.Instance+"'")
	prometheusResponse := Query(pql, request.Time, request.TenantId)
	result := prometheusResponse.Data.Result
	var instanceList []forms.PrometheusInstance
	for i := range result {
		instanceDTO := forms.PrometheusInstance{
			Instance: result[i].Metric[constant.INSTANCE],
			Value:    result[i].Value[1].(string),
		}
		instanceList = append(instanceList, instanceDTO)
	}
	return instanceList
}

func (mpd *MonitorReportFormService) GetAxisData(request forms.PrometheusRequest) forms.PrometheusAxis {
	pql := getPql(request)
	prometheusResponse := QueryRange(pql, strconv.Itoa(request.Start), strconv.Itoa(request.End), strconv.Itoa(request.Step), request.TenantId)
	result := prometheusResponse.Data.Result

	labels := strings.Split(request.Labels, ",")
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
		YAxis: yAxisFillEmptyData(result, timeList, label),
	}
	return prometheusAxis
}

func yAxisFillEmptyData(Result []forms.PrometheusResult, timeList []string, label string) map[string][]string {
	resultMap := make(map[string][]string)
	for i := range Result {
		timeMap := map[string]string{}
		for j := range Result[i].Values {
			key := strconv.Itoa(int(Result[i].Values[j][0].(float64)))
			timeMap[key] = Result[i].Values[j][1].(string)
		}
		var key string
		var arr []string
		for k := range timeList {
			arr = append(arr, changeDecimal(timeMap[timeList[k]]))
		}
		if Result[i].Metric[label] == "" {
			key = Result[i].Metric["__name__"]

		} else {
			key = Result[i].Metric[label]
		}
		resultMap[key] = arr
	}
	return resultMap
}

func getPql(request forms.PrometheusRequest) string {
	metricLabels := constant.INSTANCE + "='" + request.Instance + "'," + constant.FILTER
	statistics := ""
	scope := ""
	if request.Statistics != "" {
		statistics = request.Statistics + constant.STATISTICS
		scope = "[" + request.Scope + "]"
	}
	return fmt.Sprintf("%s(%s{%s}%s)", statistics, request.Name, metricLabels, scope)
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
	v, _ := strconv.ParseFloat(value, 64)
	return fmt.Sprintf("%.2f", v)
}
