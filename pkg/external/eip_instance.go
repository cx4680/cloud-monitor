package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type EipInstanceService struct {
	service.InstanceServiceImpl
}
type EipQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"orderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	UserCode string `json:"userCode,omitempty"`
}

type EipQueryParam struct {
	IpAddress    string `json:"ipAddress,omitempty"`
	Status       []int  `json:"status,omitempty"`
	InstanceType int    `json:"instanceType,omitempty"`
	Uid          string `json:"uid,omitempty"`
	InstanceUid  string `json:"instanceUid,omitempty"`
	RegionCode   string `json:"regionCode,omitempty"`
	UserCode     string `json:"userCode,omitempty"`
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

func (eip *EipInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	queryParam := EipQueryParam{
		IpAddress:   form.ExtraAttr["ip"],
		InstanceUid: form.InstanceId,
		UserCode:    form.TenantId,
	}
	if strutil.IsNotBlank(form.StatusList) {
		queryParam.Status = toIntList(form.StatusList)
	}

	return EipQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		UserCode:  form.TenantId,
	}
}

func (eip *EipInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := httputil.HttpPostJson(url, form, nil)
	if err != nil {
		return nil, err
	}
	var resp EipResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (eip *EipInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EipResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.Uid,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "eipAddress",
					Value: d.IpAddress,
				}, {
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "bandWidth",
					Value: strconv.Itoa(d.BandWidth),
				}, {
					Name:  "bindInstanceId",
					Value: d.InstanceUid,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
