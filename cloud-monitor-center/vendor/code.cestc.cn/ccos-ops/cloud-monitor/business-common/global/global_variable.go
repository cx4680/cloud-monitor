package global

import "gorm.io/gorm"

var (
	DB               *gorm.DB
	NoticeChannelMap map[string]string
)
