package util

import "net/url"

// UrlResolve 获取绝对地址
func UrlResolve(s, path string) string {
	u, err := url.Parse(s)
	if err != nil {
		return path
	}
	u2, _ := u.Parse(path)
	return u2.String()
}
