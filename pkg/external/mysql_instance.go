package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type MySqlInstanceService struct {
	service.InstanceServiceImpl
}

type MySqlQueryPageForm struct {
	TenantId        string `json:"tenantId"`
	DbInstanceId    string `json:"dbInstanceId"`
	DbInstanceName  string `json:"dbInstanceName"`
	DbInstanceState int    `json:"dbInstanceState"`
	DbInstanceMode  string `json:"dbInstanceMode"`
	PageNum         int    `json:"pageNum"`
	PageSize        int    `json:"pageSize"`
}

type MySqlQueryPageVO struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"requestId"`
	Data      MySqlPageVO `json:"data"`
}

type MySqlPageVO struct {
	TotalCount  int       `json:"totalCount"`
	TotalPage   int       `json:"totalPage"`
	CurrentPage int       `json:"currentPage"`
	PageSize    int       `json:"pageSize"`
	List        []MySqlVO `json:"list"`
}

type MySqlVO struct {
	DbInstanceId string `json:"dbInstanceId"`
	Name         string `json:"name"`
	State        int    `json:"state"`
	InstanceMode string `json:"instanceMode"`
}

func (mysql *MySqlInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	param := "?tenantId=" + form.TenantId + "&pageNum=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&dbInstanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&dbInstanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.StatusList) {
		param += "&dbInstanceState=" + form.StatusList
	}
	return param
}

func (mysql *MySqlInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(string)
	respStr, err := httputil.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp MySqlQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (mysql *MySqlInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(MySqlQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.List {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.DbInstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.State),
				}, {
					Name:  "instanceMode",
					Value: d.InstanceMode,
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}

func (mysql *MySqlInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
	param := "?tenantId=" + form.TenantId + "&pageNum=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	if strutil.IsNotBlank(form.InstanceName) {
		param += "&dbInstanceName=" + form.InstanceName
	}
	if strutil.IsNotBlank(form.InstanceId) {
		param += "&dbInstanceId=" + form.InstanceId
	}
	if strutil.IsNotBlank(form.StatusList) {
		param += "&dbInstanceState=" + form.StatusList
	}
	return param
}

func (mysql *MySqlInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(string)
	respStr, err := httputil.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp MySqlQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (mysql *MySqlInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(MySqlQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.List {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.DbInstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.State),
				}, {
					Name:  "instanceMode",
					Value: d.InstanceMode,
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}
