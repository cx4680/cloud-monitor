package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/region"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"github.com/pkg/errors"
	"log"
)

func AddSyncJobs(bt *task.BusinessTaskImpl) error {
	list := dao.MonitorProduct.GetMonitorProduct()
	for _, product := range *list {
		if strutil.IsNotBlank(product.Cron) {
			abbreviation := product.Abbreviation
			err := bt.Add(task.BusinessTaskDTO{
				Cron: product.Cron,
				Name: product.Name,
				Task: func() {
					_ = Run(abbreviation)
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Run(productType string) error {
	var (
		current   = 1
		size      = 100
		totalPage = 1
	)

	for current <= totalPage {
		tenantPage := dao.AlarmInstance.SelectTenantIdList(productType, current, size)
		if tenantPage.Total <= 0 {
			break
		}
		totalPage = tenantPage.Pages
		tenantIds := tenantPage.Records.(*[]string)
		for _, tenantId := range *tenantIds {
			dbInstanceList := dao.AlarmInstance.SelectInstanceList(tenantId, productType)
			remoteInstanceList, err := GetRemoteProductInstanceList(productType, tenantId)
			logger.Logger().Infof(" sync list ,db: %+v,remote:%+v", dbInstanceList, remoteInstanceList)
			if err != nil {
				logger.Logger().Error("查询出错", err)
				continue
			}
			sync(tenantId, dbInstanceList, remoteInstanceList)
		}
		current++
	}
	return nil

}

func GetRemoteProductInstanceList(productType string, tenantId string) ([]*model.AlarmInstance, error) {
	var is = external.ProductInstanceServiceMap[productType]
	if is == nil {
		return nil, errors.New("未配置instanceService")
	}
	var (
		current   = 1
		size      = 100
		totalPage = 1
	)

	stage := is.(commonService.InstanceStage)
	var instances []*model.AlarmInstance
	for current <= totalPage {
		page, err := is.GetPage(commonService.InstancePageForm{Current: current, PageSize: size, Product: productType, TenantId: tenantId}, stage)
		if err != nil {
			return nil, err
		}
		if page.Total <= 0 {
			break
		}
		totalPage = page.Pages
		vos := page.Records.([]commonService.InstanceCommonVO)
		for _, vo := range vos {
			instances = append(instances, &model.AlarmInstance{InstanceID: vo.InstanceId, InstanceName: vo.InstanceName})
		}
		current++
	}
	return instances, nil
}

func sync(tenantId string, dbList, remoteList []*model.AlarmInstance) {
	if len(remoteList) > 0 {
		syncInstanceName(remoteList)
	}
	deleteNotExistsInstances(tenantId, dbList, remoteList)
}

func syncInstanceName(list []*model.AlarmInstance) {
	for _, instance := range list {
		instance.RegionCode = config.Cfg.Common.RegionName
		info := GetRegionInfo(instance.RegionCode)
		instance.RegionName = info.Name
	}
	dao.AlarmInstance.UpdateBatchInstanceName(list)

	producer.SendInstanceJobMsg(sys_rocketmq.InstanceTopic, list)
}

func deleteNotExistsInstances(tenantId string, dbInstanceList []*model.AlarmInstance, instanceInfoList []*model.AlarmInstance) {
	var deletedList []*model.AlarmInstance
	for _, oldInstance := range dbInstanceList {
		exist := false
		if oldInstance.RegionCode != config.Cfg.Common.RegionName {
			continue
		}
		for _, newInstance := range instanceInfoList {
			if IsEqual(oldInstance, newInstance) {
				exist = true
				break
			}
		}
		if !exist {
			deletedList = append(deletedList, oldInstance)
		}
	}
	logger.Logger().Infof(" delete list :%+v", deletedList)
	if len(deletedList) != 0 {
		dao.AlarmInstance.DeleteInstanceList(tenantId, deletedList)
		instance := &dto.Instance{
			TenantId: tenantId,
			List:     deletedList,
		}
		producer.SendInstanceJobMsg(sys_rocketmq.DeleteInstanceTopic, instance)
		service.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	}
}

func IsEqual(A, B interface{}) bool {
	if _, ok := A.(*model.AlarmInstance); ok {
		if _, ok := B.(*model.AlarmInstance); ok {
			return A.(*model.AlarmInstance).InstanceID == B.(*model.AlarmInstance).InstanceID
		}
	}
	return false
}

func GetRegionInfo(regionCode string) region.ConfigItemVO {
	regionInfo, err2 := sys_redis.Get(regionCode)
	if err2 != nil {
		log.Printf("获取region缓存错误, key=%v,err:%v", regionCode, err2)
	}
	vo := region.ConfigItemVO{}
	if regionInfo != "" {
		jsonutil.ToObject(regionInfo, &vo)
		return vo
	}
	list, err := region.RegionService.GetRegionList("1")
	if err != nil {
		log.Printf("查询region list 错误, key=%v,err:%v", regionCode, err2)
		return vo
	}
	for _, vo := range list {
		if vo.Code == regionCode {
			if e := sys_redis.Set(regionCode, jsonutil.ToString(vo)); e != nil {
				log.Printf("缓存region错误, key=%v,err:%v", regionCode, err2)
			}
			return vo
		}
	}
	return vo
}
