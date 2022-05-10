package domain

type Response struct {
	ErrorMsg   string      `json:"errorMsg"`
	ErrorCode  string      `json:"errorCode"`
	Success    bool        `json:"success"`
	Module     interface{} `json:"module"`
	AllowRetry bool        `json:"allowRetry"`
	ErrorArgs  interface{} `json:"errorArgs"`
}
