package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"github.com/pkg/errors"
)

type InstanceLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InstanceCommonVO struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	Labels []InstanceLabel `json:"labels"`
}

type InstancePageForm struct {
	TenantId     string   `form:"tenantId"`
	InstanceId   string   `form:"instanceId"`
	InstanceName string   `form:"instanceName"`
	StatusList   []string `form:"statusList"`
	Current      int      `form:"current,default=1"`
	PageSize     int      `form:"pageSize,default=10"`
	Product      string   `form:"product"`
	ExtraAttr    map[string]string
}

type InstanceStage interface {
	convertRealForm(InstancePageForm) interface{}
	doRequest(string, interface{}) (interface{}, error)
	convertResp(realResp interface{}) (int, []InstanceCommonVO)
}

type InstanceService interface {
	GetPage(InstancePageForm, InstanceStage) (*vo.PageVO, error)
}

type InstanceServiceImpl struct {
}

func (is *InstanceServiceImpl) GetPage(form InstancePageForm, stage InstanceStage) (*vo.PageVO, error) {
	var err error
	f := stage.convertRealForm(form)

	url, err := is.getRequestUrl(form.Product)
	if err != nil {
		return nil, err
	}

	resp, err := stage.doRequest(url, f)
	if err != nil {
		return nil, err
	}

	total, list := stage.convertResp(resp)
	return &vo.PageVO{
		Records: list,
		Total:   total,
		Size:    form.PageSize,
		Current: form.Current,
		Pages:   (total / form.PageSize) + 1,
	}, nil
}

func (is *InstanceServiceImpl) getRequestUrl(product string) (string, error) {
	p := dao.MonitorProduct.GetByDesc(global.DB, product)
	if p == nil {
		return "", errors.New("产品配置有误")
	}
	//TODO 返回url字段拼接
	return "", nil
}
