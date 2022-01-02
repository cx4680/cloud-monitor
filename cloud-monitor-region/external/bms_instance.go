package external

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
	"strings"
)

type BmsInstanceService struct {
	service.InstanceServiceImpl
}

type BmsQueryPageRequest struct {
	PageNumber       int    `json:"pageNumber,omitempty"`
	PageSize         int    `json:"pageSize,omitempty"`
	TenantId         string `json:"tenantId,omitempty"`
	AllTenants       string `json:"all_tenants,omitempty"`
	Id               string `json:"id,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

type BmsResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    BmsQueryPageResult `json:"data"`
}

type BmsQueryPageResult struct {
	Servers    []BmsServers `json:"servers"`
	TotalCount int          `json:"total_count"`
}

type BmsServers struct {
	FlavorId string `json:"flavor_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Id       string `json:"id"`
}

type SerialResponse struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Data    SerialQueryPageResult `json:"data"`
}

type SerialQueryPageResult struct {
	Bmcnodes   []Bmcnodes `json:"bmcnodes"`
	TotalCount int        `json:"total_count"`
}

type Bmcnodes struct {
	SerialNumber      string `json:"serial_number"`
	EcsId             string `json:"ecs_id"`
	Manufactory       string `json:"manufactory"`
	Productor         string `json:"productor"`
	Cpu               int    `json:"cpu"`
	CpuManufactory    string `json:"cpu_manufactory"`
	CpuProductor      string `json:"cpu_productor"`
	CpuArch           string `json:"cpu_arch"`
	Memory            int    `json:"memory"`
	SystemDisk        int    `json:"system_disk"`
	DataDisk          int    `json:"data_disk"`
	Nic               string `json:"nic"`
	GpuType           string `json:"gpu_type"`
	FirmwareVersion   string `json:"firmware_version"`
	IdcRoom           string `json:"idc_room"`
	CabinetNum        string `json:"cabinet_num"`
	RankNum           string `json:"rank_num"`
	AvailabilityZone  string `json:"availability_zone"`
	IpmiBmcIp         string `json:"ipmi_bmc_ip"`
	IpmiBmcUsername   string `json:"ipmi_bmc_username"`
	IpmiBmcPassword   string `json:"ipmi_bmc_password"`
	SmartnicManagerIp string `json:"smartnic_manager_ip"`
	Smartniciqn       string `json:"smartniciqn"`
	HealthyStatus     bool   `json:"healthy_status"`
	UnhealthyReason   string `json:"unhealthy_reason"`
	UserSystemStatus  string `json:"user_system_status"`
}

func (bms *BmsInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	var param = "/tenants/" + form.TenantId + "/servers?pageNumber=" + strconv.Itoa(form.Current) + "&pageSize=" + strconv.Itoa(form.PageSize)
	var filterList []string
	if strutil.IsNotBlank(form.InstanceName) {
		filterList = append(filterList, "name:lk:"+form.InstanceName)
	}
	if strutil.IsNotBlank(form.InstanceId) {
		filterList = append(filterList, "id:lk:"+form.InstanceId)
	}
	if strutil.IsNotBlank(form.StatusList) {
		filterList = append(filterList, "status:in:"+form.StatusList)
	}
	if len(filterList) > 0 {
		filter := strings.Join(filterList, "|")
		param += "&filter=" + filter
	}
	return param
}

func (bms *BmsInstanceService) DoRequest(url string, f interface{}) (interface{}, error) {
	var param = f.(string)
	respStr, err := httputil.HttpGet(url + param)
	if err != nil {
		return nil, err
	}
	var resp BmsResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (bms *BmsInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(BmsResponse)
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
					Name:  "serialNumber",
					Value: getSerialNumber(d.Status),
				}},
			})
		}
	}
	return vo.Data.TotalCount, list
}

//获取机器SN号
func getSerialNumber(bmsId string) string {
	p := commonDao.MonitorProduct.GetByAbbreviation(global.DB, BMS)
	url := p.Host + p.PageUrl + "/bmcnodes?filter=ecs_id:eq:" + bmsId
	respStr, err := httputil.HttpGet(url)
	if err != nil {
		logger.Logger().Error("BMS error", err)
		return ""
	}
	var resp SerialResponse
	jsonutil.ToObject(respStr, &resp)
	if resp.Data.TotalCount <= 0 {
		logger.Logger().Error("Not found BMS_ID")
		return ""
	}
	return resp.Data.Bmcnodes[0].SerialNumber
}
