package external

import (
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
	"strings"
)

type EcsInstanceService struct {
	commonService.InstanceServiceImpl
}

func (ecs *EcsInstanceService) ConvertRealForm(f commonService.InstancePageForm) interface{} {
	param := form.EcsQueryPageForm{
		TenantId:     f.TenantId,
		Current:      f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
	}
	if strutil.IsNotBlank(f.StatusList) {
		param.StatusList = toIntList(f.StatusList)
	}
	return param
}

func (ecs *EcsInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	respStr, err := httputil.HttpPostJson(url, f, nil)
	if err != nil {
		return nil, err
	}
	var resp form.EcsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (ecs *EcsInstanceService) ConvertResp(realResp interface{}) (int, []commonService.InstanceCommonVO) {
	vo := realResp.(form.EcsQueryPageVO)
	var list []commonService.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, commonService.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.InstanceName,
				Labels: []commonService.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "ecsCpuUsage",
					Value: getMonitorData("ecs_cpu_usage", d.InstanceId),
				}, {
					Name:  "ecsMemoryUsage",
					Value: getMonitorData("ecs_memory_usage", d.InstanceId),
				}, {
					Name:  "osType",
					Value: d.OsType,
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func getMonitorData(metricName string, instanceId string) string {
	var request = form.PrometheusRequest{
		Name:     metricName,
		Instance: instanceId,
	}
	data, err := service.NewMonitorReportFormService().GetData(request)
	if err != nil || data == nil {
		return ""
	}
	return data.Value
}

func toIntList(s string) []int {
	statusList := strings.Split(s, ",")
	var list []int
	for _, v := range statusList {
		i, err := strconv.Atoi(v)
		if err == nil {
			list = append(list, i)
		}
	}
	return list
}
