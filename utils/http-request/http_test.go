package http_request

import (
	"testing"
)

func TestSls(t *testing.T) {

}

func TestGlog(t *testing.T) {
	c := New()
	c.Client.OnAfterResponse(MustCode200)
	c.SetLog(NewGocoreLog())
	c.Client.SetHostURL("http://baidu.com")
	_, err := c.Client.R().Get("/")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
}
