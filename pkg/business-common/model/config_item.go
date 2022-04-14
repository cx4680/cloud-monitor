package model

type ConfigItem struct {
	Id     uint64 `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	BizId  string `gorm:"column:biz_id;size=50" json:"bizId"`
	PBizId string `gorm:"column:p_biz_id;size=50" json:"pBizId"`
	Name   string `gorm:"column:name;size=100" json:"name"`     //配置名称
	Code   string `gorm:"column:code;size=100" json:"code"`     //配置编码
	Data   string `gorm:"column:data;size=100" json:"data"`     //配置值
	SortId uint32 `gorm:"column:sort_id" json:"sortId"`         //排序
	Remark string `gorm:"column:remark;size=200" json:"remark"` //备注
}

func (m *ConfigItem) TableName() string {
	return "t_config_item"
}
