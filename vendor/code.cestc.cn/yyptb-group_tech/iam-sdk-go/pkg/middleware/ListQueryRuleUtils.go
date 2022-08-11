package middleware

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/handler/authhttp"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"strings"
)

func GenerateSqlQuery(rule *authhttp.ListQueryRule) string {
	// 查询规则不能为空
	if rule == nil {
		return ""
	}

	// 查询规则集合不能都为空
	if (rule.AllowList == nil || len(rule.AllowList) == 0) && (rule.DenyList == nil || len(rule.DenyList) == 0) {
		return ""
	}

	// 定义返回结果
	var result strings.Builder

	// 处理AllowList
	if rule.AllowList != nil && len(rule.AllowList) > 0 {
		result.WriteString(generateSQLCondition(rule.AllowList, true))
	}

	// 处理DenyList
	if rule.DenyList != nil && len(rule.DenyList) > 0 {
		if result.Len() > 0 {
			result.WriteString(" and ")
		}
		result.WriteString(generateSQLCondition(rule.DenyList, false))
	}

	logger.Logger().Infof("【IAM SDK】 generateSqlQuery sql:%s", result.String())
	return result.String()
}

func generateSQLCondition(conditions []*authhttp.ListQueryRuleCondition, allow bool) string {
	if conditions == nil || len(conditions) == 0 {
		return ""
	}

	var join = " and "
	if allow {
		join = " or "
	}

	// 定义返回结果
	var result strings.Builder
	result.WriteString("(")

	for i := 0; i < len(conditions); i++ {
		result.WriteString("(")

		var condition = conditions[i]

		result.WriteString(joinSqlColumn(condition.Region))
		if condition.Region != nil {
			result.WriteString(" and ")
		}
		result.WriteString(joinSqlColumn(condition.CloudAccountId))
		if condition.CloudAccountId != nil {
			result.WriteString(" and ")
		}
		result.WriteString(joinSqlColumn(condition.ResourceType))
		if condition.ResourceType != nil {
			result.WriteString(" and ")
		}
		result.WriteString(joinSqlColumn(condition.ResourceId))

		result.WriteString(")")
		if i < len(conditions)-1 {
			result.WriteString(join)
		}
	}

	if result.Len() == 1 {
		return ""
	} else {
		result.WriteString(")")
		return result.String()
	}
}

func joinSqlColumn(query *authhttp.ListQueryRuleField) string {
	if query == nil {
		return ""
	}

	// 定义返回结果
	var result strings.Builder
	result.WriteString(query.Field)
	result.WriteString(" ")
	result.WriteString(query.Condition)
	result.WriteString(" '")
	result.WriteString(strings.Replace(query.Value, "*", "%", -1))
	result.WriteString("'")
	return result.String()
}
