package jwt

import (
	"context"
	"dongguanquandao_server/library/lang"
	"dongguanquandao_server/library/token"
	"dongguanquandao_server/library/xdb"
	"dongguanquandao_server/library/xresp"
	"strings"
	"sync"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/net/ghttp"
)

var (
	checkLock sync.Mutex
)

// 鉴权中间件，只有登录成功之后才能通过
func JWT(r *ghttp.Request) {
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := r.Request.Header.Get("Authorization")
	// 按空格分割
	tokenStr := ""
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) > 1 {
			tokenStr = parts[1]
		}
	} else {
		tokenStr = r.Request.Header.Get("token")
	}
	if tokenStr == "" {
		tokenStr = r.GetQueryString("token")
	}
	if tokenStr == "" {
		//"获取token失败，请刷新页面"
		xresp.Error(r).SetMsg(lang.T(r.GetCtx()).T("NotFondToKen")).
			WriteJsonExit()
		return
	}
	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	t, err := token.VerifyAuthToken(tokenStr)
	if err != nil {
		xresp.Error(r).SetCode(401).SetMsg(lang.T(r.GetCtx()).T("NotFondToKen")).
			WriteJsonExit()
		return
	}
	// 将当前请求的uid信息保存到请求的上下文c上
	r.SetCtxVar("uid", t.Claim.StandardClaims.Id)
	// 将当前请求的部门信息保存到请求的上下文c上
	r.SetCtxVar("deptId", t.Claim.DepartmentId)
	//判断是否有新的token生成
	if t.NewToken != "" {
		r.Response.Header().Set("nt", t.NewToken)
	}
	r.Middleware.Next()
}
func check(sub, obj, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()
	// 检查策略
	pass, _ := xdb.MasterDB.Enforcer.Enforce(sub, obj, act)
	return pass
}

func JWT_WX(r *ghttp.Request) {
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := r.Request.Header.Get("Authorization")
	// 按空格分割
	tokenStr := ""
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) > 1 {
			tokenStr = parts[1]
		}
	} else {
		tokenStr = r.Request.Header.Get("token")
	}
	if tokenStr == "" {
		tokenStr = r.GetQueryString("token")
	}
	if tokenStr == "" {
		//"获取token失败，请刷新页面"
		//xresp.Error(r).SetMsg(lang.T(r.GetCtx()).T("NotFondToKen")).
		//	WriteJsonExit()
		//沒有token的 直接像是請登錄 401報錯
		xresp.Error(r).SetCode(401).SetMsg(lang.T(r.GetCtx()).T("NotFondToKen")).
			WriteJsonExit()
		return
	}
	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	t, err := token.VerifyAuthTokenWX(tokenStr)
	if err != nil {
		xresp.Error(r).SetCode(401).SetMsg(lang.T(r.GetCtx()).T("NotFondToKen")).
			WriteJsonExit()
		return
	}
	// 将当前请求的uid信息保存到请求的上下文c上
	r.SetCtxVar("uid", t.Claim.StandardClaims.Id)
	//判断是否有新的token生成
	if t.NewToken != "" {
		r.Response.Header().Set("nt", t.NewToken)
	}
	r.Middleware.Next()
}

func GetUserID(ctx context.Context) int64 {
	return gconv.Int64(ctx.Value("uid"))

}
