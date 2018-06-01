package utils

import "net/url"

// URL编码
func UrlEncode(urls string) (string, error) {
	//UrlEnCode编码
	return url.QueryEscape(urls), nil
}

// URL解码
func UrlDecode(urls string) (string, error) {
	//UrlEnCode解码
	urlStr, err := url.QueryUnescape(urls)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
