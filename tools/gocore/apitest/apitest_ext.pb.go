// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v1.0.3
// - protoc            v4.24.2
// source: apitest/hello.proto

package apitest

import (
	sonic "github.com/bytedance/sonic"
	binding "github.com/gin-gonic/gin/binding"
	api "github.com/sunmi-OS/gocore/v2/api"
	ecode "github.com/sunmi-OS/gocore/v2/api/ecode"
	utils "github.com/sunmi-OS/gocore/v2/utils"
	math "math"
	http "net/http"
)

type TResponse[T any] struct {
	Code int32  `json:"code"`
	Data *T     `json:"data"`
	Msg  string `json:"msg"`
}

var defaultValidateErr error = api.ErrorBind
var releaseShowDetail bool
var disableValidate bool
var validateCode int = math.MaxInt

// set you error or use api.ErrorBind(diable:是否启用自动validate, 如果启用则返回 defaultValidateErr or 原始错误)
func SetAutoValidate(disable bool, validatErr error, releaseShowDetail bool) {
	disableValidate = disable
	defaultValidateErr = validatErr
	releaseShowDetail = releaseShowDetail
}

func SetValidateCode(code int) {
	validateCode = code
}

func checkValidate(err error) error {
	if disableValidate || err == nil {
		return nil
	}
	if utils.IsRelease() && !releaseShowDetail {
		return defaultValidateErr
	}

	if validateCode != math.MaxInt {
		return ecode.NewV2(validateCode, err.Error())
	}
	return err
}

const customReturnKey = "x-md-local-customreturn"

func SetCustomReturn(ctx *api.Context) {
	c := ctx.Request.Context()
	c = utils.SetMetaData(c, customReturnKey, "true")
	ctx.Request = ctx.Request.WithContext(c)
}

func setRetJSON(ctx *api.Context, resp interface{}, err error) {
	if utils.GetMetaData(ctx.Request.Context(), customReturnKey) != "" {
		return
	}
	ctx.RetJSON(resp, err)
}

func setRetOrigin(ctx *api.Context, resp interface{}) {
	if utils.GetMetaData(ctx.Request.Context(), customReturnKey) != "" {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func parseReq(ctx *api.Context, req interface{}) (err error) {
	if ctx.ContentType() == binding.MIMEPOSTForm {
		err = ctx.Request.ParseForm()
		if err != nil {
			return err
		}
		params := ctx.Request.FormValue("params")
		err = sonic.UnmarshalString(params, req)
		if err != nil {
			return err
		}
		err = binding.Validator.ValidateStruct(req)
	} else {
		err = ctx.ShouldBindJSON(req)
	}
	return
}
