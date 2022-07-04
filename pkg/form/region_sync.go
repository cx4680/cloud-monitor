package form

import "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"

type ContactSync struct {
	SyncTime           *model.SyncTime             `json:"syncTime"`
	Contact            []*model.Contact            `json:"contact"`
	ContactGroup       []*model.ContactGroup       `json:"contactGroup"`
	ContactInformation []*model.ContactInformation `json:"contactInformation"`
	ContactGroupRel    []*model.ContactGroupRel    `json:"contactGroupRel"`
}

type AlarmRecordSync struct {
	SyncTime    *model.SyncTime      `json:"syncTime"`
	AlarmRecord []*model.AlarmRecord `json:"alarmRule"`
	AlarmInfo   []*model.AlarmInfo   `json:"alarmInfo"`
}

type AlarmRuleSync struct {
	SyncTime             *model.SyncTime               `json:"syncTime"`
	AlarmRule            []*model.AlarmRule            `json:"alarmRule"`
	AlarmItem            []*model.AlarmItem            `json:"alarmItem"`
	AlarmNotice          []*model.AlarmNotice          `json:"alarmNotice"`
	ResourceGroup        []*model.ResourceGroup        `json:"resourceGroup"`
	AlarmRuleGroupRel    []*model.AlarmRuleGroupRel    `json:"alarmRuleGroupRel"`
	AlarmInstance        []*model.AlarmInstance        `json:"alarmInstance"`
	AlarmRuleResourceRel []*model.AlarmRuleResourceRel `json:"alarmRuleResourceRel"`
	AlarmHandler         []*model.AlarmHandler         `json:"alarmHandler"`
}
