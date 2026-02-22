package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"wild_bluebell/dao/mysql"
	"wild_bluebell/logic"
	"wild_bluebell/models"
)

// SignUpHandler 用户注册
// @Summary 用户注册接口
// @Description 用户注册，传入用户名、密码和确认密码进行注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.ParamsSignUp true "用户注册参数"
// @Success 200 {object} ResponseData "注册成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1002 {object} ResponseData "用户名已存在"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamsSignUp)
	// 绑定参数
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		// 记录日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errMsg := ""
		for _, v := range removeTopStruct(errs.Translate(trans)) {
			errMsg = v
			break
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errMsg)
		return
	}
	// 手动对请求参数进行详细的业务规则校验，（这部分逻辑被validator库代替）
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SingUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	fmt.Println(p)
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录
// @Summary 用户登录接口
// @Description 用户登录，传入用户名和密码，返回用户信息和 JWT token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.ParamsLogin true "用户登录参数"
// @Success 200 {object} ResponseData{data=map[string]string} "登录成功，返回 user_id, user_name, token"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1003 {object} ResponseData "用户名不存在"
// @Failure 1004 {object} ResponseData "用户名或密码错误"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	p := new(models.ParamsLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid params", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errMsg := ""
		for _, v := range removeTopStruct(errs.Translate(trans)) {
			errMsg = v
			break
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errMsg)
		return
	}
	// 2. 业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		// 记录日志
		zap.L().Error("login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
