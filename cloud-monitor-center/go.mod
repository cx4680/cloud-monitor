module code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center

go 1.15

require (
	code.cestc.cn/ccos-ops/cloud-monitor/business-common v0.0.0-00010101000000-000000000000
	code.cestc.cn/ccos-ops/cloud-monitor/common v0.0.0-20211028062752-e559c17fe0f2
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gohouse/converter v0.0.3
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tidwall/gjson v1.9.4 // indirect
	golang.org/x/sys v0.0.0-20211020064051-0ec99a608a1b // indirect
	gorm.io/gorm v1.22.1
)

replace code.cestc.cn/ccos-ops/cloud-monitor/common => ../common

replace code.cestc.cn/ccos-ops/cloud-monitor/business-common => ../business-common
