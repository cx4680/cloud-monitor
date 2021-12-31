package model

type NotificationRecord struct {
	Id               string `gorm:"id" json:"id"`
	SenderId         string `gorm:"sender_id" json:"sender_id"`                 // 发送人
	SourceId         string `gorm:"source_id" json:"source_id"`                 // 通知源
	SourceType       int    `gorm:"source_type" json:"source_type"`             // 源类型
	TargetAddress    string `gorm:"target_address" json:"target_address"`       // 通知对象地址
	NotificationType int    `gorm:"notification_type" json:"notification_type"` // 通知类型，1：短信，2：邮箱
	Result           int    `gorm:"result" json:"result"`                       // 通知结果 0:失败1:成功
	CreateTime       string `gorm:"create_time" json:"create_time"`
}

func (*NotificationRecord) TableName() string {
	return "notification_record"
}
