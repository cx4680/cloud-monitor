package global

type Resp struct {
	ErrorCode  string        `json:"errorCode"`
	ErrorMsg   string        `json:"errorMsg"`
	Success    bool          `json:"success"`
	Module     interface{}   `json:"module"`
	AllowRetry bool          `json:"allowRetry"`
	ErrorArgs  []interface{} `json:"errorArgs"`
}

func NewResp(code, msg string, success bool, module interface{}) *Resp {
	return &Resp{
		ErrorCode: code,
		ErrorMsg:  msg,
		Success:   success,
		Module:    module,
	}
}

func NewError(msg string) *Resp {
	return NewResp(ErrorServer, msg, false, nil)
}

func NewSuccess(msg string, data interface{}) *Resp {
	return NewResp(SuccessServer, msg, true, data)
}
