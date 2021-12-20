package external

import (
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"strconv"
	"strings"
)

type EcsInstanceService struct {
	commonService.InstanceServiceImpl
}

func (ecs *EcsInstanceService) ConvertRealForm(form commonService.InstancePageForm) interface{} {
	param := forms.EcsQueryPageForm{
		TenantId:     form.TenantId,
		Current:      form.Current,
		PageSize:     form.PageSize,
		InstanceName: form.InstanceName,
		InstanceId:   form.InstanceId,
	}
	if tools.IsNotBlank(form.StatusList) {
		param.StatusList = toIntList(form.StatusList)
	}
	return param
}

func (ecs *EcsInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := tools.HttpPostJson(url, form, nil)
	if err != nil {
		return nil, err
	}
	var resp forms.EcsQueryPageVO
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (ecs *EcsInstanceService) ConvertResp(realResp interface{}) (int, []commonService.InstanceCommonVO) {
	vo := realResp.(forms.EcsQueryPageVO)
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
	var request = forms.PrometheusRequest{
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
