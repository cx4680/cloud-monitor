package data_sync

import "errors"

type Synchronizer interface {
	StartSync() error
}

type CombinedSynchronizer struct {
	Synchronizes []Synchronizer
}

type ContactSynchronizer struct {
}

type AlarmRuleSynchronizer struct {
}

type AlarmRecordSynchronizer struct {
}

func (cs *CombinedSynchronizer) NewCombinedSynchronizer(synchronizes []Synchronizer) (Synchronizer, error) {
	if len(synchronizes) == 0 {
		return nil, errors.New("同步器不能为空")
	}
	return &CombinedSynchronizer{
		Synchronizes: synchronizes,
	}, nil
}

func (cs *CombinedSynchronizer) StartSync() error {
	for _, synchronize := range cs.Synchronizes {
		if err := synchronize.StartSync(); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ContactSynchronizer) StartSync() error {
	return nil
}
func (cs *AlarmRuleSynchronizer) StartSync() error {
	return nil
}
func (cs *AlarmRecordSynchronizer) StartSync() error {
	return nil
}
