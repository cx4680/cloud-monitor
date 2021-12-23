package messageCenter

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/google/uuid"
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
		logger.Logger().Infof("send msg  is empty, \n")
		return nil
	}
	var list []MessagesBean
	//获取消息模板
	eventCode := AlarmNoticeTemplateMap[GetTemplateMapKey(msgList[0].Type, msgList[0].SourceType)]
	for _, msg := range msgList {
		if msg.Targets == nil || len(msg.Targets) <= 0 {
			logger.Logger().Infof("send msg target is empty, %s\n", tools.ToString(msg))
			continue
		}
		var recvList = make([]RecvObjectBean, len(msg.Targets))
		for i, addr := range msg.Targets {
			recvList[i] = RecvObjectBean{
				RecvObjectType: msg.Type,
				RecvObject:     addr,
				NoticeContent:  msg.Content,
			}
		}
		list = append(list, MessagesBean{
			MsgEventCode:   eventCode,
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
		logger.Logger().Info("send message is empty")
		return nil
	}
	resp, err := tools.HttpPostJson(config.GetCommonConfig().SmsCenterPath, *smsMessageReqDTO, nil)
	if err != nil {
		logger.Logger().Errorf("send message to msgCenter fail:%v", err)
		return err
	}
	logger.Logger().Info(smsMessageReqDTO.Messages)
	logger.Logger().Info("send message to msgCenter resp=" + resp)
	return nil
}
