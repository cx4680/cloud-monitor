package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"gorm.io/gorm"
)

type AlertRecordDao struct {
}

var AlertRecord = new(AlertRecordDao)

func (dao *AlertRecordDao) InsertBatch(db *gorm.DB, records []model.AlertRecord) {
	db.Create(&records)
}

func (dao *AlertRecordDao) FindAlertRuleBindResourceNum(ruleId, resourceId string) int {
	var num int
	global.DB.Raw("select count(1) from t_alarm_rule_resource_rel t LEFT JOIN t_alarm_rule t2 ON t2.id = t.alarm_rule_id where t.resource_id = ? and t.alarm_rule_id = ? and t2.enabled = 1 and t2.deleted = 0", resourceId, ruleId).Scan(&num)
	return num
}

func (dao *AlertRecordDao) FindAlertRuleBindGroupNum(ruleId, resourceGroupId string) int {
	var num int
	global.DB.Raw("select count(1) from t_alarm_rule_group_rel t LEFT JOIN t_alarm_rule t2 ON t2.id = t.alarm_rule_id where t.resource_group_id = ?  and t.alarm_rule_id = ? and t2.enabled = 1 and t2.deleted = 0", resourceGroupId, ruleId).Scan(&num)
	return num
}

func (dao *AlertRecordDao) FindContactInfoByGroupIds(groupIds []string) []*dto.ContactGroupInfo {
	var groupList []*dto.ContactGroupInfo
	global.DB.Raw("SELECT g.id as groupId, g.name as groupName FROM alert_contact_group g WHERE id IN(?)", groupIds).Scan(&groupList)
	for _, group := range groupList {
		var contacts []dto.UserContactInfo
		global.DB.Raw("select t2.`name` as contactName  ,t2.id as contactId,GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.no  END) as phone,GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.no  END) as mail from alert_contact_group_rel  t LEFT JOIN alert_contact t2 on t2.id = t.contact_id LEFT JOIN alert_contact_information t3 on ( t3.contact_id = t2.id and  t3.is_certify=1) where t.group_id=?  and t2.`status`=1   GROUP BY contactId  order by contactName", group.GroupId).Scan(&contacts)
		group.Contacts = contacts
	}
	return groupList
}

func (dao *AlertRecordDao) FindFirstInstanceInfo(instanceId string) *model.AlarmInstance {
	var alarmInstance model.AlarmInstance
	global.DB.Raw("select * from t_alarm_instance where instance_id=? limit 1", instanceId).Scan(&alarmInstance)
	//todo 实例是否关联规则 或直接判断规则是否存在
	return &alarmInstance
}
