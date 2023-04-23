package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type CgwInstanceService struct {
	service.InstanceServiceImpl
}

type CgwQueryPageForm struct {
	TenantId     string          `json:"tenantId"`
	InstanceId   string          `json:"instanceId"`
	InstanceName string          `json:"instanceName"`
	PageNum      int             `json:"pageNum"`
	PageSize     int             `json:"pageSize"`
	Status       string          `json:"status"`
	IamInfo      service.IamInfo `json:"-"`
}

type CgwQueryPageVO struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data CgwData `json:"data"`
}

type CgwData struct {
	Total    int          `json:"total"`
	PageNum  int          `json:"pageNum"`
	PageSize int          `json:"pageSize"`
	Records  []CgwRecords `json:"records"`
}

type CgwRecords struct {
	PaasInstanceId string `json:"paasInstanceId"`
	InstanceId     string `json:"instanceId"`
	InstanceName   string `json:"instanceName"`
	InstanceSpec   string `json:"instanceSpec"`
	InstanceType   int    `json:"instanceType"`
	Status         int    `json:"status"`
	VpcInfo        string `json:"vpcInfo"`
	Eip            string `json:"eip"`
}

func (cgw *CgwInstanceService) ConvertRealForm(f service.InstancePageForm) interface{} {
	param := CgwQueryPageForm{
		TenantId:     f.TenantId,
		PageNum:      f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
		Status:       f.StatusList,
	}
	return param
}

func (cgw *CgwInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var form = f.(CgwQueryPageForm)
	param := "?pageNum=" + strconv.Itoa(form.PageNum) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&instanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&paasInstanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.Status) {
		param += "&status=" + form.Status
	}
	respStr, err := httputil.HttpHeaderGet(url+param, map[string]string{"CECLOUD-CSP-USER": "{\"tenantId\":\"" + form.TenantId + "\"}"})
	if err != nil {
		return nil, err
	}
	var resp CgwQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (cgw *CgwInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(CgwQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.PaasInstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "instanceSpec",
					Value: d.InstanceSpec,
				}, {
					Name:  "instanceType",
					Value: strconv.Itoa(d.InstanceType),
				}, {
					Name:  "vpcInfo",
					Value: d.VpcInfo,
				}, {
					Name:  "eip",
					Value: d.Eip,
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func (cgw *CgwInstanceService) ConvertRealAuthForm(f service.InstancePageForm) interface{} {
	param := CgwQueryPageForm{
		TenantId:     f.TenantId,
		PageNum:      f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
		Status:       f.StatusList,
		IamInfo:      f.IamInfo,
	}
	return param
}

func (cgw *CgwInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var form = f.(CgwQueryPageForm)
	param := "?pageNum=" + strconv.Itoa(form.PageNum) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&instanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&paasInstanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.Status) {
		param += "&status=" + form.Status
	}
	headerMap := map[string]string{"CECLOUD-CSP-USER": "{\"tenantId\":\"" + form.TenantId + "\",\"userId\":\"" + form.IamInfo.UserId + "\",\"userName\":\"" + form.IamInfo.UserName + "\",\"userType\":\"" + form.IamInfo.UserType + "\",\"roleName\":\"" + form.IamInfo.CloudAccountOrganizeRoleName + "\",\"organizationRoleName\":\"" + form.IamInfo.OrganizeAssumeRoleName + "\"}"}
	respStr, err := httputil.HttpHeaderGet(url+param, headerMap)
	if err != nil {
		return nil, err
	}
	var resp CgwQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (cgw *CgwInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(CgwQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Records {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.PaasInstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "instanceSpec",
					Value: d.InstanceSpec,
				}, {
					Name:  "instanceType",
					Value: strconv.Itoa(d.InstanceType),
				}, {
					Name:  "vpcInfo",
					Value: d.VpcInfo,
				}, {
					Name:  "eip",
					Value: d.Eip,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
