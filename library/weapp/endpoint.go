package weapp

import "net/url"

var _ Endpoint = (*Endpoint1)(nil)

// Endpoint1 实现了 github.com/chanxuehong/wechat/oauth2.Endpoint 接口.
type Endpoint1 struct {
	AppId     string
	AppSecret string
}

func NewEndpoint(AppId, AppSecret string) *Endpoint1 {
	return &Endpoint1{
		AppId:     AppId,
		AppSecret: AppSecret,
	}
}

func (p *Endpoint1) ExchangeTokenURL(code string) string {
	return "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(p.AppId) +
		"&secret=" + url.QueryEscape(p.AppSecret) +
		"&code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
}

func (p *Endpoint1) RefreshTokenURL(refreshToken string) string {
	return "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + url.QueryEscape(p.AppId) +
		"&grant_type=refresh_token&refresh_token=" + url.QueryEscape(refreshToken)
}

func (p *Endpoint1) SessionCodeUrl(code string) string {

	return "https://api.weixin.qq.com/sns/jscode2session?appid=" + url.QueryEscape(p.AppId) +
		"&secret=" + url.QueryEscape(p.AppSecret) +
		"&js_code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
}
