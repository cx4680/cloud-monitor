package form

type NoticeChannel struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Data int8   `json:"data"`
}

type NoticeCenter struct {
	VerifyCodeIsOpen string `json:"verifyCodeIsOpen"`
	MsgIsOpen        string `json:"msgIsOpen"`
	MsgServiceIsOpen string `json:"msgServiceIsOpen"`
	MsgChannel       string `json:"msgChannel"`
}
