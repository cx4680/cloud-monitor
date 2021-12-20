package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestChan(t *testing.T) {
	var waitGroup = &sync.WaitGroup{}
	alertListChan := make(chan []*forms.AlertDTO, 1)
	waitGroup.Add(2)
	go fun1(waitGroup, alertListChan)
	go fun2(waitGroup, alertListChan)
	go func() {
		log.Println("-------------wait ")
		waitGroup.Wait()
		close(alertListChan)
		log.Println("-------------close ")
	}()
	log.Println("ing ")
	var alertList []*forms.AlertDTO
	for list := range alertListChan {
		log.Printf("apend----")
		alertList = append(alertList, list...)
	}
	log.Printf("%+v", alertList)
}

func fun1(wg *sync.WaitGroup, alertListChan chan []*forms.AlertDTO) {
	defer wg.Done()
	dtos := []*forms.AlertDTO{
		{RuleType: "11"},
		{RuleType: "22"},
	}
	alertListChan <- dtos
	log.Println("1-------------end ")
}

func fun2(wg *sync.WaitGroup, alertListChan chan []*forms.AlertDTO) {
	defer wg.Done()
	dtos := []*forms.AlertDTO{
		{RuleType: "33"},
		{RuleType: "44"},
	}
	alertListChan <- dtos
	fmt.Printf("2-------------end ")
}
