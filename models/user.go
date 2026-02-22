package models

// User 用户信息
type User struct {
	UserID   int64  `json:"user_id" db:"user_id" example:"123456789"`                          // 用户ID
	Username string `json:"username" db:"username" example:"zhangsan"`                         // 用户名
	Password string `json:"password,omitempty" db:"password"`                                  // 密码
	Token    string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT token
}
