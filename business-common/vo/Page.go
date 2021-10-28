package vo

type PageVO struct {
	Records          interface{}   `json:"records"`
	Total            int           `json:"total"`
	Size             int           `json:"size"`
	Current          int           `json:"current"`
	Orders           []interface{} `json:"orders"`
	OptimizeCountSql bool          `json:"optimizeCountSql"`
	HitCount         bool          `json:"hitCount"`
	CountId          interface{}   `json:"countId"`
	MaxLimit         interface{}   `json:"maxLimit"`
	SearchCount      bool          `json:"searchCount"`
	Pages            int           `json:"pages"`
}
