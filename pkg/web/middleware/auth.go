package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	global2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	openapi2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
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
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI

		for _, i := range ignoreList {
			match, err := path.Match(i, uri)
			if err != nil {
				if openapi2.OpenApiRouter(c) {
					c.JSON(http.StatusUnauthorized, openapi2.NewRespError(openapi2.AuthorizedFailed, c))
				} else {
					c.JSON(http.StatusUnauthorized, global2.NewError("权限认证失败"))
				}
				c.Abort()
			}
			if match {
				req := c.Request
				tenantId := req.Header.Get("tenantId")
				c.Set(global2.TenantId, tenantId)
				c.Set(global2.UserId, tenantId)
				c.Next()
				return
			}
		}
		if err := ParsingAndSetUserInfo(c); err != nil {
			if openapi2.OpenApiRouter(c) {
				c.JSON(http.StatusUnauthorized, openapi2.NewRespError(openapi2.AuthorizedFailed, c))
			} else {
				c.JSON(http.StatusUnauthorized, global2.NewError("权限认证失败"))
			}
			c.Abort()
		}
		c.Next()
	}

}

func ParsingAndSetUserInfo(c *gin.Context) error {
	if config.Cfg.Common.Env == "local" {
		c.Set(global2.UserType, "1")
		c.Set(global2.TenantId, "1")
		c.Set(global2.UserId, "1")
		c.Set(global2.UserName, "jim")
		return nil
	}

	req := c.Request
	userInfoEncode := req.Header.Get("user-info")
	var userMap map[string]string
	if len(userInfoEncode) != 0 {
		userInfoDecode, err := base64.StdEncoding.DecodeString(userInfoEncode)
		if err != nil {
			return err
		}
		if strutil.IsNotBlank(string(userInfoDecode)) {
			jsonutil.ToObject(string(userInfoDecode), &userMap)
			c.Set(global2.UserType, userMap["userTypeCode"])
			c.Set(global2.TenantId, userMap["cloudLoginId"])
			c.Set(global2.UserId, userMap["loginId"])
			c.Set(global2.UserName, userMap["loginCode"])
			return nil
		}
	}
	userIdentity := req.Header.Get("userIdentity")
	if len(userIdentity) != 0 {
		jsonutil.ToObject(userIdentity, &userMap)
		c.Set(global2.UserType, userMap["typeCode"])
		c.Set(global2.TenantId, userMap["accountId"])
		c.Set(global2.UserId, userMap["principalId"])
		return nil
	}
	return errors.New("无用户信息")

}
