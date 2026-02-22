package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamsSignUp 用户注册请求参数
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required" example:"zhangsan"`                   // 用户名
	Password   string `json:"password" binding:"required" example:"123456"`                     // 密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password" example:"123456"` // 确认密码
}

// ParamsLogin 用户登录请求参数
type ParamsLogin struct {
	Username string `json:"username" binding:"required" example:"zhangsan"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"`   // 密码
}

// ParamVoteData 投票请求参数
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required" example:"1234567890123456789"` // 帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1" example:"1"`      // 投票方向：赞成票(1)、反对票(-1)、取消投票(0)
}

// ParamPostList 帖子列表查询参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page" example:"1"`      // 页码，默认1
	Size  int64  `json:"size" form:"size" example:"10"`     // 每页数量，默认10
	Order string `json:"order" form:"order" example:"time"` // 排序方式：time(时间)或score(分数)
}

// ParamCommunityPostList 社区帖子列表查询参数
type ParamCommunityPostList struct {
	*ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id" example:"1"` // 社区ID
}
