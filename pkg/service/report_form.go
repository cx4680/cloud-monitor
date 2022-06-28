package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ReportFormService struct {
}

func NewReportFormService() *ReportFormService {
	return &ReportFormService{}
}

func (s *ReportFormService) GetMonitorData(param form.ReportFormParam) ([]form.ReportForm, error) {
	var reportForm []form.ReportForm
	var instanceList []string
	instanceMap := make(map[string]form.InstanceForm)
	for _, v := range param.InstanceList {
		instanceList = append(instanceList, v.InstanceId)
		instanceMap[v.InstanceId] = v
	}
	instances := strings.Join(instanceList, "|")
	for _, v := range param.ItemList {
		item := dao.MonitorItem.GetMonitorItemCacheByName(v)
		pql := strings.ReplaceAll(item.MetricsLinux, constant.MetricLabel, constant.INSTANCE+"=~'"+instances+"'")
		result := QueryRange(pql, strconv.Itoa(param.Start), strconv.Itoa(param.End), strconv.Itoa(param.Step)).Data.Result
		if len(result) == 0 {
			continue
		}
		labels := strings.Split(item.Labels, ",")
		for _, v1 := range labels {
			if v1 != "instance" {

			}
		}
		if len(param.Statistics) == 0 {
			for _, v1 := range result {
				for _, v2 := range v1.Values {
					f := form.ReportForm{
						Region:       param.RegionCode,
						InstanceName: instanceMap[v1.Metric["instance"]].InstanceName,
						InstanceId:   v1.Metric["instance"],
						Status:       instanceMap[v1.Metric["instance"]].Status,
						ItemName:     item.Name,
						Time:         util.TimestampToFullTimeFmtStr(int64(v2[0].(float64))),
						Timestamp:    int64(v2[0].(float64)),
						Value:        changeDecimal(v2[1].(string)),
					}
					for _, v3 := range labels {
						if v3 != "instance" && strutil.IsNotBlank(v1.Metric[v3]) {
							f.InstanceId = f.InstanceId + " - " + v1.Metric[v3]
						}
					}
					reportForm = append(reportForm, f)
				}
			}
		} else {
			maxResult := make(map[string]form.PrometheusResult)
			minResult := make(map[string]form.PrometheusResult)
			avgResult := make(map[string]form.PrometheusResult)
			start := param.Start + (86400 - (param.Start-57600)%86400)
			end := param.End + (86400 - (param.End-57600)%86400)
			for _, v1 := range param.Statistics {
				expr := fmt.Sprintf("%s_over_time((%s)[1d:1m])", v1, pql)
				result = QueryRange(expr, strconv.Itoa(start), strconv.Itoa(end), "86400").Data.Result
				if v1 == "max" {
					for _, v2 := range result {
						maxResult[v2.Metric["instance"]] = v2
					}
				}
				if v1 == "min" {
					for _, v2 := range result {
						minResult[v2.Metric["instance"]] = v2
					}
				}
				if v1 == "avg" {
					for _, v2 := range result {
						avgResult[v2.Metric["instance"]] = v2
					}
				}
			}
			for _, v1 := range result {
				for i, v2 := range v1.Values {
					f := form.ReportForm{
						Region:       param.RegionCode,
						InstanceName: instanceMap[v1.Metric["instance"]].InstanceName,
						InstanceId:   v1.Metric["instance"],
						Status:       instanceMap[v1.Metric["instance"]].Status,
						ItemName:     item.Name,
						Time:         util.TimestampToDayTimeFmtStr(int64(v2[0].(float64)) - 1),
						Timestamp:    int64(v2[0].(float64) - 1),
					}
					if len(maxResult[v1.Metric["instance"]].Values) != 0 {
						f.MaxValue = changeDecimal(maxResult[v1.Metric["instance"]].Values[i][1].(string))
					}
					if len(minResult[v1.Metric["instance"]].Values) != 0 {
						f.MinValue = changeDecimal(minResult[v1.Metric["instance"]].Values[i][1].(string))
					}
					if len(avgResult[v1.Metric["instance"]].Values) != 0 {
						f.AvgValue = changeDecimal(avgResult[v1.Metric["instance"]].Values[i][1].(string))
					}
					for _, v3 := range labels {
						if v3 != "instance" && strutil.IsNotBlank(v1.Metric[v3]) {
							f.InstanceId = f.InstanceId + " - " + v1.Metric[v3]
						}
					}
					reportForm = append(reportForm, f)
				}
			}
		}
	}
	return reportForm, nil
}

func (s *ReportFormService) Export(param form.ReportFormParam, userInfo string) error {
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.Export
	asyncParams := []form.AsyncExportParam{
		{
			SheetSeq:   0,
			SheetName:  "云监控",
			SheetParam: jsonutil.ToString(param),
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

func (s *ReportFormService) QueryExportRecords(param form.ReportFormParam, userInfo string) (form.ExportRecords, error) {
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.ExportRecords + "?current=" + strconv.Itoa(param.Current) + "&pageSize=" + strconv.Itoa(param.PageSize)
	header := map[string]string{"user-info": userInfo}
	result, err := httputil.HttpHeaderGet(url, header)
	if err != nil {
		logger.Logger().Infof("AsyncExportError：%v", err)
		return form.ExportRecords{}, errors.NewBusinessError("获取下载列表失败")
	}
	var exportRecords form.ExportRecords
	jsonutil.ToObject(result, &exportRecords)
	return exportRecords, nil
}

func (s *ReportFormService) DownloadFile(param form.ReportFormParam, userInfo string) (io.ReadCloser, error) {
	url := config.Cfg.AsyncExport.Url + config.Cfg.AsyncExport.DownloadFile + "?fileId=" + param.FileId
	header := map[string]string{"user-info": userInfo}

	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger().Errorf("download error:%v", err)
	}
	defer resp.Body.Close()
	return resp.Body, nil
}
