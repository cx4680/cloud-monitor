package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"gorm.io/gorm"
)

func UnbindInstance(tx *gorm.DB, param interface{}, dd interface{}) error {
	instanceDao := dd.(*dao.InstanceDao)
	dto := param.(*forms.UnBindRuleParam)
	instanceDao.UnbindInstance(tx, dto)
	mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.UnbindRule, param)
	return nil
}

func BindInstance(tx *gorm.DB, param interface{}, dd interface{}) error {
	instanceDao := dd.(*dao.InstanceDao)
	dto := param.(*forms.InstanceBindRuleDTO)
	instanceDao.BindInstance(tx, dto)
	mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.BindRule, param)
	return nil
}
