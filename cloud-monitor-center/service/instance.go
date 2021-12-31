package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
	"gorm.io/gorm"
)

func UnbindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.UnBindRuleParam)
	instanceDao.UnbindInstance(tx, dto)
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.UnbindRule, param)
}

func BindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.InstanceBindRuleDTO)
	instanceDao.BindInstance(tx, dto)
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.BindRule, param)
}
