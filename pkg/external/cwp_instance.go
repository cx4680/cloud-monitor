package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strconv"
)

type CwpInstanceService struct {
	service.InstanceServiceImpl
}

func (cwp *CwpInstanceService) ConvertRealForm(f service.InstancePageForm) interface{} {
	param := EcsQueryPageForm{
		TenantId:     f.TenantId,
		Current:      f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
	}
	if strutil.IsNotBlank(f.StatusList) {
		param.StatusList = toIntList(f.StatusList)
	}
	return param
}

func (cwp *CwpInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	respStr, err := httputil.HttpPostJson(url, f, nil)
	if err != nil {
		return nil, err
	}
	var resp EcsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (cwp *CwpInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EcsQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "osType",
					Value: d.OsType,
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func (cwp *CwpInstanceService) ConvertRealAuthForm(f service.InstancePageForm) interface{} {
	param := EcsQueryPageForm{
		TenantId:     f.TenantId,
		Current:      f.Current,
		PageSize:     f.PageSize,
		InstanceName: f.InstanceName,
		InstanceId:   f.InstanceId,
		IamInfo:      f.IamInfo,
	}
	if strutil.IsNotBlank(f.StatusList) {
		param.StatusList = toIntList(f.StatusList)
	}
	return param
}

func (cwp *CwpInstanceService) DoAuthRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(EcsQueryPageForm)
	respStr, err := httputil.HttpPostJson(url, f, cwp.GetIamHeader(&param.IamInfo))
	if err != nil {
		return nil, err
	}
	var resp EcsQueryPageVO
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (cwp *CwpInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(EcsQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.InstanceId,
				InstanceName: d.InstanceName,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: strconv.Itoa(d.Status),
				}, {
					Name:  "osType",
					Value: d.OsType,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
