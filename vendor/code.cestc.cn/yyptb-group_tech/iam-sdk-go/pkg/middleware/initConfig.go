package middleware

import "code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"

func InitIamConfig(identityUrl, regionId, logDir string, v ...string) {
	config.InitConfig(identityUrl, regionId, logDir)
}
