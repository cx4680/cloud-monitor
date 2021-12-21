package models

type ConfigItem struct {
	Id     string `gorm:"column:id" json:"id"`
	Pid    string `gorm:"column:pid" json:"pid"`        //上级Id
	Name   string `gorm:"column:name" json:"name"`      //配置名称
	Code   string `gorm:"column:code" json:"code"`      //配置编码
	Data   string `gorm:"column:data" json:"data"`      //配置值
	SortId int    `gorm:"column:sort_id" json:"sortId"` //排序
	Remark string `gorm:"column:remark" json:"remark"`  //备注
}

func (m *ConfigItem) TableName() string {
	return "config_item"
}
