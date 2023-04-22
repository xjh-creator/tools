package xresp

var (
	Ok = 200
)

// 通用api响应
type CommonRes struct {
	Code int         `json:"code"`    //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg  string      `json:"message"` //消息
	Data interface{} `json:"data"`    //数据内容
	//Total int64      `json:"total"`
}

// 通用api响应
type ApiResp struct {
	c *CommonRes
	r *ghttp.Request
}

func Parse(r *ghttp.Request, pointer interface{}) {
	m := r.GetMap()
	//for i, _ := range m {
	//	if i[0] == '$' {
	//		delete(m, i)
	//	}
	//}
	err := gconv.Struct(m, pointer)
	//err := gconv.Structs(m, pointer)
	if err != nil {
		Error(r).SetMsg(err.Error()).WriteJsonExit()
	}
	cherr := gvalid.CheckStruct(pointer, nil)
	if cherr != nil {
		Error(r).SetMsg(cherr.Error()).WriteJsonExit()
	}
}

//返回一个成功的消息体
func Success(r *ghttp.Request) *ApiResp {
	msg := CommonRes{
		Code: Ok,
		//Msg:  "操作成功",
		Msg: "",
	}
	var a = ApiResp{
		c: &msg,
		r: r,
	}
	return &a
}

//返回一个错误的消息体
func Error(r *ghttp.Request) *ApiResp {
	msg := CommonRes{
		Code: 500,
		Msg:  "操作失败",
	}
	var a = ApiResp{
		c: &msg,
		r: r,
	}
	return &a
}

//返回一个拒绝访问的消息体
func Forbidden(r *ghttp.Request) *ApiResp {
	msg := CommonRes{
		Code: 403,
		Msg:  "无操作权限",
	}
	var a = ApiResp{
		c: &msg,
		r: r,
	}
	return &a
}

//设置消息体的内容
func (resp *ApiResp) SetMsg(msg string) *ApiResp {
	resp.c.Msg = msg
	return resp
}

//设置消息体的编码
func (resp *ApiResp) SetCode(code int) *ApiResp {
	resp.c.Code = code
	return resp
}

//设置消息体的数据
func (resp *ApiResp) SetData(data interface{}) *ApiResp {
	resp.c.Data = data
	return resp
}

//设置消息体的编码
func (resp *ApiResp) SetStatusCode(status int) *ApiResp {
	resp.r.Response.Status = status
	return resp
}
func (resp *ApiResp) SetPageData(count int64, data interface{}) *ApiResp {
	//resp.c.Total = count
	resp.c.Data = g.Map{
		"records": data,
		"total":   count,
	}
	return resp
}

//输出json到客户端
func (resp *ApiResp) WriteJsonExit() {
	if resp.r.Response.Status > 200 {
		resp.c.Code = resp.r.Response.Status
	}
	if err := resp.r.Response.WriteJsonExit(resp.c); err != nil {
		glog.Error(err)
	}
}

// HTML 输出 html 到客户端
func (resp *ApiResp) HTML(str string) {
	if resp.r.Response.Status > 200 {
		resp.c.Code = resp.r.Response.Status
	}
	resp.r.Header.Set("Content-Type", "text/html")
	resp.r.Response.Write(str)
}
