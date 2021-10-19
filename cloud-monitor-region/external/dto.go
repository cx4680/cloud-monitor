package external

type QueryPageRequest struct {
	/**
	 * 排序名称
	 */
	OrderName string `json:"orderName,omitempty"`

	/**
	 * 排序规则
	 */
	OrderRule string `json:"OrderRule,omitempty"`

	/**
	 * 当前页
	 */
	PageIndex int `json:"pageIndex,omitempty"`

	/**
	 * 每页显示的记录数
	 */
	PageSize int `json:"pageSize,omitempty"`

	/**
	 * 查询条件
	 */
	Data interface{} `json:"data,omitempty"`
}
