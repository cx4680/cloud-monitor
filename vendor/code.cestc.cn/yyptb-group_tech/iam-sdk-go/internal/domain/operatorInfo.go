package domain

type OperatorInfo struct {
	AccountId string

	/**
	 * 用户中心云账号ID
	 */
	CloudAccountId string

	/**
	 * 请求的方式
	 */
	RequestType RequestType

	/**
	 * IAM or ACCOUNT
	 */
	UserTypeCode string
	RoleCrn      string
	Token        string
	Cid          string
	/**
	 * 云账号之上组织角色
	 */
	OrganizeAssumeRoleName string
	/**
	 * 云账号之下组织角色
	 */
	CloudAccountOrganizeRoleName string
}

type RequestType int32

const (
	Request_Http RequestType = iota + 1
)
