package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type KafkaInstanceService struct {
	service.InstanceServiceImpl
}

type KafkaQueryPageForm struct {
	TenantId   string `json:"tenantId"`
	SearchName string `json:"searchName"`
	PageNum    int    `json:"pageNum"`
	PageSize   int    `json:"pageSize"`
}

type KafkaQueryPageVO struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    KafkaPageVO `json:"data"`
}

type KafkaPageVO struct {
	Total    int       `json:"total"`
	PageNum  int       `json:"pageNum"`
	PageSize int       `json:"pageSize"`
	Records  []KafkaVO `json:"records"`
}

type KafkaVO struct {
	InstanceID    string `json:"instanceID"`
	ClusterName   string `json:"clusterName"`
	Storage       int    `json:"storage"`
	MqVersion     string `json:"mqVersion"`
	State         string `json:"state"`
	ResourceClass string `json:"resourceClass"`
}

type KafkaDTO struct {
	Param    string
	TenantId string
	UserInfo string
}

func (kafka *KafkaInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	param := "?pageNum=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&searchName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&searchID" + form.InstanceId
	}
	if strutil.IsNotBlank(form.StatusList) {
		param += "&searchState" + form.StatusList
	}
	return KafkaDTO{Param: param, TenantId: form.TenantId}
}

func (kafka *KafkaInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(KafkaDTO)
	respStr, err := httputil.HttpHeaderGet(url+param.Param, map[string]string{"CECLOUD-CSP-USER": "{\"tenantId\":\"" + param.TenantId + "\"}"})
	if err != nil {
		return nil, err
	}
	var resp KafkaQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (kafka *KafkaInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(KafkaQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceID,
				InstanceName: d.ClusterName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.State,
				}, {
					Name:  "mqVersion",
					Value: d.MqVersion,
				}, {
					Name:  "storage",
					Value: strconv.Itoa(d.Storage),
				}, {
					Name:  "resourceClass",
					Value: d.ResourceClass,
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func (kafka *KafkaInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
	param := "?pageNum=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&searchName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&searchID" + form.InstanceId
	}
	if strutil.IsNotBlank(form.StatusList) {
		param += "&searchState" + form.StatusList
	}
	return KafkaDTO{Param: param, TenantId: form.TenantId, UserInfo: form.IamInfo.UserInfo}
}

func (kafka *KafkaInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(KafkaDTO)
	respStr, err := httputil.HttpHeaderGet(url+param.Param, map[string]string{"user-info": param.UserInfo})
	if err != nil {
		return nil, err
	}
	var resp KafkaQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (kafka *KafkaInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(KafkaQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceID,
				InstanceName: d.ClusterName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.State,
				}, {
					Name:  "mqVersion",
					Value: d.MqVersion,
				}, {
					Name:  "storage",
					Value: strconv.Itoa(d.Storage),
				}, {
					Name:  "resourceClass",
					Value: d.ResourceClass,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
