package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/sync/publisher"
	"gorm.io/gorm"
)

type SyncService interface {
	PersistenceLocal(*gorm.DB, interface{}) (string, error)
	SyncRemote(sys_rocketmq.Topic, string) error
}

type AbstractSyncServiceImpl struct {
}

func (s *AbstractSyncServiceImpl) SyncRemote(topic sys_rocketmq.Topic, msg string) error {
	return publisher.GlobalPublisher.Pub(publisher.PubMessage{
		Topic: topic,
		Data:  msg,
	})
}

func (s *AbstractSyncServiceImpl) Persistence(c SyncService, topic sys_rocketmq.Topic, param interface{}) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		_, err := c.PersistenceLocal(tx, param)
		if err != nil {
			return err
		}
		//if err = s.SyncRemote(topic, content); err != nil {
		//	return err
		//}
		return nil
	})
}

func (s *AbstractSyncServiceImpl) PersistenceInner(db *gorm.DB, c SyncService, topic sys_rocketmq.Topic, param interface{}) error {
	_, err := c.PersistenceLocal(db, param)
	if err != nil {
		return err
	}
	//if err = s.SyncRemote(topic, content); err != nil {
	//	return err
	//}
	return nil
}
