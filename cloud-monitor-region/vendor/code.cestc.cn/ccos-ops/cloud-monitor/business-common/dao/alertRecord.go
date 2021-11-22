package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
)

type AlertRecordDao struct {
}

var AlertRecord = new(AlertRecordDao)


func (dao *AlertRecordDao) InsertBatch(records []*models.AlertRecord) {
	var list []models.AlertRecord
	for _, record := range records {
		list = append(list, *record)
	}
	database.GetDb().Create(&records)
}

func (dao *AlertRecordDao) FindAlertRuleBindNum(ruleId, instanceId string) int {
	var num int
	database.GetDb().Raw("select count(1) from t_alarm_instance t LEFT JOIN t_alarm_rule t2 ON t2.id = t.alarm_rule_id where t.instance_id = ? and t.alarm_rule_id = ? and t2.enabled = 1 and t2.deleted = 0", instanceId, ruleId).Scan(&num)
	return num
}

func (dao *AlertRecordDao) FindContactInfoByGroupIds(groupIds []string) []dtos.ContactGroupInfo {
	var groupList []dtos.ContactGroupInfo
	database.GetDb().Raw("SELECT g.id as groupId, g.name as groupName FROM alert_contact_group g WHERE id IN(?)", groupIds).Scan(&groupList)
	for _, group := range groupList {
		var contacts []dtos.UserContactInfo
		database.GetDb().Raw("select t2.`name` as contactName  ,t2.id as contactId,GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.no  END) as phone,GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.no  END) as mail from alert_contact_group_rel  t LEFT JOIN alert_contact t2 on t2.id = t.contact_id LEFT JOIN alert_contact_information t3 on ( t3.contact_id = t2.id and  t3.is_certify=1) where t.group_id=?  and t2.`status`=1   GROUP BY contactId  order by contactName", group.GroupId).Scan(&contacts)
		group.Contacts = contacts
	}
	return groupList
}

func (dao *AlertRecordDao) FindFirstInstanceInfo(instanceId string) *models.AlarmInstance {
	var alarmInstance models.AlarmInstance
	database.GetDb().Raw("select * from t_alarm_instance where instance_id=? limit 1", instanceId).Scan(&alarmInstance)
	return &alarmInstance
}
