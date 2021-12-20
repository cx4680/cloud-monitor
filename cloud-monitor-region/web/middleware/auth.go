package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"path"
)

//忽略认证的路径列表
var ignoreList = []string{"/hawkeye/alertContact/certifyAlertContact?*", "/inner/alertRecord/**", "/actuator/**", "/swagger/**"}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI

		for _, i := range ignoreList {
			match, err := path.Match(i, uri)
			if err != nil {
				c.JSON(http.StatusUnauthorized, global.NewError("权限认证失败"))
				c.Abort()
			}
			if match {
				c.Next()
				return
			}
		}
		if err := ParsingAndSetUserInfo(c); err != nil {
			c.JSON(http.StatusUnauthorized, global.NewError("权限认证失败"))
			c.Abort()
		}
		c.Next()
	}

}

func ParsingAndSetUserInfo(c *gin.Context) error {
	req := c.Request
	userInfoEncode := req.Header.Get("user-info")
	userInfoDecode, err := base64.StdEncoding.DecodeString(userInfoEncode)
	if err != nil {
		return err
	}
	if tools.IsBlank(string(userInfoDecode)) {
		return errors.New("解析用户信息出错")
	}
	var userMap map[string]string
	tools.ToObject(string(userInfoDecode), &userMap)
	c.Set(global.UserType, userMap["userTypeCode"])
	c.Set(global.TenantId, userMap["cloudLoginId"])
	c.Set(global.UserId, userMap["loginId"])
	return nil
}
