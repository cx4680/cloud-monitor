package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"fmt"
	"strconv"
)

type DnsInstanceService struct {
	service.InstanceServiceImpl
}

type DnsQueryPageForm struct {
	TenantId     string          `json:"tenantId"`
	InstanceId   string          `json:"instanceId"`
	InstanceName string          `json:"instanceName"`
	PageNumber   int             `json:"pageNumber"`
	PageSize     int             `json:"pageSize"`
	IamInfo      service.IamInfo `json:"-"`
	Status       string          `json:"statusList"`
}

type DnsQueryPageVO struct {
	Code string  `json:"code"`
	Msg  string  `json:"message"`
	Data DnsData `json:"data"`
}

type DnsData struct {
	Total   int          `json:"totalCount"`
	Records []DnsRecords `json:"list"`
}

type DnsRecords struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	TenantId     string `json:"tenantId"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func (dns *DnsInstanceService) ConvertRealForm(f service.InstancePageForm) interface{} {
	param := DnsQueryPageForm{
		TenantId:     f.TenantId,
		PageNumber:   f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
		Status:       f.StatusList,
	}
	return param
}

func (dns *DnsInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var form = f.(DnsQueryPageForm)
	param := "?pageNumber=" + strconv.Itoa(form.PageNumber) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&instanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&instanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.TenantId) {
		param += "&tenantId=" + form.TenantId
	}
	if strutil.IsNotBlank(form.Status) {
		param += "&status=" + form.Status
	}

	respStr, err := httputil.HttpGet(url + param)
	fmt.Println("dns接口：", url+param)
	fmt.Println("dns接口-respStr：", respStr)
	if err != nil {
		return nil, err
	}
	var resp DnsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (dns *DnsInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(DnsQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: convertStatus(d.Status),
				}, {
					Name:  "tenantId",
					Value: d.TenantId,
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func (dns *DnsInstanceService) ConvertRealAuthForm(f service.InstancePageForm) interface{} {
	param := DnsQueryPageForm{
		TenantId:     f.TenantId,
		PageNumber:   f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
		IamInfo:      f.IamInfo,
	}
	return param
}

func (dns *DnsInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var form = f.(DnsQueryPageForm)
	param := "?pageNumber=" + strconv.Itoa(form.PageNumber) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&instanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&instanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.TenantId) {
		param += "&tenantId=" + form.TenantId
	}
	respStr, err := httputil.HttpGet(url + param)
	fmt.Println("dns-im接口：", url+param)

	if err != nil {
		return nil, err
	}
	var resp DnsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (dns *DnsInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(DnsQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: convertStatus(d.Status),
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func convertStatus(status string) string {
	res := "0"
	if status == "running" {
		res = "2"
	} else if status == "deleting" {
		res = "3"
	} else if status == "creating" {
		res = "1"
	}
	return res
}
