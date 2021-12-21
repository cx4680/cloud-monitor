package tools

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
)

func ToString(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		logger.Logger().Errorf("序列化json字符串失败, error:%v, data:%v", err, obj)
	}
	return string(str)
}

func ToObject(str string, obj interface{}) {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		logger.Logger().Errorf("反序列化json失败, error:%v, data:%v", err, str)
	}
}
