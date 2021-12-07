package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"strconv"
)

type EcsInstanceService struct {
	service.InstanceServiceImpl
}

type EcsQueryPageForm struct {
	TenantId     string
	Current      int
	PageSize     int
	InstanceName string
	InstanceId   string
	Status       int
	StatusList   []int
}

type EcsQueryPageVO struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    EcsPageVO `json:"data"`
}

type EcsPageVO struct {
	Total int      `json:"total"`
	Rows  []*ECSVO `json:"rows"`
}

type ECSVO struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Region       string `json:"region"`
	Ip           string `yaml:"ip"`
	Status       int    `yaml:"status"`
}

func (ecs *EcsInstanceService) convertRealForm(form service.InstancePageForm) interface{} {
	param := EcsQueryPageForm{
		TenantId:     form.TenantId,
		Current:      form.Current,
		PageSize:     form.PageSize,
		InstanceName: form.InstanceName,
		InstanceId:   form.InstanceId,
	}
	if form.StatusList != nil && len(form.StatusList) > 0 {
		var list []int
		for _, s := range form.StatusList {
			i, err := strconv.Atoi(s)
			if err != nil {
				list = append(list, i)
			}
		}
		param.StatusList = list
	}
	return param
}

func (ecs *EcsInstanceService) doRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := tools.HttpPostJson(url, form, nil)
	if err != nil {
		return nil, err
	}
	var resp EcsQueryPageVO
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (ecs *EcsInstanceService) convertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EcsQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				Id:   d.InstanceId,
				Name: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: string(rune(d.Status)),
				}},
			})
		}
	}
	return vo.Data.Total, list
}
