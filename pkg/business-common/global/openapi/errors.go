package openapi

type ErrorCode struct {
	Code      string
	MessageEn string
	MessageCn string
}

var OperationDenied = &ErrorCode{Code: "FailedOperation.OperationDenied", MessageEn: "The request has failed due to the operation deny", MessageCn: "对不起，您没有操作权限"}
var AuthException = &ErrorCode{Code: "FailedOperation.AuthException", MessageEn: "The request has failed due to IAM authentication service is abnormal", MessageCn: "IAM鉴权服务异常"}
var UserInformationError = &ErrorCode{Code: "FailedOperation.UserInformationError", MessageEn: "The request has failed due to get user information error", MessageCn: "获取用户信息失败"}
var DataQueryFailed = &ErrorCode{Code: "FailedOperation.DataQueryFailed", MessageEn: "Data query failed", MessageCn: "数据查询失败"}
var ContactNumberExceeded = &ErrorCode{Code: "FailedOperation.ContactNumberExceeded", MessageEn: "The request has failed due to the contacts number exceeded 100", MessageCn: "联系人最多创建100个"}
var GroupNumberExceeded = &ErrorCode{Code: "FailedOperation.GroupNumberExceeded", MessageEn: "The request has failed due to the  contact group number exceeded 10", MessageCn: "联系组最多创建10个"}
var ContactGroupNumberExceeded = &ErrorCode{Code: "FailedOperation.ContactGroupNumberExceeded", MessageEn: "The request has failed due to the  contacts can join up to 5 contact groups", MessageCn: "联系人最多加入5个联系组"}
var GroupContactNumberExceeded = &ErrorCode{Code: "FailedOperation.GroupContactNumberExceeded", MessageEn: "The request has failed due to the  contact groups can join up to 100 contacts", MessageCn: "联系组最多加入100个联系人"}
var TenantHaveNotContact = &ErrorCode{Code: "FailedOperation.TenantNotHaveContact", MessageEn: "This tenant does not have this contact", MessageCn: "该租户无此联系人"}
var TenantHaveNotGroup = &ErrorCode{Code: "FailedOperation.TenantNotHaveGroup", MessageEn: "This tenant does not have this group", MessageCn: "该租户无此联系组"}
var InvalidActivationCode = &ErrorCode{Code: "FailedOperation.InvalidActivationCode", MessageEn: "This tenant does not have this group", MessageCn: "无效激活码"}
var AlarmNoticeChannelInvalid = &ErrorCode{Code: "InvalidParameter.AlarmNoticeChannel", MessageEn: "The specified parameter “HandlerType” is not valid", MessageCn: "告警渠道不存在"}
var RuleIdInvalid = &ErrorCode{Code: "InvalidParameter.RuleId", MessageEn: "The specified parameter “RuleId” is not valid", MessageCn: "告警规则不存在"}
var ComparisonOperatorInvalid = &ErrorCode{Code: "InvalidParameter.ComparisonOperator", MessageEn: "The specified parameter “ComparisonOperator” is not valid", MessageCn: "比较符不存在"}
var StatisticsInvalid = &ErrorCode{Code: "InvalidParameter.Statistics", MessageEn: "The specified parameter “Statistics” is not valid", MessageCn: "统计类型不存在"}
var InvalidParameter = &ErrorCode{Code: "InvalidParameter.InvalidParameter", MessageEn: "Invalid Parameter", MessageCn: "无效参数"}
var MetricCodeInvalid = &ErrorCode{Code: "InvalidParameter.MetricCode", MessageEn: "The specified parameter “MetricCode” is not valid", MessageCn: "指标不存在"}
var ProductAbbreviationInvalid = &ErrorCode{Code: "InvalidParameter.ProductAbbreviation", MessageEn: "The specified parameter “ProductAbbreviation” is not valid", MessageCn: "无效的产品简称"}
var SystemError = &ErrorCode{Code: "InternalError.System", MessageEn: "The request has failed due to inner system  error  of the server", MessageCn: "系统错误"}
var GroupNameRepeat = &ErrorCode{Code: "InvalidParameter.GroupNameRepeat", MessageEn: "The specified parameter “GroupName” is not valid", MessageCn: "联系组名重复"}
var MissingContact = &ErrorCode{Code: "MissingParameter.Contract", MessageEn: "The request least one contact", MessageCn: "至少选择一个联系人"}
var MissingContactId = &ErrorCode{Code: "MissingParameter.ContactId", MessageEn: "The required parameter 'ContactId' is not supplied", MessageCn: "联系人ID不能为空"}
var MissingContactName = &ErrorCode{Code: "MissingParameter.ContactName", MessageEn: "The required parameter 'ContactName' is not supplied", MessageCn: "联系人名字不能为空"}
var MissingGroupId = &ErrorCode{Code: "MissingParameter.GroupId", MessageEn: "The required parameter 'GroupId' is not supplied", MessageCn: "联系组ID不能为空"}
var MissingGroupName = &ErrorCode{Code: "MissingParameter.GroupName", MessageEn: "The required parameter 'GroupName' is not supplied", MessageCn: "联系组名字不能为空"}
var ContactNameFormatError = &ErrorCode{Code: "MissingParameter.ContactNameFormatError", MessageEn: "The contact name format error", MessageCn: "联系人名字格式错误"}
var GroupNameFormatError = &ErrorCode{Code: "MissingParameter.GroupNameFormatError", MessageEn: "The group name format error", MessageCn: "联系组名字格式错误"}
var MissingParameter = &ErrorCode{Code: "MissingParameter.MissingParameter", MessageEn: "Missing Parameter", MessageCn: "缺少参数"}
var PathNotFound = &ErrorCode{Code: "PathsNotFound", MessageEn: "The path not found ", MessageCn: "接口不存在"}
var AuthorizedFailed = &ErrorCode{Code: "AuthorizedFailed", MessageEn: "The request authorized failed", MessageCn: "权限认证失败"}
var PhoneFormatError = &ErrorCode{Code: "InvalidParameter.PhoneFormatError", MessageEn: "The phone format error", MessageCn: "手机号格式错误"}
var MailFormatError = &ErrorCode{Code: "InvalidParameter.MailFormatError", MessageEn: "The mail format error", MessageCn: "邮箱格式错误"}
var InformationMissing = &ErrorCode{Code: "InvalidParameter.InformationMissing", MessageEn: "Phone number and email address must be filled in", MessageCn: "手机号和邮箱必须填写一项"}
var MissingResource = &ErrorCode{Code: "MissingParameter.ResourceId", MessageEn: "The required parameter 'ResourceId' is not supplied", MessageCn: "实例ID不能为空"}
var ResourceError = &ErrorCode{Code: "FailedOperation.Resource", MessageEn: "The tenant does not have this resource", MessageCn: "该租户无此实例"}
var TimeParameterError = &ErrorCode{Code: "InvalidParameter.TimeParameter", MessageEn: "Time parameter error", MessageCn: "时间参数错误"}
var AuthorizedNoPermission = &ErrorCode{Code: "AuthorizedNoPermission", MessageEn: "The request authorized failed, no permission", MessageCn: "您没有操作权限"}
var RuleNameInvalid = &ErrorCode{Code: "InvalidParameter.RuleName", MessageEn: "The specified parameter “RuleName” is not valid", MessageCn: "规则名称无效"}
var ThresholdInvalid = &ErrorCode{Code: "InvalidParameter.Threshold", MessageEn: "The specified parameter “Threshold” is not valid", MessageCn: "阈值无效"}
var RuleLevelMissing = &ErrorCode{Code: "InvalidParameter.RuleLevel", MessageEn: "The specified parameter “Level” is not valid", MessageCn: "告警等级不能为空"}
var RuleSilencesTimeMissing = &ErrorCode{Code: "InvalidParameter.RuleSilencesTime", MessageEn: "The specified parameter “SilencesTime” is not valid", MessageCn: "告警间隔不能为空"}
var RulePeriodMissing = &ErrorCode{Code: "InvalidParameter.RulePeriod", MessageEn: "The specified parameter “Period” is not valid", MessageCn: "数据周期不能为空"}
var RuleTimesMissing = &ErrorCode{Code: "InvalidParameter.RuleTimes", MessageEn: "The specified parameter “Times” is not valid", MessageCn: "持续周期不能为空"}
var RuleCombinationMissing = &ErrorCode{Code: "InvalidParameter.RuleCombination", MessageEn: "The specified parameter “Combination” is not valid", MessageCn: "告警条件关系有误"}
var MetricCodeMissing = &ErrorCode{Code: "InvalidParameter.MetricCode", MessageEn: "The specified parameter “MetricCode” is not valid", MessageCn: "指标不能为空"}
var MissingResources = &ErrorCode{Code: "InvalidParameter.Resources", MessageEn: "The specified parameter “Resources” is not valid", MessageCn: "需要至少选择一个实例"}
