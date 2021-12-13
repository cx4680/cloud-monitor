package models

type Identity struct {
	Product string

	/**
	 * 动作
	 */
	Action string
	/**
	操作资源的region
	*/
	Region string
	/**
	  白名单
	*/
	IsWhite bool

	/**
	 * 资源类型
	 */
	ResourceType string

	/**
	 * 资源ID
	 */
	ResourceId string
}
