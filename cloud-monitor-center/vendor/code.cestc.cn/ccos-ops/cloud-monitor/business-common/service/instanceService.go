package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"github.com/pkg/errors"
)

type InstanceLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InstanceCommonVO struct {
	InstanceId   string          `json:"instanceId"`
	InstanceName string          `json:"instanceName"`
	Labels       []InstanceLabel `json:"labels"`
}

type InstancePageForm struct {
	TenantId     string            `form:"tenantId"`
	InstanceId   string            `form:"instanceId"`
	InstanceName string            `form:"instanceName"`
	StatusList   string            `form:"statusList"`
	Current      int               `form:"current,default=1"`
	PageSize     int               `form:"pageSize,default=10"`
	Product      string            `form:"product"`
	ExtraAttr    map[string]string `form:"extraAttr"`
}

type InstanceStage interface {
	ConvertRealForm(InstancePageForm) interface{}
	DoRequest(string, interface{}) (interface{}, error)
	ConvertResp(realResp interface{}) (int, []InstanceCommonVO)
}

type InstanceService interface {
	GetPage(InstancePageForm, InstanceStage) (*vo.PageVO, error)
}

type InstanceServiceImpl struct {
}

func (is *InstanceServiceImpl) GetPage(form InstancePageForm, stage InstanceStage) (*vo.PageVO, error) {
	var err error
	f := stage.ConvertRealForm(form)

	url, err := is.getRequestUrl(form.Product)
	if err != nil {
		return nil, err
	}
	logger.Logger().Infof(" request  %+v ,%s", form, url)
	resp, err := stage.DoRequest(url, f)
	if err != nil {
		return nil, err
	}
	logger.Logger().Infof(" resp:%+v", resp)
	total, list := stage.ConvertResp(resp)
	return &vo.PageVO{
		Records: list,
		Total:   total,
		Size:    form.PageSize,
		Current: form.Current,
		Pages:   (total / form.PageSize) + 1,
	}, nil
}

func (is *InstanceServiceImpl) getRequestUrl(product string) (string, error) {
	p := dao.MonitorProduct.GetByAbbreviation(global.DB, product)
	if p == nil {
		return "", errors.New("产品配置有误")
	}
	if tools.IsBlank(p.Host) || tools.IsBlank(p.PageUrl) {
		return "", errors.New("产品配置有误")
	}
	return p.Host + p.PageUrl, nil
}
