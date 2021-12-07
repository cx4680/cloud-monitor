package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
)

type NatInstanceService struct {
	service.InstanceServiceImpl
}

type NatQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"OrderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	TenantId string
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
	ReleaseTime   string
	OrderId       string
	Specification string
	Description   string
	PayModel      int
	UpdateTime    string
	FreezeTime    string
	UserCode      string
	SubnetName    string
	EipNum        int
	SubnetCidr    string
	Uid           string
	VpcUid        string
	ExpireTime    string
	ResourceCode  string
	StatusList    string
	ExpireDay     int
	CreateTime    string
	SubnetUid     string
	Name          string
	Id            int
	AdminStateUp  bool
	Status        string
}

func (nat *NatInstanceService) convertRealForm(form service.InstancePageForm) interface{} {
	queryParam := NatQueryParam{
		Name: form.InstanceName,
		Uid:  form.InstanceId,
	}
	return NatQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		TenantId:  form.TenantId,
	}
}

func (nat *NatInstanceService) doRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(NatQueryPageRequest)
	respStr, err := tools.HttpPostJson(url, form, map[string]string{"userCode": f.TenantId})
	if err != nil {
		return nil, err
	}
	var resp NatResponse
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (nat *NatInstanceService) convertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(NatResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				Id:   string(rune(d.Id)),
				Name: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "subnetName",
					Value: d.SubnetName,
				}, {
					Name:  "status",
					Value: d.Status,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
