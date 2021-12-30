package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"strconv"
	"strings"
)

type BmsInstanceService struct {
	service.InstanceServiceImpl
}

func (bms *BmsInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	var param = "/tenants/" + form.TenantId + "/servers?pageNumber=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	var filterList []string
	if tools.IsNotBlank(form.InstanceName) {
		filterList = append(filterList, "name:lk:"+form.InstanceName)
	}
	if tools.IsNotBlank(form.InstanceId) {
		filterList = append(filterList, "id:lk:"+form.InstanceId)
	}
	if tools.IsNotBlank(form.StatusList) {
		filterList = append(filterList, "status:in:"+form.StatusList)
	}
	if len(filterList) > 0 {
		filter := strings.Join(filterList, "|")
		param += "&filter=" + filter
	}
	return param
}

func (bms *BmsInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	var param = form.(string)
	respStr, err := tools.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp forms.BmsResponse
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (bms *BmsInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(forms.BmsResponse)
	var list []service.InstanceCommonVO
	if vo.Data.TotalCount > 0 {
		for _, d := range vo.Data.Servers {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.Id,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "status",
					Value: d.Status,
				}, {
					Name:  "bmsCpuUsage",
					Value: getMonitorData("bms_cpu_usage", getSerialNumber(d.Id)),
				}, {
					Name:  "bmsMemoryUsage",
					Value: getMonitorData("bms_memory_usage", getSerialNumber(d.Id)),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}

func getSerialNumber(bmsId string) string {
	p := dao.MonitorProduct.GetByAbbreviation(global.DB, constant.Bms)
	url := p.Host + p.PageUrl + "/bmcnodes?filter=ecs_id:eq:" + bmsId
	respStr, err := tools.HttpGet(url)
	if err != nil {
		logger.Logger().Error("BMS error", err)
		return ""
	}
	var resp forms.SerialResponse
	tools.ToObject(respStr, &resp)
	if resp.Data.TotalCount <= 0 {
		logger.Logger().Error("Not found BMS_ID")
		return ""
	}
	return resp.Data.Bmcnodes[0].SerialNumber
}
