package http_request

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/sunmi-OS/gocore/v2/utils"
	"google.golang.org/grpc/metadata"
)

func TestSls(t *testing.T) {

}

func TestNew(t *testing.T) {
	c := New()
	c.Client.OnAfterResponse(MustCode200)
	c.SetLog(NewGocoreLog())
	c.Client.SetBaseURL("http://baidu.com")
	_, err := c.Client.R().Get("/")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
}

func TestNewWithParams(t *testing.T) {
	// the value of APP_NAME is after global variables in glog/logx/ctx.go, so cannot get APP_NAME with os.Setenv
	os.Setenv("APP_NAME", "http_test")

	c := New(WithRetryWaitTime(10), WithRetryMaxWaitTime(10), WithRetryCount(10))
	assert.Equal(t, 10, c.Client.RetryCount)
	assert.Equal(t, 10*time.Second, c.Client.RetryWaitTime)
	assert.Equal(t, 10*time.Second, c.Client.RetryMaxWaitTime)

	traceId := uuid.NewString()
	md := metadata.New(map[string]string{
		utils.XB3TraceId: traceId,
	})

	ctx := metadata.NewIncomingContext(context.Background(), md)
	req := c.R(ctx)
	assert.Equal(t, traceId, req.Header.Get(utils.XB3TraceId))
	//assert.Equal(t, "http_test", req.Header.Get(utils.XClientApp))
}
