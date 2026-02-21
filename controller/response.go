package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 该模块的做用是将Gin框架中返回前端的内容做封装，使程序更加简洁
// 在前后端分离的项目中，通过定义标准的响应体内容来写作
/*
{
	"code": 10000,  // 程序中的错误码
	"msg": xxxx,    // 提示信息
	"data": {},     // 传递的数据
}
*/

// ResponseData 定义返回响应的结构体
type ResponseData struct {
	Code ResCode `json:"code"`
	// 提示信息和数据，因为不确定是什么类型，所以使用接口
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 返回错误
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 返回错误信息内容
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseSuccess 返回成功
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
