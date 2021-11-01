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
	RecvObjectType int
	RecvObject     string
	NoticeContent  string
}

type NoticeMsg struct {
	SourceId      string
	TenantId      string
	MsgEvent      MsgEvent
	RevObjectBean RecvObjectBean
}
