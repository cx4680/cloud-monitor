package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"net/http"
	"strings"
)

func Query(pql string, time string) forms.PrometheusResponse {
	var cfg = config.GetPrometheusConfig()
	url := cfg.Url + cfg.Query
	if time != "" {
		pql += "&time=" + time
	}
	return sendRequest(url, pql)
}

func QueryRange(pql string, start string, end string, step string) forms.PrometheusResponse {
	var cfg = config.GetPrometheusConfig()
	url := cfg.Url + cfg.QueryRange
	pql += "&start=" + start + "&end=" + end + "&step=" + step
	return sendRequest(url, pql)
}

func sendRequest(url string, pql string) forms.PrometheusResponse {
	logger.Logger().Infof("url:%v\n", url+pql)
	response, err := http.Get(url + strings.ReplaceAll(pql, " ", ""))
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
