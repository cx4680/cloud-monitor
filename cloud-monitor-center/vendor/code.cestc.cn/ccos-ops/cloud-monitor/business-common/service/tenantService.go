package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRedis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"fmt"
	"time"
)

const RedisTimeOutOneHour = 60 * 60

type TenantService struct {
}

func NewTenantService() *TenantService {
	return &TenantService{}
}

func (s *TenantService) GetTenantInfo(tenantId string) dtos.TenantDTO {
	var tenant dtos.TenantDTO
	key := fmt.Sprintf(constants.TenantInfoKey, tenantId)

	value, err := sysRedis.Get(key)
	if err != nil {
		logger.Logger().Error("获取缓存出错, key=" + key)
	}
	if tools.IsNotBlank(value) {
		tools.ToObject(value, &tenant)
		return tenant
	}
	tenantFromRemote := getTenantFromServer(tenantId)
	if tenantFromRemote == nil {
		logger.Logger().Info("获取租户信息为空")
		return tenant
	}
	if e := sysRedis.SetByTimeOut(key, tools.ToString(tenantFromRemote), RedisTimeOutOneHour*time.Second); e != nil {
		logger.Logger().Info("设置redis出错, key=" + key)
	}
	return *tenantFromRemote
}

func getTenantFromServer(tenantId string) *dtos.TenantDTO {
	logger.Logger().Info("tenantId:", tenantId)
	m := make(map[string]string, 1)
	m["loginId"] = tenantId
	resp, err := tools.HttpPostJson(config.GetCommonConfig().TenantUrl, m, nil)
	if err != nil {
		logger.Logger().Error("查询租户信息失败:", resp, err)
		return nil
	}
	var result map[string]map[string]map[string]string
	var loginName, serialNumber string
	tools.ToObject(resp, &result)
	if result["module"] != nil {
		loginName = result["module"]["login"]["loginCode"]
		serialNumber = result["module"]["login"]["serialNumber"]
	} else {
		loginName = result["result"]["login"]["loginCode"]
		serialNumber = result["result"]["login"]["serialNumber"]
	}
	return &dtos.TenantDTO{
		Name:  loginName,
		Phone: serialNumber,
	}
}
