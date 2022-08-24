package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
)

type NatEInstanceService struct {
	service.InstanceServiceImpl
}

type NatEQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"orderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	UserCode string          `json:"userCode,omitempty"`
	IamInfo  service.IamInfo `json:"-"`
}

type NatEQueryParam struct {
	Name      string `json:"name"`
	Uid       string `json:"uid"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	NoVpcName bool   `json:"noVpcName"` //true:优化查询实例列表，但vpcName字段值为空，false:正常查询，vpcName字段有值
}

type NatEResponse struct {
	Code string
	Msg  string
	Data NatEQueryPageResult
}
type NatEQueryPageResult struct {
	Total int
	Rows  []*NatEInfoBean
}

type NatEInfoBean struct {
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

func (nat *NatEInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	queryParam := NatEQueryParam{
		Name:   form.InstanceName,
		Uid:    form.InstanceId,
		Status: form.StatusList,
		Type:   "nat-e",
		//NoVpcName: true,
	}
	return NatEQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		UserCode:  form.TenantId,
	}
}

func (nat *NatEInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := httputil.HttpPostJson(url, form, nil)
	if err != nil {
		return nil, err
	}
	var resp NatEResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (nat *NatEInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(NatEResponse)
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

func (nat *NatEInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
	queryParam := NatEQueryParam{
		Name:   form.InstanceName,
		Uid:    form.InstanceId,
		Status: form.StatusList,
		Type:   "nat-e",
		//NoVpcName: true,
	}
	return NatEQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		UserCode:  form.TenantId,
		IamInfo:   form.IamInfo,
	}
}

func (nat *NatEInstanceService) DoAuthRequest(url string, form interface{}) (interface{}, error) {
	f := form.(NatEQueryPageRequest)
	respStr, err := httputil.HttpPostJson(url, form, nat.GetIamHeader(&f.IamInfo))
	if err != nil {
		return nil, err
	}
	var resp NatEResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (nat *NatEInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(NatEResponse)
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
