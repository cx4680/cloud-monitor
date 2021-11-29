package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"gorm.io/gorm"
)

type SyncService interface {
	PersistenceLocal(*gorm.DB, interface{}) (string, error)
	SyncRemote(sysRocketMq.Topic, string) error
}

type AbstractSyncServiceImpl struct {
}

func (s *AbstractSyncServiceImpl) SyncRemote(topic sysRocketMq.Topic, msg string) error {
	return sysRocketMq.SendRocketMqMsg(sysRocketMq.RocketMqMsg{
		Topic:   topic,
		Content: msg,
	})
}

func (s *AbstractSyncServiceImpl) Persistence(c SyncService, topic sysRocketMq.Topic, param interface{}) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		content, err2 := c.PersistenceLocal(tx, param)
		if err2 != nil {
			return err2
		}
		if err := s.SyncRemote(topic, content); err != nil {
			return err
		}
		return nil
	})
}

func (s *AbstractSyncServiceImpl) PersistenceInner(db *gorm.DB, c SyncService, topic sysRocketMq.Topic, param interface{}) error {
	content, err := c.PersistenceLocal(db, param)
	if err != nil {
		return err
	}
	if err := s.SyncRemote(topic, content); err != nil {
		return err
	}
	return nil
}
