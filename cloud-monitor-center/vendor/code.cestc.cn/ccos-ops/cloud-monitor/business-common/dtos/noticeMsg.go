package dtos

type ReceiveType int

const (
	Phone ReceiveType = iota + 1
	Email
)

type MsgSource int

const (
	ALERT_OPEN   MsgSource = iota + 1 //告警产生
	ALERT_CANCEL                      //告警恢复
	VERIFY                            //验证
	SMS_LACK                          //短信余量
)

type MsgEvent struct {
	Type   ReceiveType
	Source MsgSource
}

type RecvObjectBean struct {
	RecvObjectType ReceiveType
	RecvObject     string
	NoticeContent  string
}

type NoticeMsgDTO struct {
	SourceId      string
	TenantId      string
	MsgEvent      MsgEvent
	RevObjectBean RecvObjectBean
}

type MsgSourceDTO struct {
	Type     MsgSource
	SourceId string
}

type SmsMessageReqDTO struct {
	BusinessId string
	InModeCode string
	Messages   []MessagesBean
	ReferTime  string
}

type MessagesBean struct {
	MsgEventCode   string
	RecvObjectList []RecvObjectBean
}
