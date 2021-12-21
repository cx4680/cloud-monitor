package domain

type UserIdentity struct {
	AccountId string `json:"accountId"`

	/**
	 * 用户ID
	 */
	PrincipalId string `json:"principalId"`

	/**
	 * 用户类型
	 */
	Type string `json:"type"`
}
