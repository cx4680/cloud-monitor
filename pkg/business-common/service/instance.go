package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
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
	Product      string            `form:"product" binding:"required"`
	ExtraAttr    map[string]string `form:"extraAttr"`
	IamInfo      IamInfo           `form:"iamInfo"`
}

type IamInfo struct {
	UserInfo                     string
	SID                          string
	CurrentTime                  string
	SecureTransport              string
	SourceIp                     string
	UserId                       string
	UserName                     string
	UserType                     string
	CloudAccountOrganizeRoleName string
	OrganizeAssumeRoleName       string
}

type InstanceStage interface {
	ConvertRealForm(InstancePageForm) interface{}
	DoRequest(string, interface{}) (interface{}, error)
	ConvertResp(realResp interface{}) (int, []InstanceCommonVO)

	ConvertRealAuthForm(InstancePageForm) interface{}
	DoAuthRequest(string, interface{}) (interface{}, error)
	ConvertAuthResp(realResp interface{}) (int, []InstanceCommonVO)
}

type InstanceService interface {
	GetPage(InstancePageForm, InstanceStage) (*vo.PageVO, error)
	GetPageByAuth(InstancePageForm, InstanceStage) (*vo.PageVO, error)
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

func (is *InstanceServiceImpl) GetPageByAuth(form InstancePageForm, stage InstanceStage) (*vo.PageVO, error) {
	var err error
	f := stage.ConvertRealAuthForm(form)

	url, err := is.getAuthRequestUrl(form.Product)
	if err != nil {
		return nil, err
	}
	logger.Logger().Infof(" request  %+v ,%s", form, url)
	resp, err := stage.DoAuthRequest(url, f)
	if err != nil {
		return nil, err
	}
	logger.Logger().Infof(" resp:%+v", resp)
	total, list := stage.ConvertAuthResp(resp)
	return &vo.PageVO{
		Records: list,
		Total:   total,
		Size:    form.PageSize,
		Current: form.Current,
		Pages:   (total / form.PageSize) + 1,
	}, nil
}

func (is *InstanceServiceImpl) getRequestUrl(product string) (string, error) {
	p := dao.MonitorProduct.GetByProductCode(global.DB, product)
	if p == nil {
		return "", errors.New("产品配置有误")
	}
	if strutil.IsBlank(p.Host) || strutil.IsBlank(p.PageUrl) {
		return "", errors.New("产品配置有误")
	}
	return p.Host + p.PageUrl, nil
}

func (is *InstanceServiceImpl) getAuthRequestUrl(product string) (string, error) {
	p := dao.MonitorProduct.GetByProductCode(global.DB, product)
	if p == nil {
		return "", errors.New("产品配置有误")
	}
	if strutil.IsBlank(p.Host) || strutil.IsBlank(p.IamPageUrl) {
		return "", errors.New("产品配置有误")
	}
	return p.Host + p.IamPageUrl, nil
}

func (is *InstanceServiceImpl) GetIamHeader(info *IamInfo) map[string]string {
	var headerParams = make(map[string]string)
	if info != nil {
		headerParams[global.UserInfo] = info.UserInfo
		headerParams[global.Cookie] = global.SID + "=" + info.SID
		headerParams[global.CsCurrentTime] = info.CurrentTime
		headerParams[global.CsSecureTransport] = info.SecureTransport
		headerParams[global.CsSourceIp] = info.SourceIp
		headerParams[global.CloudAccountOrganizeRoleName] = info.CloudAccountOrganizeRoleName
		headerParams[global.OrganizeAssumeRoleName] = info.OrganizeAssumeRoleName
	}
	return headerParams
}
