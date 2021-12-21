package messageCenter

type ReceiveType int

const (
	Phone ReceiveType = iota + 1
	Email
)

const (
	SingleResourceThresholdAlarmReminder     = "single_resource_threshold_alarm_reminder"
	SingleResourceThresholdAlarmReminderMail = "single_resource_threshold_alarm_reminder_mail"
	SingleResourceRecoveryReminder           = "single_resource_recovery_reminder"
	SingleResourceRecoveryReminderMail       = "single_resource_recovery_reminder_mail"
	AddAlarmContact                          = "add_alarm_contact"
	AddAlarmContactMail                      = "add_alarm_contact_mail"
	CloudMonitorAlarmSmsLack                 = "cloud_monitor_alarm_sms_lack"
)

type MsgSource int

const (
	ALERT_OPEN   MsgSource = iota + 1 //告警产生
	ALERT_CANCEL                      //告警恢复
	VERIFY                            //验证
	SMS_LACK                          //短信余量
)

type MessageTargetDTO struct {
	Addr string
	Type ReceiveType
}

type MessageSendDTO struct {
	SenderId     string
	Target       []MessageTargetDTO
	SourceType   MsgSource
	Content      string
	MsgEventCode string
}

// SmsMessageReqDTO 消息中心接口入参
type SmsMessageReqDTO struct {
	BusinessId string         `json:"businessId"`
	InModeCode string         `json:"inModeCode"`
	Messages   []MessagesBean `json:"messages"`
	ReferTime  string         `json:"referTime"`
}

// MessagesBean 消息中心接口入参
type MessagesBean struct {
	MsgEventCode   string           `json:"msgEventCode"`
	RecvObjectList []RecvObjectBean `json:"recvObjectList"`
}

// RecvObjectBean 消息中心接口入参
type RecvObjectBean struct {
	RecvObjectType ReceiveType `json:"recvObjectType"`
	RecvObject     string      `json:"recvObject"`
	NoticeContent  string      `json:"noticeContent"`
}
