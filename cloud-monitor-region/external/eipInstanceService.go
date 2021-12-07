package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"strconv"
)

type EipInstanceService struct {
	service.InstanceServiceImpl
}
type EipQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"OrderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	TenantId string
}

type EipQueryParam struct {
	IpAddress    string `json:"ipAddress,omitempty"`
	Status       int    `json:"status,omitempty"`
	InstanceType int    `json:"instanceType,omitempty"`
	Uid          string `json:"uid,omitempty"`
	RegionCode   string `json:"regionCode,omitempty"`
	StatusList   []int  `json:"statusList,omitempty"`
}

type EipResponse struct {
	Code string
	Msg  string
	Data EipQueryPageResult
}
type EipQueryPageResult struct {
	Total int
	Rows  []EipInfoBean
}

type EipInfoBean struct {
	BandWidth     int         `json:"bandWidth"`
	ReleaseTime   string      `json:"releaseTime"`
	Instance      interface{} `json:"instance"`
	Data          string      `json:"data"`
	Isp           string      `json:"isp"`
	Description   interface{} `json:"description"`
	PayModel      int         `json:"payModel"`
	FreezeTime    string      `json:"freezeTime"`
	UserCode      string      `json:"userCode"`
	InstanceUid   string      `json:"instanceUid"`
	Uid           string      `json:"uid"`
	RegionCode    string      `json:"regionCode"`
	ResourceCode  string      `json:"resourceCode"`
	ExpireDay     int         `json:"expireDay"`
	Id            int         `json:"id"`
	BandWidthUid  string      `json:"bandWidthUid"`
	FloatingIp    string      `json:"floatingIp"`
	InstanceType  interface{} `json:"instanceType"`
	IpAddress     string      `json:"ipAddress"`
	FloatingIpUid interface{} `json:"floatingIpUid"`
	PortUid       interface{} `json:"portUid"`
	ExpireTime    string      `json:"expireTime"`
	CreateTime    string      `json:"createTime"`
	Name          string      `json:"name"`
	Status        int         `json:"status"`
}

func (eip *EipInstanceService) convertRealForm(form service.InstancePageForm) interface{} {
	queryParam := EipQueryParam{
		IpAddress: form.ExtraAttr["ip"],
		Uid:       form.InstanceId,
	}
	if form.StatusList != nil && len(form.StatusList) > 0 {
		var list []int
		for _, s := range form.StatusList {
			i, err := strconv.Atoi(s)
			if err != nil {
				list = append(list, i)
			}
		}
		queryParam.StatusList = list
	}

	return EipQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		TenantId:  form.TenantId,
	}
}

func (eip *EipInstanceService) doRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(EipQueryPageRequest)
	respStr, err := tools.HttpPostJson(url, form, map[string]string{"userCode": f.TenantId})
	if err != nil {
		return nil, err
	}
	var resp EipResponse
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (eip *EipInstanceService) convertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EipResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				Id:   d.InstanceUid,
				Name: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "ipAddress",
					Value: d.IpAddress,
				}, {
					Name:  "status",
					Value: string(rune(d.Status)),
				}},
			})
		}
	}
	return vo.Data.Total, list
}
