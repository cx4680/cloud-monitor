package external

import commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"

//已接入的产品简称
const (
	ECS        = "ecs"
	CBR        = "cbr"
	EIP        = "eip"
	NAT        = "nat"
	SLB        = "slb"
	BMS        = "bms"
	EBMS       = "ebms"
	MYSQL      = "mysql"
	DM         = "dm"
	POSTGRESQL = "postgresql"
	KAFKA      = "kafka"
)

var ProductInstanceServiceMap = map[string]commonService.InstanceService{
	CBR: &CbrInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	ECS: &EcsInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	EIP: &EipInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	NAT: &NatInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	SLB: &SlbInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	BMS: &BmsInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	EBMS: &EbmsInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	MYSQL: &MySqlInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	DM: &MySqlInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	POSTGRESQL: &MySqlInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	KAFKA: &KafkaInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
}
