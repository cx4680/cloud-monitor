package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"encoding/json"
	"net/http"
	"net/url"
)

func Query(pql string, time string) form.PrometheusResponse {
	var cfg = config.Cfg.Prometheus
	requestUrl := cfg.Url + cfg.Query
	logger.Logger().Info(requestUrl + pql)
	pql = url.QueryEscape(pql)
	if strutil.IsNotBlank(time) {
		pql += "&time=" + time
	}
	return sendRequest(requestUrl, pql)
}

func QueryRange(pql string, start string, end string, step string) form.PrometheusResponse {
	var cfg = config.Cfg.Prometheus
	requestUrl := cfg.Url + cfg.QueryRange
	logger.Logger().Info(requestUrl + pql)
	pql = url.QueryEscape(pql) + "&start=" + start + "&end=" + end + "&step=" + step
	return sendRequest(requestUrl, pql)
}

func sendRequest(requestUrl string, pql string) form.PrometheusResponse {
	response, err := http.Get(requestUrl + pql)
	if err != nil {
		logger.Logger().Errorf("error:%v\n", err)
		return form.PrometheusResponse{}
	}
	var prometheusResponse form.PrometheusResponse
	err = json.NewDecoder(response.Body).Decode(&prometheusResponse)
	if err != nil {
		logger.Logger().Errorf("error:%v\n", err)
	}
	return prometheusResponse
}
