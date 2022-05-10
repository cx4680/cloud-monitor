package authhttp

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/domain"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/errortypes"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
)

// RoleTokenInvalidNotRetry 令牌失效已重试，不再重试
const RoleTokenInvalidNotRetry = false

type AuthResponse struct {
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode string `json:"errorCode"`
	Success   bool   `json:"success"`
	Module    *struct {
		Allowed          bool                      `json:"allowed"`
		IsHit            bool                      `json:"isHit"`
		IdentityResource *IdentityResourceResponse `json:"identityResource"`
		ListQueryRule    *ListQueryRule            `json:"listQueryRule"`
	} `json:"module"`
}

type IdentityResourceResponse struct {
	ResourceType     string   `json:"resourceType"`
	DenyResourceIDs  []string `json:"denyResourceIDs"`
	AllowResourceIDs []string `json:"allowResourceIDs"`
}

type ListQueryRule struct {
	AllowList []*ListQueryRuleCondition `json:"allowList"`
	DenyList  []*ListQueryRuleCondition `json:"denyList"`
}

type ListQueryRuleCondition struct {
	Region         *ListQueryRuleField `json:"region"`
	CloudAccountId *ListQueryRuleField `json:"cloudAccountId"`
	ResourceType   *ListQueryRuleField `json:"resourceType"`
	ResourceId     *ListQueryRuleField `json:"resourceId"`
}

type ListQueryRuleField struct {
	Field     string `json:"field"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

func (resp *AuthResponse) AuthResultRole(action *domain.ActionInfo, hasTry bool, operatorInfo *domain.OperatorInfo) (*AuthResponse, error) {
	if err := resp.authResultForNull(); err != nil {
		return nil, err
	}

	if len(resp.ErrorCode) > 0 && autherrorenum.IamRoleStsTokenInvalid == resp.ErrorCode {
		if hasTry {
			logger.Logger().Errorf("【IAM SDK】 AuthResultRole Token已失效")
			return nil, errortypes.IAMSDKError(autherrorenum.IamRoleStsTokenInvalid)
		}

		token, err := handler.GetToken(operatorInfo.RoleCrn, operatorInfo.Cid)
		if err != nil {
			return nil, err
		}

		request := &RoleAuthRequest{RoleCrn: operatorInfo.RoleCrn, Token: token}
		result, err := request.RoleAuth()
		if err != nil {
			return nil, err
		} else {
			result.AuthResultRole(action, RoleTokenInvalidNotRetry, operatorInfo)
		}
	}

	if err := resp.authResultForSuccess(); err != nil {
		return nil, err
	}

	if err := resp.authResultForModel(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (resp *AuthResponse) AuthResult() error {
	if err := resp.authResultForNull(); err != nil {
		return err
	}
	if err := resp.authResultForSuccess(); err != nil {
		return err
	}
	if err := resp.authResultForModel(); err != nil {
		return err
	}
	return nil
}

func (resp *AuthResponse) authResultForNull() error {
	if resp == nil {
		logger.Logger().Error("【IAM SDK】 authResultForNull,resp is nil")
		return errortypes.IAMSDKError(autherrorenum.AuthResponseError)
	}
	return nil
}

func (resp *AuthResponse) authResultForSuccess() error {
	if !resp.Success {
		logger.Logger().Errorf("【IAM SDK】 authResultForSuccess,resp is not success")
		return errortypes.IAMSDKError(autherrorenum.AuthError)

	}
	return nil
}

func (resp *AuthResponse) authResultForModel() error {
	if resp.Module == nil {
		logger.Logger().Errorf("【IAM SDK】 authResultForModel,resp.Module is nil")
		return errortypes.IAMSDKError(autherrorenum.AuthResponseError)
	}

	if !resp.Module.Allowed {
		logger.Logger().Errorf("【IAM SDK】 authResultForModel,resp.Module.Allowed is false")
		return errortypes.IAMSDKError(autherrorenum.OperNotAllowed)
	}

	return nil
}
