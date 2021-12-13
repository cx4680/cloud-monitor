package vo

type PageVO struct {
	Records interface{} `json:"records"`
	Total   int         `json:"total"`
	Size    int         `json:"size"`
	Current int         `json:"current"`
	Pages   int         `json:"pages"`
}
