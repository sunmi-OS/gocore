package web

type Config struct {
	// http export port. :8080
	Port string
}

type RecoverInfo struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   string      `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}

type CommonRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
