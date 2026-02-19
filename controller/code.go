package controller

// 定义程序中可能出现的错误码

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeInvalidToken
	CodeNeedLogin
)

// 错误码映射到错误信息
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "Success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeInvalidToken:    "无效的token",
	CodeNeedLogin:       "需要登录",
}

// Msg 根据错误码返回对应消息
func (rc ResCode) Msg() string {
	msg, ok := codeMsgMap[rc]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
