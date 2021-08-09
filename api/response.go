package api

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

var defaultResponse Response

func init() {
	defaultResponse = Response{
		Code: 1,
		Data: nil,
		Msg:  "",
	}
}

func NewResponse() *Response {
	return &defaultResponse
}

func SetDefaultCode(code int) {
	defaultResponse.Code = code
}

func SetDefaultData(data interface{}) {
	defaultResponse.Data = data
}

func SetDefaultMsg(msg string) {
	defaultResponse.Msg = msg
}
