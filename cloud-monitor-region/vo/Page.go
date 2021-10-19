package vo

type PageVO struct {
	Records []interface{}
	Total   int
	Size    int
	Current int
}
