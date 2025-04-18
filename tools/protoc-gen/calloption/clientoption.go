package calloption

import (
	"github.com/go-resty/resty/v2"
)

type CallOption func(request *resty.Request)

func WithHeaders(headers map[string]string) CallOption {
	return func(o *resty.Request) {
		o.SetHeaders(headers)
	}
}
