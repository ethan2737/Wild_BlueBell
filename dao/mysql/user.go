package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"wild_bluebell/models"
)

// 把每一步数据库操作封装成函数
// 待 logic 层根据业务需要调用

const secret = "liwenzhou.com"

// CheckUserExist 查询指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	// 判断用户是否存在
	sqlStr := `select count(user_id) from user where username = ?`

	// 定义一个计数器，如果查询到的用户数量大于0，说明存在
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 1. 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 2. 执行sql语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 对密码加密
func encryptPassword(oPassword string) string {
	// 通过md5加盐实现加密
	h := md5.New()
	h.Write([]byte(secret))
	// 将加密后的字节转成16进制
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 用户登录
func Login(user *models.User) (err error) {
	// 记录用户登录的密码
	oPassword := user.Password
	// 查询用户信息
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	// 判断用户是否存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	// 查询数据库失败
	if err != nil {
		return err
	}
	// 将用户登录的密码加密后，与数据库中的密码比较，判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserByID 根据Id获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
