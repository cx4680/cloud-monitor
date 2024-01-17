package constant

const (
	MaxContactNum   = 1000 //每个租户限制可创建联系人数量
	MaxGroupNum     = 100  //每个租户限制创建联系组数量
	MaxContactGroup = 100  //单个联系人可加入联系组数量

	Phone = 1 //手机
	Email = 2 //邮箱

	Activated = 1 //已激活
	NotActive = 0 //未激活

	PhoneSize    = 11  //手机号长度限制
	MaxEmailSize = 100 //邮箱最大长度限制

	DefaultContact = "云账号告警联系人" //默认联系人/组
)
