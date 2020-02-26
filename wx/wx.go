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
		redis       *redis.Client
	}

	GetUnLimitQRCodeRequest struct {
		Scene     string `json:"scene"`
		Page      string `json:"page"`
		AutoColor bool   `json:"auto_color"`
		IsHyaline bool   `json:"is_hyaline"`
		Width     int64  `json:"width"`
	}
)

const (
	// access_token 地址
	AccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token"
	//获取无限制小程序二维码
	CreateUqrcodeUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)

// @desc 初始化
// @auth liuguoqiang 2020-02-25
// @param
// @return
func NewWx(appId, secret, grantType string, redis *redis.Client) *Wx {
	return &Wx{
		appId:     appId,
		secret:    secret,
		grantType: grantType,
		redis:     redis,
	}
}

// @desc 根据access_token值进行授权
// @auth liuguoqiang 2020-02-25
// @param $isFresh 是否刷新access_token
// @return
func (s Wx) InitAuthToken(isFresh bool) (string, error) {
	//查询缓存
	tokenKey := "wechat:applet:token:" + s.appId
	accessToken := s.redis.Get(tokenKey).Val()
	if accessToken != "" && !isFresh {
		s.accessToken = accessToken
		return accessToken, nil
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
		err := s.redis.Set(tokenKey, accessToken, 7000*time.Second).Err()
		if err != nil {
			return "", err
		}
		s.accessToken = accessToken
		return accessToken.(string), nil
	}
}

// @desc 获取二维码
// @auth liuguoqiang 2020-02-25
// @param
// @return
func (s Wx) GetUnLimitQRCode(param *GetUnLimitQRCodeRequest, isFresh bool) ([]byte, error) {
	if s.accessToken == "" || isFresh {
		_, err := s.InitAuthToken(isFresh)
		if err != nil {
			return nil, err
		}
	}

	req := httplib.Post(CreateUqrcodeUrl + "?access_token=" + s.accessToken)
	req, err := req.JSONBody(param)
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
		if errcode, ok := data["errcode"]; ok {
			if errcode.(float64) == 40001 {
				dataByte, err = s.GetUnLimitQRCode(param, true)
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
