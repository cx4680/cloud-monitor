package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
)

func GetInstanceList(productId string, tenantId string) ([]string, error) {
	switch productId {
	case constant.Ecs:
		f := form.EcsQueryPageForm{
			TenantId: tenantId,
			Current:  1,
			PageSize: 10000,
		}
		url := getRequestUrl(constant.Ecs)
		respStr, err := httputil.HttpPostJson(url, f, nil)
		if err != nil {
			return nil, err
		}
		var resp form.EcsQueryPageVO
		jsonutil.ToObject(respStr, &resp)
		var instanceList []string
		for _, ecsVO := range resp.Data.Rows {
			if ecsVO.InstanceId != "" {
				instanceList = append(instanceList, ecsVO.InstanceId)
			}
		}
		return instanceList, nil
	//case constant.Eip:
	//	form := service.InstancePageForm{
	//		TenantId: "210011082310350",
	//		Current:  1,
	//		PageSize: 1000,
	//		Product:  "eip",
	//	}
	//	instanceService := external.ProductInstanceServiceMap[external.EIP]
	//	page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
	//	if err != nil {
	//		return nil, err
	//	}
	//	if page.Total <= 0 {
	//		return nil, nil
	//	}
	//	var instanceList []string
	//	for _, ecsVO := range page.Records.([]service.InstanceCommonVO) {
	//		if ecsVO.Id != "" {
	//			instanceList = append(instanceList, ecsVO.Id)
	//		}
	//	}
	//	return instanceList, nil
	//case constant.Slb:
	//	form := service.InstancePageForm{
	//		TenantId: "210011082310350",
	//		Current:  1,
	//		PageSize: 1000,
	//		Product:  "slb",
	//	}
	//	instanceService := external.ProductInstanceServiceMap[external.SLB]
	//	page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
	//	if err != nil {
	//		return nil, err
	//	}
	//	if page.Total <= 0 {
	//		return nil, nil
	//	}
	//	var instanceList []string
	//	for _, ecsVO := range page.Records.([]service.InstanceCommonVO) {
	//		if ecsVO.Id != "" {
	//			instanceList = append(instanceList, ecsVO.Id)
	//		}
	//	}
	//	return instanceList, nil
	case constant.Bms:
		param := tenantId + "/servers?pageNumber=1&pageSize=10000"
		url := getRequestUrl(constant.Bms)
		respStr, err := httputil.HttpGet(url + param)
		if err != nil {
			return nil, err
		}
		var resp form.BmsResponse
		jsonutil.ToObject(respStr, &resp)

		if resp.Data.TotalCount <= 0 {
			return nil, nil
		}
		var instanceList []string
		for _, bmsVO := range resp.Data.Servers {
			if bmsVO.Id != "" {
				instanceList = append(instanceList, bmsVO.Id)
			}
		}
		return instanceList, nil
	default:
		return nil, errors.NewBusinessError("产品错误")
	}
}

func getRequestUrl(product string) string {
	p := dao.MonitorProduct.GetByAbbreviation(global.DB, product)
	return p.Host + p.PageUrl
}
