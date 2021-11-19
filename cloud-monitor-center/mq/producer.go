package mq

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"encoding/json"
	"fmt"
)

func SendMsg(topic string, eventEum enums.EventEum, module interface{}) error {
	var mqMsg = forms.MqMsg{
		EventEum: eventEum,
		Data:     module,
	}
	// msg对象转json ([]byte)
	jsonBytes, err := json.Marshal(mqMsg)
	if err != nil {
		fmt.Println(err)
	}
	mq.SendMsg(topic, string(jsonBytes))
	return nil
}
