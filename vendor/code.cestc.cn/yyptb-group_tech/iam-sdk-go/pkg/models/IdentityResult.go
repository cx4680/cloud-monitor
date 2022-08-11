package models

import "code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler/authhttp"

type IdentityResult struct {
	ListQueryRule         *authhttp.ListQueryRule `json:"listQueryRule"`
	ListResource          *authhttp.ListResource  `json:"listResource"`
	IsFullDisplayResource bool                    `json:"isFullDisplayResource"`
}
