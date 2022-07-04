package model

type SyncTime struct {
	Name       string `gorm:"column:name;primary_key" json:"name" `
	UpdateTime string `gorm:"column:update_time" json:"update_time"`
}

func (*SyncTime) TableName() string {
	return "t_sync_time"
}
