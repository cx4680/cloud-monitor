package models

import "time"

type ResourceGroup struct {
	Id           string           `gorm:"column:id;primary_key"`
	Name         string           `gorm:"column:name"`
	TenantId     string           `gorm:"column:tenant_id"`
	ProductId    int              `gorm:"column:product_id"`
	SourceType   int              `gorm:"column:source"` // 来源:1 用户 2 弹性伸缩
	CreateTime   time.Time        `gorm:"column:create_time;autoCreateTime;type:datetime"`
	UpdateTime   time.Time        `gorm:"column:update_time;autoUpdateTime;type:datetime"`
	ResourceList []*AlarmInstance `gorm:"many2many:t_resource_resource_group_rel;joinForeignKey:ResourceGroupId;joinReferences:ResourceId"`
}

func (m *ResourceGroup) TableName() string {
	return "t_resource_group"
}

type AlarmRuleGroupRel struct {
	Id              uint64    `gorm:"column:id;primary_key;autoIncrement"`
	AlarmRuleId     string    `gorm:"column:alarm_rule_id"`
	ResourceGroupId string    `gorm:"column:resource_group_id"`
	CalcMode        int       `gorm:"column:calc_mode"` // 计算模式：1 单个实例 2 按组
	TenantId        string    `gorm:"column:tenant_id"`
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime;type:datetime"` // 创建时间
}

func (m *AlarmRuleGroupRel) TableName() string {
	return "t_alarm_rule_group_rel"
}

type AlarmRuleResourceRel struct {
	Id          uint64    `gorm:"column:id;primary_key;autoIncrement"`
	AlarmRuleId string    `gorm:"column:alarm_rule_id;"`
	ResourceId  string    `gorm:"column:resource_id;"`
	TenantId    string    `gorm:"column:tenant_id"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime;type:datetime"` // 创建时间
}

func (m *AlarmRuleResourceRel) TableName() string {
	return "t_alarm_rule_resource_rel"
}

type ResourceResourceGroupRel struct {
	Id              uint64    `gorm:"column:id;primary_key;autoIncrement"`
	ResourceGroupId string    `gorm:"column:resource_group_id;"`
	ResourceId      string    `gorm:"column:resource_id;"`
	TenantId        string    `gorm:"column:tenant_id"`
	CreateTime      time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime;type:datetime"` // 创建时间
}

func (m *ResourceResourceGroupRel) TableName() string {
	return "t_resource_resource_group_rel"
}

type AlarmHandler struct {
	Id           string    `gorm:"column:id;primary_key"`
	AlarmRuleId  string    `gorm:"column:alarm_rule_id;"`
	HandleType   int       `gorm:"column:handle_type"`   // 1邮件；2 短信；3弹性
	HandleParams string    `gorm:"column:handle_params"` //回调地址
	TenantId     string    `gorm:"column:tenant_id"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime;type:datetime"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime;type:datetime"` // 创建时间
}

func (m *AlarmHandler) TableName() string {
	return "t_alarm_handler"
}
