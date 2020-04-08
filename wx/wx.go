package wx

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/sunmi-OS/gocore/httplib"
	"gopkg.in/redis.v5"
)

type (
	Wx struct {
		appId       string
		secret      string
		grantType   string
		accessToken string
		getRedis    func() *redis.Client
	}

	GetUnLimitQRCodeRequest struct {
		Scene     string `json:"scene"`
		Page      string `json:"page"`
		AutoColor bool   `json:"auto_color"`
		IsHyaline bool   `json:"is_hyaline"`
		Width     int64  `json:"width"`
	}
	SendRequest struct {
		Openid          string `json:"touser"`
		TemplateId      string `json:"template_id"`
		Page            string `json:"page"`
		FormId          string `json:"form_id"`
		Data            string `json:"data"`
		EmphasisKeyword string `json:"emphasis_keyword"`
	}
	CheckLoginResponse struct {
		OpenId     string `json:"openId"`
		SessionKey string `json:"session_key"`
	}
)

const (
	// access_token 地址
	AccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	//获取无限制小程序二维码
	CreateUqrcodeUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	//授权$code 访问地址
	CodeAccessUrl    = "https://api.weixin.qq.com/sns/jscode2session"
	TemplatedSendUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
)

// @desc 初始化
// @auth liuguoqiang 2020-02-25
// @param
// @return
func NewWx(appId, secret, grantType string, getRedis func() *redis.Client) *Wx {
	return &Wx{
		appId:     appId,
		secret:    secret,
		grantType: grantType,
		getRedis:  getRedis,
	}
}

// @desc 根据access_token值进行授权
// @auth liuguoqiang 2020-02-25
// @param $isFresh 是否刷新access_token
// @return
func (s *Wx) InitAuthToken(isFresh bool) (string, error) {
	//查询缓存
	tokenKey := "wechat:applet:token:" + s.appId
	accessToken := s.getRedis().Get(tokenKey).Val()
	if accessToken != "" && !isFresh {
		s.accessToken = accessToken
		return s.accessToken, nil
	}

	// 获取token
	req := httplib.Get(AccessTokenUrl + "?grant_type=client_credential&appid=" + s.appId + "&secret=" + s.secret)
	data := make(map[string]interface{})
	err := req.ToJSON(&data)
	if err != nil {
		return "", err
	}
	if accessToken, ok := data["access_token"]; !ok {
		return "", errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
	} else {
		err := s.getRedis().Set(tokenKey, accessToken, 7000*time.Second).Err()
		if err != nil {
			return "", err
		}
		s.accessToken = accessToken.(string)
		return s.accessToken, nil
	}
}

// @desc 获取二维码
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s *Wx) GetUnLimitQRCode(params *GetUnLimitQRCodeRequest, isFresh bool) ([]byte, error) {
	return s.Request(params, CreateUqrcodeUrl, isFresh)
}

// @desc 微信小程序模板消息推送
// @auth liuguoqiang 2020-02-25
// @param $openid 接收者（用户）的 openid
// @param $template_id 所需下发的模板消息的id
// @param $page 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
// @param $form_id 表单提交场景下，为 submit 事件带上的 formId；支付场景下，为本次支付的 prepay_id
// @param $data type:object 模板内容，不填则下发空模板。具体格式请参考示例。
// @param $emphasis_keyword 模板需要放大的关键词，不填则默认无放大
func (s *Wx) Send(params *SendRequest, isFresh bool) ([]byte, error) {
	return s.Request(params, TemplatedSendUrl, isFresh)
}

func (s *Wx) Request(params interface{}, url string, isFresh bool) ([]byte, error) {
	if s.accessToken == "" || isFresh {
		_, err := s.InitAuthToken(isFresh)
		if err != nil {
			return nil, err
		}
	}

	req := httplib.Post(url + "?access_token=" + s.accessToken)
	req, err := req.JSONBody(params)
	if err != nil {
		return nil, err
	}
	dataByte, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err == nil {
		if _, ok := data["errcode"]; ok {
			if !isFresh {
				dataByte, err = s.Request(params, url, true)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
			}
		}
	}
	return dataByte, nil
}

// @desc 根据微信code获取授权信息
// @auth liuguoqiang 2020-04-08
// @param
// @return
func (s *Wx) CheckLogin(code string) (*CheckLoginResponse, error) {
	params := make(map[string]interface{})
	params["appid"] = s.appId
	params["secret"] = s.secret
	params["js_code"] = code
	params["grant_type"] = s.grantType
	req := httplib.Post(CodeAccessUrl)
	req, err := req.JSONBody(params)
	if err != nil {
		return nil, err
	}
	dataByte, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(dataByte, &data)
	if err == nil {
		if _, ok := data["errcode"]; ok {
			return nil, errors.New(strconv.FormatFloat(data["errcode"].(float64), 'f', -1, 64) + ":" + data["errmsg"].(string))
		}
	}
	return &CheckLoginResponse{
		OpenId:     data["openid"].(string),
		SessionKey: data["session_key"].(string),
	}, nil
}
