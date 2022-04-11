package service

import "github.com/toolkits/pkg/container/list"

type AlarmHandlerEvent struct {
	RequestId string
	Type      int
	Data      interface{}
}

var AlarmHandlerQueue = list.NewSafeListLimited(100000)
