package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"gorm.io/gorm"
)

func UnbindInstance(tx *gorm.DB, param interface{}, dd interface{}) error {
	instanceDao := dd.(*dao.InstanceDao)
	dto := param.(*forms.UnBindRuleParam)
	instanceDao.UnbindInstance(tx, dto)
	return mq.SendMsg(sysRocketMq.RuleTopic, enums.UnbindRule, param)
}

func BindInstance(tx *gorm.DB, param interface{}, dd interface{}) error {
	instanceDao := dd.(*dao.InstanceDao)
	dto := param.(*forms.InstanceBindRuleDTO)
	instanceDao.BindInstance(tx, dto)
	return mq.SendMsg(sysRocketMq.RuleTopic, enums.BindRule, param)
}
