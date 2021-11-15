module code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region

go 1.16

require (
	code.cestc.cn/ccos-ops/cloud-monitor/business-common v0.0.0-20211028091743-0cdf484dfcab
	code.cestc.cn/ccos-ops/cloud-monitor/common v0.0.0-20211028062752-e559c17fe0f2
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/gin-gonic/gin v1.7.4
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gohouse/converter v0.0.3
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/swaggo/gin-swagger v1.3.1
	github.com/swaggo/swag v1.7.3
	golang.org/x/net v0.0.0-20210928044308-7d9f5e0b762b
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.22.1
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
)

replace code.cestc.cn/ccos-ops/cloud-monitor/common => ../common

replace code.cestc.cn/ccos-ops/cloud-monitor/business-common => ../business-common
