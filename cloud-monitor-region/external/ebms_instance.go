package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
	"strings"
)

type EbmsInstanceService struct {
	service.InstanceServiceImpl
}

type EbmsRequest struct {
	TenantId string `json:"tenantId"`
	Params   string `json:"params"`
}

type EbmsResponse struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Data    EbmsQueryPageResult `json:"data"`
}

type EbmsQueryPageResult struct {
	Servers    []EbmsServers `json:"servers"`
	TotalCount int           `json:"total_count"`
}

type EbmsServers struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (ebms *EbmsInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	var params = "?pageNumber=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	var filterList []string
	if strutil.IsNotBlank(form.InstanceName) {
		filterList = append(filterList, "name:lk:"+form.InstanceName)
	}
	if strutil.IsNotBlank(form.InstanceId) {
		filterList = append(filterList, "id:lk:"+form.InstanceId)
	}
	if strutil.IsNotBlank(form.StatusList) {
		filterList = append(filterList, "status:in:"+form.StatusList)
	}
	if len(filterList) > 0 {
		filter := strings.Join(filterList, "|")
		params += "&filter=" + filter
	}
	return EbmsRequest{TenantId: form.TenantId, Params: params}
}

func (ebms *EbmsInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var params = f.(EbmsRequest)
	url = strings.ReplaceAll(url, "{tenantId}", params.TenantId)
	respStr, err := httputil.HttpGet(url + params.Params)
	if err != nil {
		return nil, err
	}
	var resp EbmsResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (ebms *EbmsInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EbmsResponse)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.Servers {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.Id,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.Status,
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}
