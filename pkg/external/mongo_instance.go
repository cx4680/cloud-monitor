package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type MongoInstanceService struct {
	service.InstanceServiceImpl
}

type MongoQueryPageForm struct {
	TenantId        string `json:"tenantId"`
	DbInstanceId    string `json:"dbInstanceId"`
	DbInstanceName  string `json:"dbInstanceName"`
	DbInstanceState int    `json:"dbInstanceState"`
	DbInstanceMode  string `json:"dbInstanceMode"`
	PageNum         int    `json:"pageNum"`
	PageSize        int    `json:"pageSize"`
}

type MongoQueryPageVO struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"requestId"`
	Data      MongoPageVO `json:"data"`
}

type MongoPageVO struct {
	TotalCount  int       `json:"totalCount"`
	TotalPage   int       `json:"totalPage"`
	CurrentPage int       `json:"currentPage"`
	PageSize    int       `json:"pageSize"`
	List        []MongoVO `json:"list"`
}

type MongoVO struct {
	DbInstanceId string `json:"dbInstanceId"`
	Name         string `json:"name"`
	State        int    `json:"state"`
}

func (mysql *MongoInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
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

func (mysql *MongoInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(string)
	respStr, err := httputil.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp MongoQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (mysql *MongoInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(MongoQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.List {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.DbInstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.State),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}

func (mysql *MongoInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
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

func (mysql *MongoInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(string)
	respStr, err := httputil.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp MongoQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (mysql *MongoInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(MongoQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.List {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.DbInstanceId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.State),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}
