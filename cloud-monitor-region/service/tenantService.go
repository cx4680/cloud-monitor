package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"log"
	"time"
)

//TODO 从环境变量中获取
var tenantUrl string = ""

const RedisTimeOutOneHour = 60 * 60

type TenantService struct {
}

func (s *TenantService) GetTenantInfo(tenantId string) dtos.TenantDTO {
	var tenant dtos.TenantDTO

	key := "userInfo:" + tenantId
	value, err := redis.Get(key)
	if err != nil {
		log.Fatalln("获取缓存出错, key=" + key)
	}
	if tools.IsNotBlank(value) {
		//
		tools.ToObject(value, &tenant)
	} else {

		err := redis.SetByTimeOut(key, "", RedisTimeOutOneHour*time.Second)
		if err != nil {
			log.Fatalln("设置redis出错, key=" + key)
		}
	}

	return tenant
}

func getTenant(tenantId string) dtos.TenantDTO {
	//TODO http request
	return dtos.TenantDTO{}
}
