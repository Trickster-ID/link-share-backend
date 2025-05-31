package model

type BaseResponse struct {
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
	TotalData     int64       `json:"total_data,omitempty"`
	Url           string      `json:"url,omitempty"`
	ErrorLog      *ErrorLog   `json:"error,omitempty"`
}
