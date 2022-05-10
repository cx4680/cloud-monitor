package domain

type ActionInfo struct {
	Product string

	/**
	 * 动作
	 */
	Action string

	/**
	 * 资源类型
	 */
	ResourceType string

	/**
	 * 资源列表
	 */
	Resources []Resource
}

type Resource struct {
	RelativeId string

	ResourceArn string
}
