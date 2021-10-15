package global

type Resp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResp(code, msg string, data interface{}) *Resp {
	return &Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewError(msg string) *Resp {
	return NewResp(ErrorServer, msg, nil)
}

func NewSuccess(msg string, data interface{}) *Resp {
	return NewResp(SuccessServer, msg, data)
}
