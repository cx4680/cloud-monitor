package vo

type ConfigItemVO struct {
	Id     string `json:"id"`
	BizId  string `json:"bizId"`
	PBizId string `json:"pBizId"` //配置名称
	Name   string `json:"name"`   //配置编码
	Code   string `json:"code"`   //配置编码
	Data   string `json:"data"`   //配置值
	Remark string `json:"remark"` //备注
}
