package transaction

import (
	"dongguanquandao_server/library/global"
	"dongguanquandao_server/library/xresp"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

// 全局事务处理中间件
func Transaction(r *ghttp.Request) {
	method := r.Request.Method
	noTransaction := false
	if method == "OPTIONS" || method == "GET" || !global.Conf.System.Transaction {
		// OPTIONS/GET方法 以及 未配置事务时不创建事务
		noTransaction = true
	}
	defer func() {
		// 获取事务对象
		tx := global.GetTx(r)
		if err := recover(); err != nil {
			// 判断是否自定义响应结果
			if resp, ok := err.(xresp.CommonRes); ok {
				if !noTransaction {
					if resp.Code == xresp.Ok {
						// 有效的请求, 提交事务
						tx.Commit()
					} else {
						// 回滚事务
						tx.Rollback()
					}
				}
				// 以json方式写入响应
				xresp.Success(r).SetStatusCode(http.StatusOK).WriteJsonExit()
				return
			}
			if !noTransaction {
				// 回滚事务
				tx.Rollback()
			}
			// 继续向上层抛出异常
			panic(err)
		} else {
			if !noTransaction {
				// 没有异常, 提交事务
				tx.Commit()
			}
		}
		// 结束请求, 避免二次调用
		r.Exit()
	}()
	if !noTransaction {
		// 开启事务, 写入当前请求
		tx := global.Mysql.Begin()
		r.SetCtxVar("tx", tx)
	}
	// 处理请求
	r.Middleware.Next()
}
