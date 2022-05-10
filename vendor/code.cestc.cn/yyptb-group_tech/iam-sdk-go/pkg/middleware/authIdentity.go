package middleware

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/identitytypeenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/domain"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler/authhttp"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
)

// RoleTokenInvalidHasRetry 令牌失效允许重试
const (
	RoleTokenInvalidHasRetry = true
	STAR                     = "*"
	ContextListQueryRuleKey  = "ContextListQueryRuleKey"
)

func AuthIdentify(c *gin.Context, identity *models.Identity, crn string) error {
	logger.Logger().Infof("【IAM SDK】 AuthIdentify iam鉴权开始")
	logger.Logger().Infof("【IAM SDK】 AuthIdentify identity:%+v, crn:%s", *identity, crn)

	if c.Request == nil {
		logger.Logger().Warnf("【IAM SDK】 AuthIdentify Request is nil -> this is an inner api")
		printEndLog()
		return nil
	}

	// 白名单不鉴权
	if identity.IsWhite {
		logger.Logger().Infof("【IAM SDK】 AuthIdentify identity:%+v,命中白名单,不进行拦截", *identity)
		printEndLog()
		return nil
	}

	// 获取用户信息
	operatorInfo, err := handler.RetrieveUserInfo(c)
	if err != nil {
		return err
	} else if operatorInfo == nil {
		logger.Logger().Warnf("【IAM SDK】 AuthIdentify operatorInfo is nil -> this is an inner api")
		printEndLog()
		return nil
	}

	// 云账号不鉴权
	if isCloudAccountRequest(operatorInfo) {
		logger.Logger().Infof("【IAM SDK】 AuthIdentify 云账号不鉴权")
		printEndLog()
		return nil
	}

	// 组装鉴权信息
	actionInfo := operationInfoConstruct(identity, crn, operatorInfo)

	// 调用鉴权接口
	var er error
	var result *authhttp.AuthResponse

	if reflect.DeepEqual(operatorInfo.UserTypeCode, strconv.Itoa(identitytypeenum.IamRole)) {
		result, er = authRole(operatorInfo, actionInfo)
	} else {
		result, er = authUser(operatorInfo, actionInfo)
	}
	if er != nil {
		printEndLog()
		return er
	}

	logger.Logger().Infof("【IAM SDK】 AuthIdentify operatorInfo:%+v, actionInfo:%+v 请求鉴权成功,用户允许该操作", *operatorInfo, *actionInfo)

	if result.Module.ListQueryRule != nil {
		setContextListQueryRule(c, result.Module.ListQueryRule)
	}

	printEndLog()
	return nil
}

func printEndLog() {
	logger.Logger().Infof("【IAM SDK】 AuthIdentify iam鉴权鉴权结束\n")
}

func GetContextListQueryRule(c *gin.Context) interface{} {
	v, e := c.Get(ContextListQueryRuleKey)
	if e == false {
		logger.Logger().Infof("【IAM SDK】 GetContextListQueryRule ContextListQueryRuleKey is not exist")
		return nil
	}

	logger.Logger().Infof("【IAM SDK】 GetContextListQueryRule ContextListQueryRuleKey:%s", v)
	return v
}

func setContextListQueryRule(c *gin.Context, a *authhttp.ListQueryRule) {
	c.Set(ContextListQueryRuleKey, a)
	logger.Logger().Debugf("【IAM SDK】 setContextListQueryRule ContextListQueryRuleKey:%+v", *a)
}

func authUser(operatorInfo *domain.OperatorInfo, actionInfo *domain.ActionInfo) (*authhttp.AuthResponse, error) {
	var resources []string
	for _, resource := range actionInfo.Resources {
		resources = append(resources, resource.ResourceArn)
	}
	request := &authhttp.UserAuthRequest{
		AccountUserId: operatorInfo.AccountId,
		Product:       actionInfo.Product,
		ActionName:    actionInfo.Action,
		Resource:      resources}
	result, err := request.UserAuth()
	if err != nil {
		return nil, err
	}

	err = result.AuthResult()
	if err != nil {
		return nil, err
	}

	logger.Logger().Infof("【IAM SDK】 authUser IAM用户鉴权成功")
	return result, nil
}

func authRole(operatorInfo *domain.OperatorInfo, actionInfo *domain.ActionInfo) (*authhttp.AuthResponse, error) {
	var resources []string
	for _, resource := range actionInfo.Resources {
		resources = append(resources, resource.ResourceArn)
	}
	request := &authhttp.RoleAuthRequest{
		Product:    actionInfo.Product,
		ActionName: actionInfo.Action,
		RoleCrn:    operatorInfo.RoleCrn,
		Token:      operatorInfo.Token,
		Resource:   resources}
	result, err := request.RoleAuth()
	if err != nil {
		return nil, err
	}

	res, er := result.AuthResultRole(actionInfo, RoleTokenInvalidHasRetry, operatorInfo)
	if er != nil {
		return nil, er
	}

	logger.Logger().Infof("【IAM SDK】 authRole 角色鉴权成功")
	return res, nil
}

func isCloudAccountRequest(operatorInfo *domain.OperatorInfo) bool {
	return reflect.DeepEqual(operatorInfo.UserTypeCode, strconv.Itoa(identitytypeenum.Account))
}

func operationInfoConstruct(identity *models.Identity, crn string, operatorInfo *domain.OperatorInfo) *domain.ActionInfo {
	if len(identity.ResourceType) == 0 {
		identity.ResourceType = STAR
	}
	if len(identity.ResourceId) == 0 {
		identity.ResourceId = STAR
	}
	if len(identity.Region) == 0 && len(config.GetConfig().AuthSdkConfig.RegionId) > 0 {
		identity.Region = config.GetConfig().AuthSdkConfig.RegionId
	}

	actionInfo := &domain.ActionInfo{Product: identity.Product, Action: identity.Action, ResourceType: identity.ResourceType}

	if len(crn) > 0 {
		resource := domain.Resource{ResourceArn: crn}
		resources := append(actionInfo.Resources, resource)
		actionInfo.Resources = resources
		return actionInfo
	}

	resourceIdArray := strings.Split(identity.ResourceId, ",")
	var resources []domain.Resource
	for _, resourceId := range resourceIdArray {
		relativeId := identity.ResourceType + "/" + resourceId
		resourceCrn := strings.Join([]string{
			config.GetConfig().AuthSdkConfig.ResourceName,
			identity.Product, identity.Region,
			operatorInfo.CloudAccountId,
			relativeId}, ":")

		resource := domain.Resource{ResourceArn: resourceCrn, RelativeId: relativeId}
		resources = append(resources, resource)
	}
	actionInfo.Resources = resources

	return actionInfo
}
