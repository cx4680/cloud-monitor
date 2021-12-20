package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"strconv"
)

type CbrInstanceService struct {
	service.InstanceServiceImpl
}

type CbrQueryParam struct {
	TenantId   string
	VaultId    string
	VaultName  string
	Status     string
	PageNumber string
	PageSize   string
}

type CbrQueryPageVO struct {
	Code        string
	Message     string
	Total_count int
	Data        []CbrInfoBean
}

type CbrInfoBean struct {
	VaultId      string
	TenantId     string
	Name         string
	Type         string
	Status       string
	Region       string
	Zone         string
	Capacity     int
	UsedCapacity int
	CreatedAt    string
	UpdatedAt    string
}

func (c *CbrInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	param := CbrQueryParam{
		TenantId:   form.TenantId,
		VaultId:    form.InstanceId,
		VaultName:  form.InstanceName,
		PageNumber: strconv.Itoa(form.Current),
		PageSize:   strconv.Itoa(form.PageSize),
	}
	if tools.IsNotBlank(form.StatusList) {
		param.Status = form.StatusList
	}
	return param
}

func (c *CbrInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	respStr, err := tools.HttpPostJson(url, form, nil)
	if err != nil {
		logger.Logger().Errorf("error:%v, url:%v, request:%v", err, url, form)
		return nil, err
	}
	var resp CbrQueryPageVO
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (c *CbrInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(CbrQueryPageVO)
	var list []service.InstanceCommonVO
	if vo.Total_count > 0 {
		for _, d := range vo.Data {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.VaultId,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "capacity",
					Value: strconv.Itoa(d.Capacity),
				}, {
					Name:  "usedCapacity",
					Value: strconv.Itoa(d.UsedCapacity),
				}},
			})
		}
	}
	return vo.Total_count, list
}
