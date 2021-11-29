package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"github.com/robfig/cron"
	"log"
)

type BusinessTask interface {
	Add(string, func()) error
	Start()
}

type BusinessTaskDTO struct {
	Cron string
	Name string
	Task func()
}

type BusinessTaskImpl struct {
	c *cron.Cron
}

func NewBusinessTaskImpl() *BusinessTaskImpl {
	return &BusinessTaskImpl{c: cron.New()}
}

func (t *BusinessTaskImpl) Add(bt BusinessTaskDTO) error {
	var err error
	if tools.IsEmpty(bt.Name) {
		err = t.c.AddFunc(bt.Cron, bt.Task)
	} else {
		err = t.c.AddFunc(bt.Cron, func() {
			log.Println(bt.Name + "start running")
			bt.Task()
			log.Println(bt.Name + "running over")
		})
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *BusinessTaskImpl) Start() {
	t.c.Start()
	//defer t.c.Stop()
}
