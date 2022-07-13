package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"gorm.io/gorm"
)

type SyncService interface {
	PersistenceLocal(*gorm.DB, interface{}) (string, error)
}

type AbstractSyncServiceImpl struct {
}

func (s *AbstractSyncServiceImpl) Persistence(c SyncService, param interface{}) error {
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

func (s *AbstractSyncServiceImpl) PersistenceInner(db *gorm.DB, c SyncService, param interface{}) error {
	_, err := c.PersistenceLocal(db, param)
	if err != nil {
		return err
	}
	//if err = s.SyncRemote(topic, content); err != nil {
	//	return err
	//}
	return nil
}
