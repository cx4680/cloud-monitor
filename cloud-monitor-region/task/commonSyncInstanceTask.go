package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq/producer"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"github.com/pkg/errors"
)

func AddSyncJobs(bt *task.BusinessTaskImpl) error {
	list := dao.MonitorProduct.SelectMonitorProductList()
	for _, product := range *list {
		err := bt.Add(task.BusinessTaskDTO{
			//TODO cron from product
			Cron: "",
			Name: product.Name,
			Task: func() {
				_ = Run(product.Description)
			},
		})
		if err != nil {
			return err
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
			remoteInstanceList, err := getRemoteProductInstanceList(productType)
			if err != nil {
				continue
			}
			sync(tenantId, dbInstanceList, remoteInstanceList)
		}
		current++
	}
	return nil

}

func getRemoteProductInstanceList(productType string) ([]*models.AlarmInstance, error) {
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
	var instances []*models.AlarmInstance
	for current <= totalPage {
		page, err := is.GetPage(commonService.InstancePageForm{Current: current, PageSize: size}, stage)
		if err != nil {
			return nil, err
		}
		if page.Total <= 0 {
			break
		}
		totalPage = page.Total
		vos := page.Records.([]commonService.InstanceCommonVO)
		for _, vo := range vos {
			instances = append(instances, &models.AlarmInstance{InstanceID: vo.Id, InstanceName: vo.Name})
		}
		current++
	}
	return instances, nil
}

func sync(tenantId string, dbList, remoteList []*models.AlarmInstance) {
	if len(remoteList) > 0 {
		syncInstanceName(remoteList)
	}
	deleteNotExistsInstances(tenantId, dbList, remoteList)
}

func syncInstanceName(list []*models.AlarmInstance) {
	dao.AlarmInstance.UpdateBatchInstanceName(list)
	producer.SendInstanceJobMsg(sysRocketMq.InstanceTopic, list)
}

func deleteNotExistsInstances(tenantId string, dbInstanceList []*models.AlarmInstance, instanceInfoList []*models.AlarmInstance) {
	var deletedList []*models.AlarmInstance
	for _, oldInstance := range dbInstanceList {
		exist := false
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
	if len(deletedList) != 0 {
		dao.AlarmInstance.DeleteInstanceList(tenantId, deletedList)
		instance := &dtos.Instance{
			TenantId: tenantId,
			List:     deletedList,
		}
		producer.SendInstanceJobMsg(sysRocketMq.DeleteInstanceTopic, instance)
		service.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	}
}

func IsEqual(A, B interface{}) bool {
	if _, ok := A.(*models.AlarmInstance); ok {
		if _, ok := B.(*models.AlarmInstance); ok {
			if A.(*models.AlarmInstance).InstanceID == B.(*models.AlarmInstance).InstanceID {
				return A.(*models.AlarmInstance).InstanceName == B.(*models.AlarmInstance).InstanceName
			} else {
				return false
			}
		}
		return false
	}
	return false
}
