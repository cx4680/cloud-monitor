package model

import "time"

type NotificationRecord struct {
	Id               uint64    `gorm:"id;primary_key;AUTO_INCREMENT" json:"id"`
	BizId            string    `gorm:"column:biz_id;size=50" json:"bizId"`
	SenderId         string    `gorm:"sender_id" json:"senderId"`                 // 发送人
	SourceId         string    `gorm:"source_id" json:"sourceId"`                 // 通知源
	SourceType       uint8     `gorm:"source_type" json:"sourceType"`             // 源类型
	TargetAddress    string    `gorm:"target_address" json:"targetAddress"`       // 通知对象地址
	NotificationType uint8     `gorm:"notification_type" json:"notificationType"` // 通知类型，1：短信，2：邮箱
	Result           uint8     `gorm:"result" json:"result"`                      // 通知结果 0:失败1:成功
	CreateTime       time.Time `gorm:"create_time" json:"createTime"`
}

func (*NotificationRecord) TableName() string {
	return "t_notification_record"
}
