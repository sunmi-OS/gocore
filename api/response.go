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

// NewResponse 获取默认返回内容
func NewResponse() Response {
	return defaultResponse
}

// SetDefaultCode 设置默认返回code码
func SetDefaultCode(code int) {
	defaultResponse.Code = code
}

// SetDefaultData 设置默认返回data内容
func SetDefaultData(data interface{}) {
	defaultResponse.Data = data
}

// SetDefaultMsg 设置默认返回msg内容
func SetDefaultMsg(msg string) {
	defaultResponse.Msg = msg
}
