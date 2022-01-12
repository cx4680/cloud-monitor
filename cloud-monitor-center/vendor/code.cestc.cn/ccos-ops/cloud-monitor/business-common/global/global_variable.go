package global

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"gorm.io/gorm"
)

var (
	DB                *gorm.DB
	NoticeChannelList []form.NoticeChannel
)
