package models

type AlarmInstance struct {
	AlarmRuleId  string `orm:"alarm_rule_id" json:"alarm_rule_id"`
	InstanceId   string `orm:"instance_id" json:"instance_id"`
	CreateTime   string `orm:"create_time" json:"create_time"`
	RegionCode   string `orm:"region_code" json:"region_code"`
	ZoneCode     string `orm:"zone_code" json:"zone_code"`
	Ip           string `orm:"ip" json:"ip"`
	RegionName   string `orm:"region_name" json:"region_name"`
	ZoneName     string `orm:"zone_name" json:"zone_name"`
	InstanceName string `orm:"instance_name" json:"instance_name"`
	TenantId     string `orm:"tenant_id" json:"tenant_id"`
}

func (*AlarmInstance) TableName() string {
	return "t_alarm_instance"
}
