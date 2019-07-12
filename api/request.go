//	PhalGo-Request
//	请求解析,获取get,post,json参数,签名加密,链式操作,并且参数验证
//	喵了个咪 <wenzhenxi@vip.qq.com> 2016/5/11
//  依赖情况:
//          "github.com/astaxie/beego/validation" 基于beego的拦截器(已经集成)
//          "github.com/labstack/echo" 依赖于echo

package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/sunmi-OS/gocore/api/validation"
	"github.com/sunmi-OS/gocore/encryption/des"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/tidwall/gjson"
)

var (
	ErrNoParams = errors.New("No params")
	ErrMD5      = errors.New("MD5 Check Error")

	ErrCustom = func(str string) error {
		return errors.New(str)
	}

	RequestBody = "requestBody"
)

type Request struct {
	Context     echo.Context
	params      *param
	Jsonparam   *Jsonparam
	valid       validation.Validation
	Json        gjson.Result
	IsJsonParam bool
	jsonTag     bool
	Debug       bool
}

type Jsonparam struct {
	key string
	val gjson.Result
}

type param struct {
	key string
	val string
	min int
	max int
}

// 初始化request
func NewRequest(c echo.Context) *Request {

	R := new(Request)
	R.Context = c
	//增加debug参数的匹配
	if R.Param("__debug__").SetDefault("").GetString() == "" {
		R.Debug = false
	} else {
		R.Debug = true
	}
	return R
}

// 清理参数
func (this *Request) Clean() {

	this.params = new(param)
	this.Jsonparam = new(Jsonparam)
}

// 返回报错信息
func (this *Request) GetError() error {

	if this.valid.HasErrors() {
		for _, v := range this.valid.Errors {

			return ErrCustom(v.Message + v.Key)
		}
	}

	return nil
}

// 进行签名验证以及DES加密验证
func (this *Request) InitDES() error {
	// debug_key 跳过签名验证以及DES加密验证
	debugKey := viper.C.GetString("system.debugKey")
	if debugKey != "" && this.Param("debug_key").GetString() == debugKey {
		return this.initWithoutDES()
	}
	params := ""
	params = this.PostParam(viper.C.GetString("system.DESParam")).GetString()

	//如果是开启了 DES加密 需要验证是否加密,然后需要验证签名,和加密内容
	if viper.C.GetBool("system.OpenDES") == true {
		if params == "" {
			return ErrNoParams
		}
	}

	if params != "" {

		if viper.C.GetBool("system.OpenSign") {

			sign := this.PostParam("sign").GetString()
			timeStamp := this.PostParam("timeStamp").GetString()
			randomNum := this.PostParam("randomNum").GetString()
			isEncrypted := this.PostParam("isEncrypted").GetString()

			if sign == "" || timeStamp == "" || randomNum == "" {
				return ErrMD5
			}

			keymd5 := md5.New()
			keymd5.Write([]byte(viper.C.GetString("system.MD5key")))
			md5key := hex.EncodeToString(keymd5.Sum(nil))

			signmd5 := md5.New()
			signmd5.Write([]byte(params + isEncrypted + timeStamp + randomNum + md5key))
			sign2 := hex.EncodeToString(signmd5.Sum(nil))

			if sign != sign2 {
				return ErrMD5
			}

			//如果是加密的params那么进行解密操作
			if isEncrypted == "1" {
				// webapi 情况下需要对params decode
				var err error
				params, err = url.PathUnescape(params)
				if err != nil {
					return err
				}
				base64params, err := base64.StdEncoding.DecodeString(params)
				if err != nil {
					return err
				}

				origData, err := des.DesDecrypt(string(base64params), viper.C.GetString("system.DESkey"), viper.C.GetString("system.DESiv"))

				if err != nil {
					return err
				}
				params = string(origData)
			}
		}
		this.Context.Set(RequestBody, string(params))
		this.Json = gjson.Parse(params)
		this.IsJsonParam = true
	}
	return nil
}

// 跳过签名和加密
func (this *Request) initWithoutDES() error {
	params := ""
	params = this.PostParam(viper.C.GetString("system.DESParam")).GetString()
	if params != "" {
		this.Context.Set(RequestBody, string(params))
		this.Json = gjson.Parse(params)
		this.IsJsonParam = true
	}
	return nil
}

// 使用Json参数传入Json字符
func (this *Request) SetJson(json string) {

	this.Json = gjson.Parse(json)
}

// 初始化restful-raw参数
func (this *Request) InitRawJson() error {

	body, err := ioutil.ReadAll(this.Context.Request().Body)
	if err != nil {
		return err
	}

	this.Context.Set(RequestBody, string(body))
	this.Json = gjson.Parse(string(body))
	this.IsJsonParam = true
	return nil
}

//--------------------------------------------------------获取参数-------------------------------------

// 获取Get参数
func (this *Request) GetParam(key string) *Request {

	this.Clean()
	str := this.Context.QueryParam(key)
	this.params.val = str
	this.params.key = key
	this.jsonTag = false

	return this
}

// 获取post参数
func (this *Request) PostParam(key string) *Request {

	this.Clean()
	str := this.Context.FormValue(key)
	this.params.val = str
	this.params.key = key
	this.jsonTag = false

	return this
}

// 获取请求参数顺序get->post
func (this *Request) Param(key string) *Request {

	var str string
	this.Clean()

	if this.IsJsonParam {

		this.Jsonparam.val = this.Json.Get(key)
		this.Jsonparam.key = key
		this.jsonTag = true
	} else {
		str = this.Context.QueryParam(key)

		if str == "" {
			str = this.Context.FormValue(key)
		}
		this.params.val = str
		this.params.key = key
		this.jsonTag = false
	}

	return this
}

func (this *Request) SetDefault(val string) *Request {
	if this.jsonTag == true {
		defJson := fmt.Sprintf(`{"index":"%s"}`, val)
		this.Jsonparam.val = gjson.Parse(defJson).Get("index")
	} else {
		this.params.val = val
	}
	return this
}

//----------------------------------------------------过滤验证------------------------------------

// GET,POST或JSON参数是否必须
func (this *Request) Require(b bool) *Request {

	this.valid.Required(this.getParamVal(), this.getParamKey()).Message("缺少必要参数,参数名称:")
	return this
}

// 设置参数最大值
func (this *Request) Max(i int) *Request {

	this.params.max = i
	return this
}

//设置参数最小值
func (this *Request) Min(i int) *Request {

	this.params.min = i
	return this
}

//--------------------------------------------GET,POST获取参数------------------------------------

// 获取并且验证参数 string类型 适用于GET或POST参数
func (this *Request) GetString() string {

	var str string

	str = this.getParamVal()
	if this.params.min != 0 {
		this.valid.MinSize(str, this.params.min, this.getParamKey()).
			Message("参数异常!参数长度为%d不能小于%d,参数名称:", len([]rune(str)), this.params.min)
	}
	if this.params.max != 0 {
		this.valid.MaxSize(str, this.params.max, this.getParamKey()).
			Message("参数异常!参数长度为%d不能大于%d,参数名称:", len([]rune(str)), this.params.max)
	}

	return str
}

// 获取并且验证参数 int类型 适用于GET或POST参数
func (this *Request) GetInt() int {
	var (
		i   int
		err error
	)

	if this.getParamVal() == "" {
		i = 0
	} else {
		i, err = strconv.Atoi(this.getParamVal())
		if err != nil {
			this.valid.SetError(this.getParamKey(), "参数异常!参数不是int类型,参数名称:")
		}
	}

	if this.params.min != 0 {
		this.valid.Min(i, this.params.min, this.getParamKey()).
			Message("参数异常!参数值为%d不能小于%d,参数名称:", i, this.params.min)
	}
	if this.params.max != 0 {
		this.valid.Max(i, this.params.max, this.getParamKey()).
			Message("参数异常!参数值为%d不能大于%d,参数名称:", i, this.params.max)
	}

	return i
}

// 获取并且验证参数 float64类型 适用于GET或POST参数
func (this *Request) GetFloat() float64 {

	var (
		i   float64
		err error
	)

	if this.getParamVal() == "" {
		i = 0
	} else {
		i, err = strconv.ParseFloat(this.getParamVal(), 64)
		if err != nil {
			this.valid.SetError(this.getParamKey(), "此参数无法转换为float64类型,参数名称:")
		}
	}

	if this.params.min != 0 {
		this.valid.Min(int(i), this.params.min, this.getParamKey()).
			Message("参数异常!参数值为%f不能小于%d,参数名称:", i, this.params.min)
	}
	if this.params.max != 0 {
		this.valid.Max(int(i), this.params.max, this.getParamKey()).
			Message("参数异常!参数值为%f不能大于%d,参数名称:", i, this.params.max)
	}

	return i
}

// 邮政编码
func (this *Request) ZipCode() *Request {

	this.valid.ZipCode(this.getParamVal(), this.getParamKey()).Message("参数异常!邮政编码验证失败,参数名称:")
	return this
}

// 手机号或固定电话号
func (this *Request) Phone() *Request {

	this.valid.Phone(this.getParamVal(), this.getParamKey()).Message("参数异常!手机号或固定电话号验证失败,参数名称:")
	return this
}

// 固定电话号
func (this *Request) Tel() *Request {

	this.valid.Tel(this.getParamVal(), this.getParamKey()).Message("参数异常!固定电话号验证失败,参数名称:")
	return this
}

// 手机号
func (this *Request) Mobile() *Request {

	this.valid.Mobile(this.getParamVal(), this.getParamKey()).Message("参数异常!手机号验证失败,参数名称:")
	return this
}

// base64编码
func (this *Request) Base64() *Request {

	this.valid.Base64(this.getParamVal(), this.getParamKey()).Message("参数异常!base64编码验证失败,参数名称:")
	return this
}

// IP格式，目前只支持IPv4格式验证
func (this *Request) IP() *Request {

	this.valid.IP(this.getParamVal(), this.getParamKey()).Message("参数异常!IP格式验证失败,参数名称:")
	return this
}

// 邮箱格式
func (this *Request) Email() *Request {

	this.valid.Email(this.getParamVal(), this.getParamKey()).Message("参数异常!邮箱格式验证失败,参数名称:")
	return this
}

// 正则匹配,其他类型都将被转成字符串再匹配(fmt.Sprintf(“%v”, obj).Match)
func (this *Request) Match(match string) *Request {

	this.valid.Match(this.getParamVal(), regexp.MustCompile(match), this.getParamKey()).Message("参数异常!正则验证失败,参数名称:")
	return this
}

// 反正则匹配,其他类型都将被转成字符串再匹配(fmt.Sprintf(“%v”, obj).Match)
func (this *Request) NoMatch(match string) *Request {

	this.valid.NoMatch(this.getParamVal(), regexp.MustCompile(match), this.getParamKey()).Message("参数异常!邮箱格式验证失败,参数名称:")
	return this
}

// 数字
func (this *Request) Numeric() *Request {

	this.valid.Numeric(this.getParamVal(), this.getParamKey()).Message("参数异常!数字格式验证失败,参数名称:")
	return this
}

// alpha字符
func (this *Request) Alpha() *Request {

	this.valid.Alpha(this.getParamVal(), this.getParamKey()).Message("参数异常!alpha格式验证失败,参数名称:")
	return this
}

// alpha字符或数字
func (this *Request) AlphaNumeric() *Request {

	this.valid.AlphaNumeric(this.getParamVal(), this.getParamKey()).Message("参数异常!AlphaNumeric格式验证失败,参数名称:")
	return this
}

// alpha字符或数字或横杠-_
func (this *Request) AlphaDash() *Request {

	this.valid.AlphaDash(this.getParamVal(), this.getParamKey()).Message("参数异常!AlphaDash格式验证失败,参数名称:")
	return this
}

// 返回解析参数的Val
func (this *Request) getParamVal() string {

	if this.jsonTag {
		return this.Jsonparam.val.String()
	} else {
		return this.params.val
	}
}

// 反回解析参数的Key
func (this *Request) getParamKey() string {
	if this.jsonTag {
		return this.Jsonparam.key
	} else {
		return this.params.key
	}
}

// 获取并且验证参数 Json类型 适用于Json参数
func (this *Request) GetJson() gjson.Result {

	return this.Jsonparam.val
}

//// 捕获panic异样防止程序终止 并且记录到日志
//func (this *Request)ErrorLogRecover() {
//
//	if err := recover(); err != nil {
//		this.Context.Response().Write([]byte("系统错误!具体原因:" + TurnString(err)))
//		LogError(err, map[string]interface{}{
//			"URL.Path":this.Context.Request().URL.Path,
//			"QueryParams":this.Context.QueryParams(),
//		})
//	}
//}
