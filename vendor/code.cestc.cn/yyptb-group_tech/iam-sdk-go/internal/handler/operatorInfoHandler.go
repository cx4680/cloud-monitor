package handler

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/identitytypeenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/domain"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/errortypes"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler/roleswitchrecordhandler"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	// WebUserType 账户类型
	WebUserType = "userType"
	// WebAccountId 云账号（当前云账号本身账号ID或IAM归属的云账号
	WebAccountId = "accountId"
	WebUserId    = "userId"

	// UserIdentity 能开 用户数据存于此键值
	UserIdentity                 = "userIdentity"
	SsoCookieTicket              = "SID"
	CloudAccountOrganizeRoleName = "cloudAccountOrganizeRoleName"
	OrganizeAssumeRoleName       = "organizeAssumeRoleName"
)

// RetrieveUserInfo 获取用户信息
func RetrieveUserInfo(c *gin.Context) (*domain.OperatorInfo, error) {
	var (
		userTypeCode                 string
		accountId                    string
		cloudAccountId               string
		roleCrn                      string
		token                        string
		requestType                  = domain.Request_Http
		cid                          string
		cloudAccountOrganizeRoleName string
		organizeAssumeRoleName       string
	)

	userInfo := getUserInfoFromContext(c)
	if !isApiRequest(c) {
		userTypeCode = userInfo.UserTypeCode
		accountId = userInfo.AccountId
		cloudAccountId = userInfo.CloudAccountId
		cloudAccountOrganizeRoleName = userInfo.CloudAccountOrganizeRoleName
		organizeAssumeRoleName = userInfo.OrganizeAssumeRoleName

		// 角色
		if strconv.Itoa(identitytypeenum.IamRole) == userTypeCode {
			cecSID := getCookie(c.Request.Cookies(), SsoCookieTicket)
			if len(cecSID) == 0 {
				cecSID = c.GetHeader(SsoCookieTicket)
				logger.Logger().Infof("【IAM SDK】 retrieveUserInfo  isWebApi,HEADER: SSO_COOKIE_TICKET:%s", cecSID)
			} else {
				logger.Logger().Infof("【IAM SDK】 retrieveUserInfo  isWebApi,COOKIE: SSO_COOKIE_TICKET:%s", cecSID)
			}

			infoRecord, err := roleswitchrecordhandler.GetRoleRecord(cecSID)
			if err != nil {
				return nil, err
			}

			roleCrn = infoRecord.RoleCrn
			token = infoRecord.SecurityToken
			cid = cecSID
		}

		logger.Logger().Infof("【IAM SDK】 retrieveUserInfo  isWebApi\nuserTypeCode:%s, accountId:%s, cloudAccountId:%s, cid:%s, roleCrn:%s, token:%s",
			userTypeCode, accountId, cloudAccountId, cid, roleCrn, token)
	} else if len(c.GetHeader(UserIdentity)) > 0 {
		userIdentityStr := c.GetHeader(UserIdentity)
		userIdentity := &domain.UserIdentity{}
		err := json.Unmarshal([]byte(userIdentityStr), userIdentity)
		if err != nil {
			logger.Logger().Errorf("【IAM SDK】 retrieveUserInfo  isOpenApi,UserIdentity信息格式错误")
			return nil, errortypes.IAMSDKError(autherrorenum.UserIdentityFormatError)
		}
		logger.Logger().Infof("【IAM SDK】 retrieveUserInfo  isOpenApi,userIdentity:%s", userIdentityStr)

		roleCrn = c.GetHeader("roleCrn")
		token = c.GetHeader("token")
		userTypeCode = strconv.Itoa(identitytypeenum.IdentityCode(userIdentity.Type))
		accountId = userIdentity.PrincipalId
		cloudAccountId = userIdentity.AccountId
		cloudAccountOrganizeRoleName = userIdentity.CloudAccountOrganizeRoleName
		organizeAssumeRoleName = userIdentity.OrganizeAssumeRoleName
		logger.Logger().Infof("【IAM SDK】 retrieveUserInfo  isOpenApi\nuserTypeCode:%s\naccountId:%s\ncloudAccountId:%s\nroleCrn:%s\ntoken:%s",
			userTypeCode, accountId, cloudAccountId, roleCrn, token)
	}

	//如果能开不携带web调用传入的identity信息 或者门户session 或者非openapi调用， 则认为是内部服务流转， 不进行鉴权
	if len(userTypeCode) == 0 || len(accountId) == 0 {
		logger.Logger().Warnf("【IAM SDK】 retrieveUserInfo  userTypeCode or accountId is null -> is inner invoke")
		return nil, nil
	}

	if strconv.Itoa(identitytypeenum.IamUser) != userTypeCode && strconv.Itoa(identitytypeenum.IamRole) != userTypeCode {
		userTypeCode = strconv.Itoa(identitytypeenum.Account)
	}

	// 正常的userId转换
	userId, err := strconv.Atoi(accountId)
	if err != nil || userId <= 0 || len(cloudAccountId) == 0 {
		logger.Logger().Errorf("【IAM SDK】 retrieveUserInfo 用户id格式错误")
		return nil, errortypes.IAMSDKError(autherrorenum.IamUserIdFormatError)
	}

	return &domain.OperatorInfo{AccountId: accountId, CloudAccountId: cloudAccountId, RequestType: requestType, UserTypeCode: userTypeCode, RoleCrn: roleCrn, Token: token, Cid: cid, OrganizeAssumeRoleName: organizeAssumeRoleName, CloudAccountOrganizeRoleName: cloudAccountOrganizeRoleName}, nil

}

func getUserInfoFromContext(c *gin.Context) *UserInfo {
	userTypeCode, _ := c.Get(WebUserType)
	accountId, _ := c.Get(WebUserId)
	cloudAccountId, _ := c.Get(WebAccountId)
	organizeAssumeRoleName, _ := c.Get(OrganizeAssumeRoleName)
	cloudAccountOrganizeRoleName, _ := c.Get(CloudAccountOrganizeRoleName)

	userInfo := &UserInfo{
		UserTypeCode:                 util.Strval(userTypeCode),
		AccountId:                    util.Strval(accountId),
		CloudAccountId:               util.Strval(cloudAccountId),
		OrganizeAssumeRoleName:       util.Strval(organizeAssumeRoleName),
		CloudAccountOrganizeRoleName: util.Strval(cloudAccountOrganizeRoleName),
	}
	return userInfo
}

func isApiRequest(c *gin.Context) bool {
	header := c.GetHeader("eventType")
	return reflect.DeepEqual("ApiCall", header)
}

func getCookie(cookies []*http.Cookie, name string) string {
	for _, cookie := range cookies {
		if strings.ToUpper(cookie.Name) == name {
			return cookie.Value
		}
	}
	return ""
}

type UserInfo struct {
	CloudAccountId               string `json:"cloudLoginId"`
	AccountId                    string `json:"loginId"`
	UserTypeCode                 string `json:"userTypeCode"`
	OrganizeAssumeRoleName       string `json:"organizeAssumeRoleName"`
	CloudAccountOrganizeRoleName string `json:"cloudAccountOrganizeRoleName"`
}
