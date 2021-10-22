package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/prometheus"
	"container/list"
	"fmt"
	"strconv"
)

type MonitorReportFormService struct {
}

func NewMonitorReportFormService() *MonitorReportFormService {
	return &MonitorReportFormService{}
}

func (mpd *MonitorReportFormService) GetData(request forms.PrometheusRequest) []forms.PrometheusValue {
	pql := getPql(request)
	prometheusResponse := prometheus.Query(pql, request.Time, request.TenantId)
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

//func (mpd *MonitorReportFormService) getAxisData(request forms.PrometheusRequest) forms.PrometheusAxis {
//	pql := getPql(request)
//	prometheusResponse := prometheus.Query(pql, request.Time, request.TenantId)
//	result := prometheusResponse.Data.Result
//
//	labels := strings.Split(request.Labels, ",")
//	var lable string
//	for i := range labels{
//		if labels[i] != "instance" {
//			lable = labels[i]
//		}
//	}
//
//	var timeList list.List
//	start := request.Start
//	end := request.End
//	step := request.Step
//	if len(result) == 0 {
//		timeList = getTimeList(start, end, step, start)
//	}else {
//		timeList = getTimeList(start, end, step, result[0].Values[0][0].(int64))
//	}
//
//	prometheusAxis := forms.PrometheusAxis{
//		XAxis: timeList,
//		YAxis: lable,
//	}
//	return prometheusAxis
//}

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

func getTimeList(start int64, end int64, step int64, firstTime int64) list.List {
	var timeList list.List
	if start > end {
		return timeList
	}
	for firstTime-step >= start {
		firstTime -= step
	}
	for firstTime <= end {
		timeList.PushBack(firstTime)
		firstTime += step
	}
	return timeList
}
