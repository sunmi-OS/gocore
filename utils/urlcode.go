package utils

import "net/url"

// URL编码
func UrlEncode(urls string) (string, error) {
	//UrlEnCode编码
	urlStr, err := url.Parse(urls)
	if err != nil {
		return "", err
	}

	return urlStr.RequestURI(), nil
}

// URL解码
func UrlDecode(urls string) (string, error) {
	//UrlEnCode解码
	urlStr, err := url.Parse(urls)
	if err != nil {
		return "", err
	}

	return urlStr.Path, nil
}
