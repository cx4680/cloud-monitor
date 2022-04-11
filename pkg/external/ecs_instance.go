package external

import (
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
	"strings"
)

type EcsInstanceService struct {
	commonService.InstanceServiceImpl
}

type EcsQueryPageForm struct {
	TenantId     string `json:"tenantId"`
	Current      int    `json:"current"`
	PageSize     int    `json:"pageSize"`
	InstanceName string `json:"instanceName"`
	InstanceId   string `json:"instanceId"`
	Status       int    `json:"status"`
	StatusList   []int  `json:"statusList"`
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
	Ip           string `json:"ip"`
	Status       int    `json:"status"`
	OsType       string `json:"osType"`
}

func (ecs *EcsInstanceService) ConvertRealForm(f commonService.InstancePageForm) interface{} {
	param := EcsQueryPageForm{
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
	var resp EcsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (ecs *EcsInstanceService) ConvertResp(realResp interface{}) (int, []commonService.InstanceCommonVO) {
	vo := realResp.(EcsQueryPageVO)
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
					Name:  "osType",
					Value: d.OsType,
				}},
			})
		}
	}
	return vo.Data.Total, list
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
