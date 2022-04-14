package model

type AlarmInfo struct {
	Id          uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	AlarmBizId  string `gorm:"column:alarm_biz_id;size=50" json:"alarmBizId"`
	Summary     string `gorm:"column:summary;text" json:"summary"`
	Expression  string `gorm:"column:expression;text" json:"expression"`
	ContactInfo string `gorm:"column:contact_info;text" json:"contactInfo"`
}

func (m *AlarmInfo) TableName() string {
	return "t_alarm_info"
}
