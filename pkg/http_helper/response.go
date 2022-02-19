package http_helper

type Response struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

type ResponseData struct {
	Response
	Data interface{} `json:"data"`
}
