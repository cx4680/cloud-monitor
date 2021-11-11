package forms

type EcsQueryPageForm struct {
	TenantId     string
	Current      int
	PageSize     int
	InstanceName string
	InstanceId   string
	Status       int
	StatusList   []int
}
