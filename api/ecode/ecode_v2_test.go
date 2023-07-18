package ecode

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sunmi-OS/gocore/v2/glog"
)

func TestEcodeWithReason(t *testing.T) {
	e := FromError(Success)
	glog.Debug(e.Error())   // error: code = 1 message = success metadata = map[] cause = <nil>
	glog.Debug(e.Code())    // 1
	glog.Debug(e.Message()) // success
	glog.Info("============================")

	e2 := FromError(nil)
	glog.Debug(e2.Error())   // error: code = 1 message = success metadata = map[] cause = <nil>
	glog.Debug(e2.Code())    // 1
	glog.Debug(e2.Message()) // success
	glog.Info("============================")

	sms := NewV2(10000, "中国电信").WithMetadata(map[string]string{
		"name":   "jerry",
		"reason": "欠话费了",
	})
	glog.Debug(sms.Error())   // error: code = 10000 message = 中国电信 metadata = map[name:jerry reason:我是metadata] cause = <nil>
	glog.Debug(sms.Code())    // 10000
	glog.Debug(sms.Message()) // 中国电信
	glog.Debug(sms.Metadata)  // map[name:jerry reason:欠话费了]
	glog.Info("============================")

	mms := NewV2(10086, "中国移动").WithCause(errors.New("我是原因"))
	glog.Debug(mms.Error())   // error: code = 10086 message = 中国移动 metadata = map[] cause = 我是原因
	glog.Debug(mms.Code())    // 10086
	glog.Debug(mms.Message()) // 中国电信
	glog.Debug(mms.Unwrap())  // 我是原因
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorV2
		err  error
		want bool
	}{
		{
			name: "true",
			e:    NewV2(404, ""),
			err:  NewV2(http.StatusNotFound, ""),
			want: true,
		},
		{
			name: "false",
			e:    NewV2(0, ""),
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.e.Is(tt.err); ok != tt.want {
				t.Errorf("ErrorV2.ErrorV2() = %v, want %v", ok, tt.want)
			}
		})
	}
}
