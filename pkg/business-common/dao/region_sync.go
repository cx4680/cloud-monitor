package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"gorm.io/gorm"
)

type RegionSyncDao struct {
}

func NewRegionSyncDao() *RegionSyncDao {
	return &RegionSyncDao{}
}

func (dao RegionSyncDao) GetContactSyncData(time string) (form.ContactSync, error) {
	currentTime := util.GetNowStr()
	contactSync := form.ContactSync{SyncTime: &model.SyncTime{Name: "contact", UpdateTime: currentTime}}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&contactSync.Contact).Error; err != nil {
			return err
		}
		if err := tx.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&contactSync.ContactGroup).Error; err != nil {
			return err
		}
		var contactList []string
		var contactGroupList []string
		for _, v := range contactSync.Contact {
			contactList = append(contactList, v.BizId)
		}
		for _, v := range contactSync.Contact {
			contactGroupList = append(contactGroupList, v.BizId)
		}
		if err := tx.Where("contact_biz_id IN (?) OR group_biz_id in (?)", contactList, contactGroupList).Find(&contactSync.ContactGroupRel).Error; err != nil {
			return err
		}
		if err := tx.Where("contact_biz_id IN (?)", contactList).Find(&contactSync.ContactInformation).Error; err != nil {
			return err
		}
		if err := tx.Updates(contactSync.SyncTime).Error; err != nil {
			return err
		}
		return nil
	})

	return contactSync, err
}

func (dao RegionSyncDao) ContactSync(db *gorm.DB, contactSync form.ContactSync) error {
	var contactList []string
	var contactGroupList []string
	if len(contactSync.Contact) != 0 {
		for _, v := range contactSync.Contact {
			contactList = append(contactList, v.BizId)
		}
		if err := db.Save(contactSync.Contact).Error; err != nil {
			return err
		}
	}
	if len(contactSync.ContactGroup) != 0 {
		for _, v := range contactSync.Contact {
			contactGroupList = append(contactGroupList, v.BizId)
		}
		if err := db.Save(contactSync.ContactGroup).Error; err != nil {
			return err
		}
	}
	if len(contactSync.ContactGroupRel) != 0 {
		if err := db.Where("contact_biz_id IN (?) OR group_biz_id in (?)", contactList, contactGroupList).Delete(&model.ContactGroupRel{}).Error; err != nil {
			return err
		}
		if err := db.Save(contactSync.ContactGroupRel).Error; err != nil {
			return err
		}
	}
	if len(contactSync.ContactInformation) != 0 {
		if err := db.Where("contact_biz_id IN (?)", contactList).Delete(&model.ContactInformation{}).Error; err != nil {
			return err
		}
		if err := db.Save(contactSync.ContactInformation).Error; err != nil {
			return err
		}
	}
	if err := db.Updates(contactSync.SyncTime).Error; err != nil {
		return err
	}
	return nil
}

func (dao RegionSyncDao) GetAlarmRuleSyncData(time string) (form.AlarmRuleSync, error) {
	currentTime := util.GetNowStr()
	alarmRuleSync := form.AlarmRuleSync{SyncTime: &model.SyncTime{Name: "alarmRule", UpdateTime: currentTime}}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRuleSync.AlarmRule).Error; err != nil {
			return err
		}
		var alarmRuleList []string
		for _, v := range alarmRuleSync.AlarmRule {
			alarmRuleList = append(alarmRuleList, v.BizId)
		}
		if err := tx.Where("rule_biz_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmItem).Error; err != nil {
			return err
		}
		if err := tx.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmNotice).Error; err != nil {
			return err
		}
		if err := tx.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRuleSync.ResourceGroup).Error; err != nil {
			return err
		}
		if err := tx.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmRuleGroupRel).Error; err != nil {
			return err
		}
		if err := tx.Where("create_time > ? AND create_time <= ?", time, currentTime).Find(&alarmRuleSync.AlarmInstance).Error; err != nil {
			return err
		}
		if err := tx.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmRuleResourceRel).Error; err != nil {
			return err
		}
		if err := tx.Where("alarm_rule_id IN (?)", alarmRuleList).Find(&alarmRuleSync.AlarmHandler).Error; err != nil {
			return err
		}
		if err := tx.Updates(alarmRuleSync.SyncTime).Error; err != nil {
			return err
		}
		return nil
	})
	return alarmRuleSync, err
}

func (dao RegionSyncDao) AlarmRuleSync(db *gorm.DB, alarmRuleSync form.AlarmRuleSync) error {
	var alarmRuleList []string
	var tenantList []string
	if len(alarmRuleSync.AlarmRule) != 0 {
		for _, v := range alarmRuleSync.AlarmRule {
			alarmRuleList = append(alarmRuleList, v.BizId)
			tenantList = append(tenantList, v.TenantID)
		}
		if err := db.Save(alarmRuleSync.AlarmRule).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmItem) != 0 {
		if err := db.Where("rule_biz_id IN (?)", alarmRuleList).Delete(&model.AlarmItem{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRuleSync.AlarmItem).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmNotice) != 0 {
		if err := db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmNotice{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRuleSync.AlarmRule).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.ResourceGroup) != 0 {
		if err := db.Save(alarmRuleSync.ResourceGroup).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmRuleGroupRel) != 0 {
		if err := db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmRuleGroupRel{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRuleSync.AlarmRuleGroupRel).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmInstance) != 0 {
		if err := db.Save(alarmRuleSync.AlarmInstance).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmRuleResourceRel) != 0 {
		if err := db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmRuleResourceRel{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRuleSync.AlarmRuleResourceRel).Error; err != nil {
			return err
		}
	}
	if len(alarmRuleSync.AlarmHandler) != 0 {
		if err := db.Where("alarm_rule_id IN (?)", alarmRuleList).Delete(&model.AlarmHandler{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRuleSync.AlarmHandler).Error; err != nil {
			return err
		}
	}
	if err := db.Updates(alarmRuleSync.SyncTime).Error; err != nil {
		return err
	}
	if len(tenantList) > 0 {
		prometheusDao := k8s.PrometheusRule
		for _, v := range tenantList {
			prometheusDao.GenerateUserPrometheusRule(v)
		}
	}
	return nil
}

func (dao RegionSyncDao) GetAlarmRecordSyncData(time string) (form.AlarmRecordSync, error) {
	currentTime := util.GetNowStr()
	alarmRecordSync := form.AlarmRecordSync{SyncTime: &model.SyncTime{Name: "alarmRecord", UpdateTime: currentTime}}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("update_time > ? AND update_time <= ?", time, currentTime).Find(&alarmRecordSync.AlarmRecord).Error; err != nil {
			return err
		}
		var alarmRecordList []string
		for _, v := range alarmRecordSync.AlarmRecord {
			alarmRecordList = append(alarmRecordList, v.BizId)
		}
		if err := tx.Where("alarm_biz_id IN (?)", alarmRecordList).Find(&alarmRecordSync.AlarmInfo).Error; err != nil {
			return err
		}
		if err := tx.Updates(alarmRecordSync.SyncTime).Error; err != nil {
			return err
		}
		return nil
	})
	return alarmRecordSync, err
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
		if err := db.Where("alarm_biz_id IN (?)", alarmRecordList).Delete(&model.AlarmInfo{}).Error; err != nil {
			return err
		}
		if err := db.Save(alarmRecordSync.AlarmInfo).Error; err != nil {
			return err
		}
	}
	if err := db.Updates(alarmRecordSync.SyncTime).Error; err != nil {
		return err
	}
	return nil
}

func (dao RegionSyncDao) GetUpdateTime(name string) model.SyncTime {
	var entity model.SyncTime
	global.DB.Where("name = ?", name).First(&entity)
	return entity
}

func (dao RegionSyncDao) UpdateTime(model model.SyncTime) {
	global.DB.Updates(model)
}
