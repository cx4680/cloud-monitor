package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/robfig/cron"
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
			logger.Logger().Info(bt.Name + "start running")
			bt.Task()
			logger.Logger().Info(bt.Name + "running over")
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

func (t *BusinessTaskImpl) Stop() {
	t.c.Stop()
}
