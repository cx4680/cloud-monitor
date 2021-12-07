package messageCenter

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"github.com/google/uuid"
	"log"
	"strings"
)

const AppCode = "monitor"

type Service struct {
}

func NewService() *Service {
	return new(Service)
}

// Send 一次发送一个事件通知，可通知多个联系人
func (s *Service) Send(msg MessageSendDTO) error {
	return s.doSend(s.buildReq([]MessageSendDTO{msg}))
}

func (s *Service) SendBatch(msgList []MessageSendDTO) error {
	return s.doSend(s.buildReq(msgList))
}

func (s *Service) buildReq(msgList []MessageSendDTO) (req *SmsMessageReqDTO) {
	if len(msgList) <= 0 {
		log.Printf("send msg  is empty, \n")
		return nil
	}
	var list []MessagesBean
	for _, msg := range msgList {
		if msg.Target == nil || len(msg.Target) <= 0 {
			log.Printf("send msg target is empty, %s\n", tools.ToString(msg))
			continue
		}
		var recvList []RecvObjectBean
		for _, t := range msg.Target {
			recvList = append(recvList, RecvObjectBean{
				RecvObjectType: t.Type,
				RecvObject:     t.Addr,
				NoticeContent:  msg.Content,
			})
		}
		list = append(list, MessagesBean{
			MsgEventCode:   string(rune(msg.SourceType)),
			RecvObjectList: recvList,
		})

	}
	req = &SmsMessageReqDTO{
		BusinessId: strings.ReplaceAll(uuid.New().String(), "-", ""),
		InModeCode: AppCode,
		Messages:   list,
		ReferTime:  "",
	}
	return
}

func (s *Service) doSend(smsMessageReqDTO *SmsMessageReqDTO) error {
	if smsMessageReqDTO == nil {
		log.Println("send message is empty")
		return nil
	}
	resp, err := tools.HttpPostJson(config.GetCommonConfig().SmsCenterPath, *smsMessageReqDTO, nil)
	if err != nil {
		log.Fatal("send message to msgCenter fail", err)
		return err
	}
	log.Println("send message to msgCenter resp=" + resp)
	return nil
}
