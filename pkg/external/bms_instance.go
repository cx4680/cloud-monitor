package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
	"strings"
)

type BmsInstanceService struct {
	service.InstanceServiceImpl
}

type BmsRequest struct {
	TenantId string `json:"tenantId"`
	Params   string `json:"params"`
}

type BmsResponse struct {
	RequestId string             `json:"RequestId"`
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	Data      BmsQueryPageResult `json:"data"`
}

type BmsQueryPageResult struct {
	Servers    []BmsServers `json:"servers"`
	TotalCount int          `json:"total_count"`
}

type BmsServers struct {
	InstanceId string `json:"InstanceId"`
	Name       string `json:"name"`
	State      string `json:"State"`
	BmType     int    `json:"BmType"`
}

func (bms *BmsInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	var params = "?pageNumber=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	var filterList []string
	if strutil.IsNotBlank(form.InstanceName) {
		filterList = append(filterList, "name:lk:"+form.InstanceName)
	}
	if strutil.IsNotBlank(form.InstanceId) {
		filterList = append(filterList, "id:lk:"+form.InstanceId)
	}
	if strutil.IsNotBlank(form.StatusList) {
		filterList = append(filterList, "state:in:"+form.StatusList)
	}
	if strutil.IsNotBlank(form.ExtraAttr["bmType"]) {
		filterList = append(filterList, "BmType:in:"+form.ExtraAttr["bmType"])
	}
	if len(filterList) > 0 {
		filter := strings.Join(filterList, "|")
		params += "&filter=" + filter
	}
	return BmsRequest{TenantId: form.TenantId, Params: params}
}

func (bms *BmsInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(BmsRequest)
	url = strings.ReplaceAll(url, "{tenantId}", param.TenantId)
	respStr, err := httputil.HttpGet(url + param.Params)
	if err != nil {
		return nil, err
	}
	var resp BmsResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (bms *BmsInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(BmsResponse)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.Servers {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.State,
				}, {
					Name:  "bmType",
					Value: strconv.Itoa(d.BmType),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}

func (bms *BmsInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
	var params = "?pageNumber=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	var filterList []string
	if strutil.IsNotBlank(form.InstanceName) {
		filterList = append(filterList, "name:lk:"+form.InstanceName)
	}
	if strutil.IsNotBlank(form.InstanceId) {
		filterList = append(filterList, "id:lk:"+form.InstanceId)
	}
	if strutil.IsNotBlank(form.StatusList) {
		filterList = append(filterList, "state:in:"+form.StatusList)
	}
	if strutil.IsNotBlank(form.ExtraAttr["bmType"]) {
		filterList = append(filterList, "BmType:in:"+form.ExtraAttr["bmType"])
	}
	if len(filterList) > 0 {
		filter := strings.Join(filterList, "|")
		params += "&filter=" + filter
	}
	return BmsRequest{TenantId: form.TenantId, Params: params}
}

func (bms *BmsInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(BmsRequest)
	url = strings.ReplaceAll(url, "{tenantId}", param.TenantId)
	respStr, err := httputil.HttpGet(url + param.Params)
	if err != nil {
		return nil, err
	}
	var resp BmsResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (bms *BmsInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(BmsResponse)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.Servers {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.State,
				}, {
					Name:  "bmType",
					Value: strconv.Itoa(d.BmType),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}
