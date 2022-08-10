package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
)

type NatInstanceService struct {
	service.InstanceServiceImpl
}

type NatQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"orderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	UserCode string `json:"userCode,omitempty"`
}

type NatQueryParam struct {
	Name string
	Uid  string
}

type NatResponse struct {
	Code string
	Msg  string
	Data NatQueryPageResult
}
type NatQueryPageResult struct {
	Total int
	Rows  []*NatInfoBean
}

type NatInfoBean struct {
	VpcName       string
	Specification string
	Uid           string
	VpcUid        string
	StatusList    string
	Name          string
	Id            int
	Status        string
	Type          string
}

func (nat *NatInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	queryParam := NatQueryParam{
		Name: form.InstanceName,
		Uid:  form.InstanceId,
	}
	return NatQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		UserCode:  form.TenantId,
	}
}

func (nat *NatInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := httputil.HttpPostJson(url, form, nil)
	if err != nil {
		return nil, err
	}
	var resp NatResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (nat *NatInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(NatResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.Uid,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.Status,
				}, {
					Name:  "specification",
					Value: d.Specification,
				}, {
					Name:  "vpcName",
					Value: d.VpcName,
				}, {
					Name:  "type",
					Value: d.Type,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
