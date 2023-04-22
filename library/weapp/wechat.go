package weapp

import (
	"dongguanquandao_server/library/util"
	"fmt"

	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/miniprogram"
	//"github.com/silenceper/wechat/miniprogram"
)

var (
	//	//weapp wechat.
	Wxa *miniprogram.MiniProgram
)

func GetConfig(programName string) wechat.Config {
	var (
		err error
	)
	wechatConfig := (*wechat.Config)(nil)
	//weCfg := gcfg.Instance().Get("wechat")
	weCfg := gcfg.Instance().Get(programName) //配置表名称
	err = gconv.Struct(weCfg, &wechatConfig)
	if err != nil {
		panic(err)
	}
	wechatConfig.Cache = cache.NewMemory()
	return *wechatConfig
}

func Init() {

	config := GetConfig("wechat")
	wc := wechat.NewWechat(&config)
	Wxa = wc.GetMiniProgram()

	show_config := fmt.Sprintf("%s,%s", util.HideStar(config.AppID), util.HideStar(config.AppSecret))

	glog.Info("微信小程序设置:", show_config)
}
