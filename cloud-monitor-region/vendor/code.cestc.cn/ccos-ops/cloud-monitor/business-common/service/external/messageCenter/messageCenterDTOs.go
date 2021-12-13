package messageCenter

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

type MessageTargetDTO struct {
	Addr string
	Type ReceiveType
}

type MessageSendDTO struct {
	SenderId   string
	Target     []MessageTargetDTO
	SourceType MsgSource
	Content    string
}

// SmsMessageReqDTO 消息中心接口入参
type SmsMessageReqDTO struct {
	BusinessId string
	InModeCode string
	Messages   []MessagesBean
	ReferTime  string
}

// MessagesBean 消息中心接口入参
type MessagesBean struct {
	MsgEventCode   string
	RecvObjectList []RecvObjectBean
}

// RecvObjectBean 消息中心接口入参
type RecvObjectBean struct {
	RecvObjectType ReceiveType
	RecvObject     string
	NoticeContent  string
}
