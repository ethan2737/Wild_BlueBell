// Package logic 存放与用户相关的业务逻辑代码
package logic

import (
	"wild_bluebell/dao/mysql"
	"wild_bluebell/models"
	"wild_bluebell/pkg/jwt"
	"wild_bluebell/pkg/snowflake"
)

// SignUp 用户注册
func SignUp(p *models.ParamsSignUp) (err error) {
	// 1. 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成ID
	userID := snowflake.GenID()

	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(user)
}

// Login 用户登录
func Login(p *models.ParamsLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成 jwt
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
