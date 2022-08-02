package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"gorm.io/gorm"
)

type AlarmRecordDao struct {
}

var AlarmRecord = new(AlarmRecordDao)

func (dao *AlarmRecordDao) InsertBatch(db *gorm.DB, records []model.AlarmRecord) {
	db.Create(&records)
}

func (dao *AlarmRecordDao) FindAlertRuleBindResourceNum(ruleId, resourceId string) int {
	var num int
	global.DB.Raw("select count(1) from t_alarm_rule_resource_rel t LEFT JOIN t_alarm_rule t2 ON t2.biz_id = t.alarm_rule_id where t.resource_id = ? and t.alarm_rule_id = ? and t2.enabled = 1 and t2.deleted = 0", resourceId, ruleId).Scan(&num)
	return num
}

func (dao *AlarmRecordDao) FindAlertRuleBindGroupNum(ruleId, resourceGroupId string) int {
	var num int
	global.DB.Raw("select count(1) from t_alarm_rule_group_rel t LEFT JOIN t_alarm_rule t2 ON t2.biz_id = t.alarm_rule_id where t.resource_group_id = ?  and t.alarm_rule_id = ? and t2.enabled = 1 and t2.deleted = 0", resourceGroupId, ruleId).Scan(&num)
	return num
}

func (dao *AlarmRecordDao) FindContactInfoByGroupIds(groupIds []string) []*dto.ContactGroupInfo {
	var groupList []*dto.ContactGroupInfo
	global.DB.Raw("SELECT g.biz_id as groupId, g.name as groupName FROM t_contact_group g WHERE g.biz_id IN(?)", groupIds).Scan(&groupList)
	for _, group := range groupList {
		var contacts []dto.UserContactInfo
		global.DB.Raw("select t2.`name` as contactName  ,t2.biz_id as contactId,GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.address END) as phone,GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.address END) as mail from t_contact_group_rel  t LEFT JOIN t_contact t2 on t2.biz_id = t.contact_biz_id LEFT JOIN t_contact_information t3 on ( t3.contact_biz_id = t2.biz_id and  t3.state=1) where t.group_biz_id=?  and t2.`state`=1   GROUP BY contactId  order by contactName", group.GroupId).Scan(&contacts)
		group.Contacts = contacts
	}
	return groupList
}

func (dao *AlarmRecordDao) FindFirstInstanceInfo(instanceId string) *model.AlarmInstance {
	var alarmInstance model.AlarmInstance
	global.DB.Raw("select * from t_resource where instance_id=? limit 1", instanceId).Scan(&alarmInstance)
	//todo 实例是否关联规则 或直接判断规则是否存在
	return &alarmInstance
}

func (dao *AlarmRecordDao) GetLevelTotalByIam(db *gorm.DB, resourcesIdList []string, startTime string, endTime string) []*form.AlarmRecordNum {
	var alarmRecordNum []*form.AlarmRecordNum
	sql := "SELECT level, count(*) AS count  FROM t_alarm_record WHERE source_id IN (?)  AND create_time BETWEEN ?  AND ?  AND rule_source_type != ?  AND status = ? GROUP BY level"
	db.Raw(sql, resourcesIdList, startTime, endTime, source_type.AutoScaling, "firing").Find(&alarmRecordNum)
	return alarmRecordNum
}

func (dao *AlarmRecordDao) GetRecordNumHistoryByIam(db *gorm.DB, resourcesIdList []string, startTime string, endTime string) []vo.RecordNumHistory {
	var sql = "SELECT COUNT(t.id) AS number, " +
		"DATE_FORMAT(t.create_time, '%Y-%m-%d') AS DayTime " +
		"FROM t_alarm_record t " +
		"WHERE source_id IN (?) " +
		" AND t.create_time between ? AND ? " +
		" and t.status='firing' " +
		" GROUP BY daytime "
	var list []vo.RecordNumHistory
	db.Raw(sql, resourcesIdList, startTime, endTime).Find(&list)
	return list
}

func (dao *AlarmRecordDao) GetProductRecordNumHistoryByIam(db *gorm.DB, resourcesIdList []string, startTime string, endTime string) []*form.ProductAlarmRecordNum {
	var productAlarmRecordNum []*form.ProductAlarmRecordNum
	sql := "SELECT t2.abbreviation AS productCode, count(*) AS count FROM t_alarm_record AS t1 LEFT JOIN t_monitor_product AS t2 ON t1.source_type = t2.name WHERE t1.source_id IN (?)  AND t1.create_time BETWEEN ? AND ? AND t1.rule_source_type != ?  AND t1.status = ? GROUP BY t1.source_type ORDER BY count(*) DESC"
	db.Raw(sql, resourcesIdList, startTime, endTime, source_type.AutoScaling, "firing").Find(&productAlarmRecordNum)
	return productAlarmRecordNum
}

func (dao *AlarmRecordDao) GetPageListByIam(db *gorm.DB, productCode string, resourcesIdList []string, startTime string, endTime string, pageNum, pageSize int) ([]*form.AlarmRecordPage, int64) {
	var alarmRecordPage []*form.AlarmRecordPage
	sql := "SELECT t2.abbreviation AS productCode,t1.source_id AS instanceId,t1.rule_name AS ruleName,t1.level AS level,t1.create_time AS time FROM t_alarm_record AS t1 LEFT JOIN t_monitor_product AS t2 ON t1.source_type = t2.name WHERE t1.source_id IN (?)  AND t1.create_time BETWEEN ?  AND ?  AND t1.rule_source_type != ?  AND t1.status = ?  AND t2.abbreviation = ? "
	var total int64
	db.Raw("select count(*) from ("+sql+" ) t ", resourcesIdList, startTime, endTime, source_type.AutoScaling, "firing", productCode).Scan(&total)
	if total == 0 {
		return alarmRecordPage, total
	}
	db.Raw(sql+" ORDER BY t1.create_time DESC LIMIT ?,?", resourcesIdList, startTime, endTime, source_type.AutoScaling, "firing", productCode, (pageNum-1)*pageSize, pageSize).Find(&alarmRecordPage)
	return alarmRecordPage, total
}
