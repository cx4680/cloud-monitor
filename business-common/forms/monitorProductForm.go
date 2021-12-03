package forms

type MonitorProductUpdateForm struct {
	Id         string `json:"id" binding:"required"`
	Name       string `json:"name"`
	CreateUser string `json:"createUser" binding:"required"`
}
