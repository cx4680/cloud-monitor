package autherrorenum

const (
	JsonFormatExceptionGetRoleRecordRequestParam = "JSON_FORMAT_EXCEPTION_GET_ROLE_RECORD_REQUEST_PARAM"
	RequestFailGetRoleRecord                     = "REQUEST_FAIL_GET_ROLE_RECORD"
	IoReadExceptionGetRoleRecord                 = "IO_READ_EXCEPTION_GET_ROLE_RECORD"
	JsonFormatExceptionGetRoleRecordResponse     = "JSON_FORMAT_EXCEPTION_GET_ROLE_RECORD_RESPONSE"
	ActionNotAllowedGetRoleRecord                = "ACTION_NOT_ALLOWED_GET_ROLE_RECORD"

	JsonFormatExceptionUserAuthRequestParam = "JSON_FORMAT_EXCEPTION_USER_AUTH_REQUEST_PARAM"
	RequestFailUserAuth                     = "REQUEST_FAIL_USER_AUTH"
	IoReadExceptionUserAuth                 = "IO_READ_EXCEPTION_USER_AUTH"
	JsonFormatExceptionUserAuthResponse     = "JSON_FORMAT_EXCEPTION_USER_AUTH_RESPONSE"

	JsonFormatExceptionRoleAuthRequestParam = "JSON_FORMAT_EXCEPTION_ROLE_AUTH_REQUEST_PARAM"
	RequestFailRoleAuth                     = "REQUEST_FAIL_ROLE_AUTH"
	IoReadExceptionRoleAuth                 = "IO_READ_EXCEPTION_ROLE_AUTH"
	JsonFormatExceptionRoleAuthResponse     = "JSON_FORMAT_EXCEPTION_ROLE_AUTH_RESPONSE"

	JsonFormatExceptionGetTokenRequestParam = "JSON_FORMAT_EXCEPTION_GET_TOKEN_REQUEST_PARAM"
	RequestFailGetToken                     = "REQUEST_FAIL_GET_TOKEN"
	IoReadExceptionGetToken                 = "IO_READ_EXCEPTION_GET_TOKEN"
	JsonFormatExceptionGetTokenResponse     = "JSON_FORMAT_EXCEPTION_GET_TOKEN_RESPONSE"
	ActionNotAllowedGetToken                = "ACTION_NOT_ALLOWED_GET_TOKEN"

	IamUserIdFormatError = "IamUserIdFormatError"

	JsonFormatException = "JSON_FORMAT_EXCEPTION"

	AuthResponseError = "AUTH_RESPONSE_ERROR"

	AuthError = "AUTH_ERROR"

	IamRoleStsTokenInvalid = "IAM_ROLE_STS_TOKEN_INVALID"

	StsInvokeError = "STS_INVOKE_ERROR"

	OperNotAllowed = "OPER_NOT_ALLOWED"

	UserIdentityFormatError = "UserIdentityFormatError"
)

var errorText = map[string]string{
	JsonFormatExceptionGetRoleRecordRequestParam: "获得角色切换记录请求参数JSON转换异常",
	RequestFailGetRoleRecord:                     "获得角色切换记录请求失败",
	IoReadExceptionGetRoleRecord:                 "获得角色切换记录响应IO读取异常",
	JsonFormatExceptionGetRoleRecordResponse:     "获得角色切换记录响应体解析失败",
	ActionNotAllowedGetRoleRecord:                "用户不允许进行该操作",

	JsonFormatExceptionUserAuthRequestParam: "IAM用户鉴权请求参数JSON转换异常",
	RequestFailUserAuth:                     "IAM用户鉴权请求失败",
	IoReadExceptionUserAuth:                 "IAM用户鉴权响应IO读取异常",
	JsonFormatExceptionUserAuthResponse:     "IAM用户鉴权响应体解析失败",

	JsonFormatExceptionRoleAuthRequestParam: "角色鉴权请求参数JSON转换异常",
	RequestFailRoleAuth:                     "角色鉴权请求失败",
	IoReadExceptionRoleAuth:                 "角色鉴权响应IO读取异常",
	JsonFormatExceptionRoleAuthResponse:     "角色鉴权响应体解析失败",

	JsonFormatExceptionGetTokenRequestParam: "获得令牌请求参数JSON转换异常",
	RequestFailGetToken:                     "获得令牌请求失败",
	IoReadExceptionGetToken:                 "获得令牌响应IO读取异常",
	JsonFormatExceptionGetTokenResponse:     "获得令牌响应体解析失败",
	ActionNotAllowedGetToken:                "用户不允许进行该操作",

	IamUserIdFormatError: "用户id格式错误",

	JsonFormatException: "json解析失败",

	AuthResponseError: "服务器端返回鉴权结果异常",

	AuthError: "鉴权服务异常",

	IamRoleStsTokenInvalid: "Token已失效",

	/**
	 * TODO从service拉过来的， 需要对service错误码做判断， sts也有引用， 统一后续提取到share里!!!
	 */
	StsInvokeError: "获取sts令牌失败",

	OperNotAllowed: "用户不允许进行该操作",

	UserIdentityFormatError: "UserIdentity信息格式错误",
}

func ErrorText(code string) string {
	return errorText[code]
}
