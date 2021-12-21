package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
)

type AlarmHandlerService struct {
	dao *dao.AlarmHandlerDao
}

func NewAlarmHandlerService() *AlarmHandlerService {
	return &AlarmHandlerService{dao: dao.AlarmHandler}
}

func (svc *AlarmHandlerService) GetAlarmHandlerListByRuleId(ruleId string) []models.AlarmHandler {
	//TODO add cache
	return svc.dao.GetHandlerListByRuleId(global.DB, ruleId)
}
