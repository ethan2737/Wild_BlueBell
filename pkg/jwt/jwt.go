package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// TokenExpireDuration token过期时间
// const TokenExpireDuration = time.Hour * 2

// 加盐（加密的字符串）
var mySecret = []byte("夏天夏天悄悄过去")

// MyClaims 自定义生命结构体并内嵌 jwt.StandardClaims
// jwt 包自带的 jwt.StandardClaims 只包含官方字段
// 我们这里需要额外记录一个userID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成jwt
func GenToken(userID int64, username string) (string, error) {
	// 创建一个自己声明的数据
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// 过期时间通过配置文件获取
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer:    "wild_bluebell",                                                                   // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的 secret 签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
