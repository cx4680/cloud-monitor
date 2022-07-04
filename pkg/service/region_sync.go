package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"gorm.io/gorm"
)

type RegionSyncService struct {
	dao *dao.RegionSyncDao
}

func NewRegionSyncService() *RegionSyncService {
	return &RegionSyncService{dao.NewRegionSyncDao()}
}

func (s RegionSyncService) GetContactSyncData(time string) (form.ContactSync, error) {
	return s.dao.GetContactSyncData(time)
}

func (s RegionSyncService) ContactSync() error {
	time := s.dao.GetUpdateTime("contact")
	response, err := httputil.HttpGet(config.Cfg.Common.CloudMonitor + "/getContactSyncData?time=" + time.UpdateTime)
	if err != nil {
		logger.Logger().Errorf("同步数据API调用失败：%v", err)
		return err
	}
	var contactSync ContactSyncResponse
	jsonutil.ToObject(response, &contactSync)
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = s.dao.ContactSync(tx, contactSync.Module)
		if err != nil {
			logger.Logger().Errorf("同步失败：%v", err)
			return errors.NewBusinessError("同步失败")
		}
		return nil
	})
	return err
}

func (s RegionSyncService) GetAlarmRuleSyncData(time string) (form.AlarmRuleSync, error) {
	return s.dao.GetAlarmRuleSyncData(time)
}

func (s RegionSyncService) AlarmRuleSync() error {
	time := s.dao.GetUpdateTime("alarmRule")
	response, err := httputil.HttpGet(config.Cfg.Common.CloudMonitor + "/getAlarmRuleSyncData?time=" + time.UpdateTime)
	if err != nil {
		logger.Logger().Errorf("同步数据API调用失败：%v", err)
		return err
	}
	var alarmRuleSync AlarmRuleSyncResponse
	jsonutil.ToObject(response, &alarmRuleSync)
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		tenantList, err := s.dao.AlarmRuleSync(tx, alarmRuleSync.Module)
		if err != nil {
			logger.Logger().Errorf("同步失败：%v", err)
			return errors.NewBusinessError("同步失败")
		}
		if len(tenantList) > 0 {
			prometheusDao := k8s.PrometheusRule
			for _, v := range tenantList {
				prometheusDao.GenerateUserPrometheusRule(v)
			}
		}
		return nil
	})
	return err
}

func (s RegionSyncService) GetAlarmRecordSyncData(time string) (form.AlarmRecordSync, error) {
	return s.dao.GetAlarmRecordSyncData(time)
}

func (s RegionSyncService) PullAlarmRecordSyncData(param form.AlarmRecordSync) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		err := s.dao.PullAlarmRecordSyncData(tx, param)
		if err != nil {
			logger.Logger().Errorf("同步失败：%v", err)
			return errors.NewBusinessError("同步失败")
		}
		return nil
	})
	return err
}

func (s RegionSyncService) AlarmRecordSync() error {
	time := s.dao.GetUpdateTime("alarmRecord").UpdateTime
	currentTime := util.GetNowStr()
	syncData, err := s.dao.GetAlarmRecordSyncData(time)
	if err != nil {
		logger.Logger().Errorf("查询失败：%v", err)
		return err
	}
	response, err := httputil.HttpPostJson(config.Cfg.Common.CloudMonitor+"/pullAlarmRecordSyncData", syncData, nil)
	logger.Logger().Info(response)
	var resp global.Resp
	jsonutil.ToObject(response, &resp)
	if err != nil || resp.ErrorCode != "200" {
		logger.Logger().Errorf("推送数据API调用失败：%v", err)
		return err
	}
	s.dao.UpdateTime(model.SyncTime{Name: "alarmRecord", UpdateTime: currentTime})
	return nil
}

type ContactSyncResponse struct {
	Module form.ContactSync `json:"module"`
}

type AlarmRuleSyncResponse struct {
	Module form.AlarmRuleSync `json:"module"`
}
