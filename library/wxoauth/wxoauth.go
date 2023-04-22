package wxoauth

import (
	_ "embed"
)

var (
	//go:embed login.html
	wxlogin string //登录页面
	//go:embed confirm.html
	wxconfirm string //确认跳转页面
	//go:embed err.html
	wxerr string //确认跳转页面

)

func Login() string {
	return wxlogin
}

func Confirm() string {
	return wxconfirm
}
func Err() string {
	return wxerr
}
