package form

type MonitorProductUpdateForm struct {
	Id         uint64 `json:"id" binding:"required"`
	Name       string `json:"name"`
	CreateUser string `json:"createUser" binding:"required"`
}
