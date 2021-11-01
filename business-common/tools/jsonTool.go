package tools

import (
	"encoding/json"
	"log"
)

func ToString(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln("序列化json字符串失败", err)
	}
	return string(str)
}

func ToObject(str string, obj interface{}) {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Fatalln("反序列化json失败")
	}
}
