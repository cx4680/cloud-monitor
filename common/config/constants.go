package config

const (
	PublicCloud      = "publicCloud"      // 公有云
	ProprietaryCloud = "proprietaryCloud" // 专有云
)

//消息中心是否开启配置常量
const (
	MsgOpen  = "open"
	MsgClose = "close"
)

//消息渠道列表
const (
	// MsgChannelEmail 邮件
	MsgChannelEmail = "mail"
	// MsgChannelSms 短信
	MsgChannelSms = "sms"
	// MsgChannelStation 站内信
	MsgChannelStation = "station"
	// MsgChannelLX 蓝信
	MsgChannelLX = "lx"
)
