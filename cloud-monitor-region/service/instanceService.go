package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
)

func GetInstanceList(productId string) ([]string, error) {
	switch productId {
	case constant.EcsProduct:
		form := service.InstancePageForm{
			TenantId: "210011082310350",
			Current:  1,
			PageSize: 1000,
			Product:  "ecs",
		}
		instanceService := external.ProductInstanceServiceMap[external.ECS]
		page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
		if err != nil {
			return nil, err
		}
		if page.Total <= 0 {
			return nil, nil
		}
		var instanceList []string
		for _, ecsVO := range page.Records.([]service.InstanceCommonVO) {
			if ecsVO.Id != "" {
				instanceList = append(instanceList, ecsVO.Id)
			}
		}
		return instanceList, nil
	case constant.EipProduct:
		form := service.InstancePageForm{
			TenantId: "210011082310350",
			Current:  1,
			PageSize: 1000,
			Product:  "eip",
		}
		instanceService := external.ProductInstanceServiceMap[external.EIP]
		page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
		if err != nil {
			return nil, err
		}
		if page.Total <= 0 {
			return nil, nil
		}
		var instanceList []string
		for _, ecsVO := range page.Records.([]service.InstanceCommonVO) {
			if ecsVO.Id != "" {
				instanceList = append(instanceList, ecsVO.Id)
			}
		}
		return instanceList, nil
	case constant.SlbProduct:
		form := service.InstancePageForm{
			TenantId: "210011082310350",
			Current:  1,
			PageSize: 1000,
			Product:  "slb",
		}
		instanceService := external.ProductInstanceServiceMap[external.SLB]
		page, err := instanceService.GetPage(form, instanceService.(service.InstanceStage))
		if err != nil {
			return nil, err
		}
		if page.Total <= 0 {
			return nil, nil
		}
		var instanceList []string
		for _, ecsVO := range page.Records.([]service.InstanceCommonVO) {
			if ecsVO.Id != "" {
				instanceList = append(instanceList, ecsVO.Id)
			}
		}
		return instanceList, nil
	default:
		return nil, errors.NewBusinessError("产品ID错误")
	}
}
