package models

// ParamsSignUp 定义请求的参数结构体
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required"` // binding 使用validator库进行参数校验，避免繁琐的逻辑判断参数是否有效
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // eqfield 判断一个字段是否等于另一个字段
}

// ParamsLogin 定义登录的参数结构体
type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
