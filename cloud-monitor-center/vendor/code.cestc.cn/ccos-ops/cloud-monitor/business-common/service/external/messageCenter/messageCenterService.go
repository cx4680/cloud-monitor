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

func checkChannelEnable(channel string, supportChannels []string) bool {
	for _, c := range supportChannels {
		if channel == c {
			return true
		}
	}
	return false
}

func (s *Service) filter(msg MessageSendDTO, channelArr []string) bool {
	if msg.Targets == nil || len(msg.Targets) <= 0 {
		logger.Logger().Info("send msg target is empty, ", tools.ToString(msg))
		return false
	}

	if msg.Type == Phone && !checkChannelEnable(config.MsgChannelSms, channelArr) {
		logger.Logger().Info("message center not support this channel, channel=", config.MsgChannelSms)
		return false
	}
	if msg.Type == Email && !checkChannelEnable(config.MsgChannelEmail, channelArr) {
		logger.Logger().Info("message center not support this channel, channel=", config.MsgChannelEmail)
		return false
	}
	//TODO check other channel
	return true
}

func (s *Service) buildReq(msgList []MessageSendDTO) (req *SmsMessageReqDTO) {
	if config.MsgOpen != config.GetCommonConfig().MsgIsOpen {
		logger.Logger().Info("this env message center is disable")
		return nil
	}
	channels := config.GetCommonConfig().MsgChannel
	if tools.IsBlank(channels) {
		logger.Logger().Info("this env message channels is empty")
		return nil
	}
	channelArr := strings.Split(channels, ",")
	if len(msgList) <= 0 {
		logger.Logger().Info("send msg  is empty")
		return nil
	}

	var list []MessagesBean
	//获取消息模板
	for _, msg := range msgList {
		if !s.filter(msg, channelArr) {
			logger.Logger().Info("send msg target is empty, ", tools.ToString(msg))
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
			MsgEventCode:   NoticeTemplateMap[GetTemplateMapKey(msg.Type, msg.SourceType)],
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
	if config.GetCommonConfig().MsgIsOpen != config.MsgOpen {
		logger.Logger().Info("message center is not open")
		return nil
	}
	if smsMessageReqDTO == nil {
		logger.Logger().Info("send message is empty")
		return nil
	}
	if tools.IsBlank(config.GetCommonConfig().SmsCenterPath) {
		logger.Logger().Info("message center path config error, it's empty")
		return nil
	}
	logger.Logger().Info(smsMessageReqDTO.Messages)
	resp, err := tools.HttpPostJson(config.GetCommonConfig().SmsCenterPath, *smsMessageReqDTO, nil)
	if err != nil {
		logger.Logger().Error("send message to msgCenter fail, ", err)
		return err
	}
	logger.Logger().Info("send message to msgCenter resp=" + resp)
	return nil
}
