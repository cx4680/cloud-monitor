package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq/handler"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func AlarmRuleHandler(msgList []*primitive.MessageExt) {
	for i := range msgList {
		logger.Logger().Infof("subscribe callback: %v \n", msgList[i])
		var MqMsg form.MqMsg
		if err := json.Unmarshal(msgList[i].Body, &MqMsg); err != nil {
			continue
		}
		data, err := json.Marshal(MqMsg.Data)
		if err != nil {
			continue
		}
		handler.HandleMsg(MqMsg.EventEum, data, true)
	}
}
