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

// ParamVoteData 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              // 帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1）反对票（-1）取消投票（0）,oneof 必须是其中一个
}
