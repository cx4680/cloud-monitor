package form

type IamUser struct {
	Module struct {
		DirectoryId string `json:"departmentId"`
	} `json:"module"`
}

type IamLoginStart struct {
	Module bool `json:"module"`
}
