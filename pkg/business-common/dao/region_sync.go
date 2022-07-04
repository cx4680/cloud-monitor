package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"gorm.io/gorm"
)

type RegionSyncDao struct {
}

func NewRegionSyncDao() *RegionSyncDao {
	return &RegionSyncDao{}
}

func (dao RegionSyncDao) GetContactSyncData(time string) form.ContactSync {
	currentTime := util.GetNowStr()
	contactSync := form.ContactSync{SyncTime: &model.SyncTime{Name: "contact", UpdateTime: currentTime}}
	global.DB.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&contactSync.Contact)
	global.DB.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&contactSync.ContactGroup)
	var contactList []string
	var contactGroupList []string
	for _, v := range contactSync.Contact {
		contactList = append(contactList, v.BizId)
	}
	for _, v := range contactSync.Contact {
		contactGroupList = append(contactGroupList, v.BizId)
	}
	global.DB.Where("contact_biz_id IN (?) OR group_biz_id in (?)", contactList, contactGroupList).Find(&contactSync.ContactGroupRel)
	global.DB.Where("contact_biz_id IN (?)", contactList).Find(&contactSync.ContactInformation)
	return contactSync
}

func (dao RegionSyncDao) ContactSync(db *gorm.DB, contactSync form.ContactSync) error {
	var contactList []string
	var contactGroupList []string
	if len(contactSync.Contact) != 0 {
		for _, v := range contactSync.Contact {
			contactList = append(contactList, v.BizId)
		}
		db.Save(contactSync.Contact)
	}
	if len(contactSync.ContactGroup) != 0 {
		for _, v := range contactSync.Contact {
			contactGroupList = append(contactGroupList, v.BizId)
		}
		db.Save(contactSync.ContactGroup)
	}
	if len(contactSync.ContactGroupRel) != 0 {
		db.Where("contact_biz_id IN (?) OR group_biz_id in (?)", contactList, contactGroupList).Delete(&model.ContactGroupRel{})
		db.Save(contactSync.ContactGroupRel)
	}
	if len(contactSync.ContactInformation) != 0 {
		db.Where("contact_biz_id IN (?)", contactList).Delete(&model.ContactInformation{})
		db.Save(contactSync.ContactInformation)
	}
	db.Updates(contactSync.SyncTime)
	return nil
}

func (dao RegionSyncDao) GetAlarmRuleSyncData(time string) form.AlarmRuleSync {
	currentTime := util.GetNowStr()
	alarmRuleSync := form.AlarmRuleSync{SyncTime: &model.SyncTime{Name: "alarmRule", UpdateTime: currentTime}}
	global.DB.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRuleSync.AlarmRule)
	var alarmRuleList []string
	for _, v := range alarmRuleSync.AlarmRule {
		alarmRuleList = append(alarmRuleList, v.BizId)
	}
	global.DB.Where("rule_biz_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmItem)
	global.DB.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmNotice)
	global.DB.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRuleSync.ResourceGroup)
	global.DB.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmRuleGroupRel)
	global.DB.Where("create_time > ? AND create_time <= ?", time, currentTime).Find(&alarmRuleSync.AlarmInstance)
	global.DB.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmRuleResourceRel)
	global.DB.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmHandler)
	return alarmRuleSync
}

func (dao RegionSyncDao) AlarmRuleSync(db *gorm.DB, alarmRuleSync form.AlarmRuleSync) error {
	var alarmRuleList []string
	if len(alarmRuleSync.AlarmRule) != 0 {
		for _, v := range alarmRuleSync.AlarmRule {
			alarmRuleList = append(alarmRuleList, v.BizId)
		}
		db.Save(alarmRuleSync.AlarmRule)
	}
	if len(alarmRuleSync.AlarmItem) != 0 {
		db.Where("rule_biz_id IN (?)", alarmRuleList).Delete(&model.AlarmItem{})
		db.Save(alarmRuleSync.AlarmItem)
	}
	if len(alarmRuleSync.AlarmNotice) != 0 {
		db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmNotice{})
		db.Save(alarmRuleSync.AlarmRule)
	}
	if len(alarmRuleSync.ResourceGroup) != 0 {
		db.Save(alarmRuleSync.ResourceGroup)
	}
	if len(alarmRuleSync.AlarmRuleGroupRel) != 0 {
		db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmRuleGroupRel{})
		db.Save(alarmRuleSync.AlarmRuleGroupRel)
	}
	if len(alarmRuleSync.AlarmInstance) != 0 {
		db.Save(alarmRuleSync.AlarmInstance)
	}
	if len(alarmRuleSync.AlarmRuleResourceRel) != 0 {
		db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmRuleResourceRel{})
		db.Save(alarmRuleSync.AlarmRuleResourceRel)
	}
	if len(alarmRuleSync.AlarmHandler) != 0 {
		db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmHandler{})
		db.Save(alarmRuleSync.AlarmHandler)
	}
	db.Updates(alarmRuleSync.SyncTime)
	return nil
}

func (dao RegionSyncDao) GetAlarmRecordSyncData(time string) form.AlarmRecordSync {
	currentTime := util.GetNowStr()
	alarmRecordSync := form.AlarmRecordSync{SyncTime: &model.SyncTime{Name: "alarmRecord", UpdateTime: currentTime}}
	global.DB.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRecordSync.AlarmRecord)
	var alarmRecordList []string
	for _, v := range alarmRecordSync.AlarmRecord {
		alarmRecordList = append(alarmRecordList, v.BizId)
	}
	global.DB.Where("alarm_biz_id IN (?)", alarmRecordList).Find(&alarmRecordSync.AlarmInfo)
	return alarmRecordSync
}

func (dao RegionSyncDao) PullAlarmRecordSyncData(db *gorm.DB, alarmRecordSync form.AlarmRecordSync) error {
	var alarmRecordList []string
	if len(alarmRecordSync.AlarmRecord) != 0 {
		for _, v := range alarmRecordSync.AlarmRecord {
			alarmRecordList = append(alarmRecordList, v.BizId)
		}
		if err := db.Save(alarmRecordSync.AlarmRecord).Error; err != nil {
			return err
		}
	}
	if len(alarmRecordSync.AlarmInfo) != 0 {
		db.Where("alarm_biz_id IN (?)", alarmRecordList).Delete(&model.AlarmInfo{})
		db.Save(alarmRecordSync.AlarmInfo)
	}
	db.Updates(alarmRecordSync.SyncTime)
	return nil
}

func (dao RegionSyncDao) GetUpdateTime() model.SyncTime {
	var entity model.SyncTime
	global.DB.First(&entity)
	return entity
}

func (dao RegionSyncDao) UpdateTime(model model.SyncTime) {
	global.DB.Updates(model)
}
