package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	openapi "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"path"
)

//忽略认证的路径列表
var ignoreList = []string{"/hawkeye/contact/activateContact?*", "/inner/alarmRecord/**", "/actuator/**",
	"/hawkeye/inner/configItem/*",
	"/hawkeye/inner/monitorItem/*",
	"/hawkeye/inner/rule/*",
	"/hawkeye/inner/configItem/*",
	"/hawkeye/inner/monitorItem/*",
	"/hawkeye/inner/notice/*",
	"/hawkeye/inner/monitorChart/*",
	"/hawkeye/inner/monitorResource/*",
	"/hawkeye/inner/reportForm/*",
	"/hawkeye/inner/regionSync/*",
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI

		for _, i := range ignoreList {
			match, err := path.Match(i, uri)
			if err != nil {
				if openapi.OpenApiRouter(c) {
					c.JSON(http.StatusUnauthorized, openapi.NewRespError(openapi.AuthorizedFailed, c))
				} else {
					c.JSON(http.StatusUnauthorized, global.NewError("权限认证失败"))
				}
				c.Abort()
			}
			if match {
				req := c.Request
				tenantId := req.Header.Get("tenantId")
				c.Set(global.TenantId, tenantId)
				c.Set(global.UserId, tenantId)
				c.Next()
				return
			}
		}
		if err := ParsingAndSetUserInfo(c); err != nil {
			if openapi.OpenApiRouter(c) {
				c.JSON(http.StatusUnauthorized, openapi.NewRespError(openapi.AuthorizedFailed, c))
			} else {
				c.JSON(http.StatusUnauthorized, global.NewError("权限认证失败"))
			}
			c.Abort()
		}
		c.Next()
	}

}

func ParsingAndSetUserInfo(c *gin.Context) error {
	if config.Cfg.Common.Env == "local" {
		c.Set(global.UserType, "1")
		c.Set(global.TenantId, "1")
		c.Set(global.UserId, "1")
		c.Set(global.UserName, "jim")
		return nil
	}

	req := c.Request
	userInfoEncode := req.Header.Get("user-info")
	var ssoUser SsoUser
	if len(userInfoEncode) != 0 {
		userInfoDecode, err := base64.StdEncoding.DecodeString(userInfoEncode)
		if err != nil {
			return err
		}
		if strutil.IsNotBlank(string(userInfoDecode)) {
			jsonutil.ToObject(string(userInfoDecode), &ssoUser)
			c.Set(global.UserType, ssoUser.UserTypeCode)
			c.Set(global.TenantId, ssoUser.CloudLoginId)
			c.Set(global.UserId, ssoUser.LoginId)
			c.Set(global.UserName, ssoUser.LoginCode)
			c.Set(global.CloudAccountOrganizeRoleName, ssoUser.RoleName)
			c.Set(global.OrganizeAssumeRoleName, ssoUser.RoleName)
			return nil
		}
	}
	userIdentity := req.Header.Get("userIdentity")
	if len(userIdentity) != 0 {
		jsonutil.ToObject(userIdentity, &ssoUser)
		c.Set(global.UserType, ssoUser.TypeCode)
		c.Set(global.TenantId, ssoUser.AccountId)
		c.Set(global.UserId, ssoUser.PrincipalId)
		c.Set(global.UserName, ssoUser.UserName)
		c.Set(global.CloudAccountOrganizeRoleName, ssoUser.RoleName)
		c.Set(global.OrganizeAssumeRoleName, ssoUser.RoleName)
		return nil
	}
	return errors.New("无用户信息")

}

type SsoUser struct {
	CustName         string        `json:"custName"`
	CustId           string        `json:"custId"`
	LoginId          string        `json:"loginId"`
	LoginCode        string        `json:"loginCode"`
	UserTypeCode     string        `json:"userTypeCode"`
	SerialNumber     string        `json:"serialNumber"`
	UserHeadPortrait string        `json:"userHeadPortrait"`
	AcctId           string        `json:"acctId"`
	CloudLoginCode   string        `json:"cloudLoginCode"`
	CloudLoginId     string        `json:"cloudLoginId"`
	IsIdentify       string        `json:"isIdentify"`
	LoginState       string        `json:"loginState"`
	RoleName         string        `json:"roleName"`
	TypeCode         string        `json:"typeCode"`
	AccountId        string        `json:"accountId"`
	PrincipalId      string        `json:"principalId"`
	UserName         string        `json:"userName"`
	Organization     *Organization `json:"organization"`
}

type Organization struct {
	RoleName        string `json:"roleName"`
	RoleDisplayName string `json:"roleDisplayName"`
}
