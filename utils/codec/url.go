package codec

import "net/url"

// UrlEncode URL编码
func UrlEncode(urls string) (string, error) {
	return url.QueryEscape(urls), nil
}

// UrlDecode URL解码
func UrlDecode(urls string) (string, error) {
	urlStr, err := url.QueryUnescape(urls)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}
