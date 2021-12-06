package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"net/http"
	"net/url"
)

func Query(pql string, time string) forms.PrometheusResponse {
	var cfg = config.GetPrometheusConfig()
	requestUrl := cfg.Url + cfg.Query
	pql = url.QueryEscape(pql)
	if time != "" {
		pql += "&time=" + time
	}
	return sendRequest(requestUrl, pql)
}

func QueryRange(pql string, start string, end string, step string) forms.PrometheusResponse {
	var cfg = config.GetPrometheusConfig()
	requestUrl := cfg.Url + cfg.QueryRange
	pql = url.QueryEscape(pql) + "&start=" + start + "&end=" + end + "&step=" + step
	return sendRequest(requestUrl, pql)
}

func sendRequest(requestUrl string, pql string) forms.PrometheusResponse {
	logger.Logger().Infof("url:%v\n", requestUrl+pql)
	response, err := http.Get(requestUrl + pql)
	if err != nil {
		logger.Logger().Errorf("error:%v\n", err)
		return forms.PrometheusResponse{}
	}
	var prometheusResponse forms.PrometheusResponse
	err = json.NewDecoder(response.Body).Decode(&prometheusResponse)
	if err != nil {
		logger.Logger().Errorf("error:%v\n", err)
	}
	return prometheusResponse
}
