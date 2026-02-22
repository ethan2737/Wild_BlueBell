package models

import "time"

// Post 帖子信息
// 内存对齐：在定义结构体时，相同类型的字段要放在一起，这样会节省内存空间
type Post struct {
	ID          int64     `json:"id,string" db:"post_id" example:"1234567890123456789"`          // 帖子ID，使用string格式解决js数字精度问题
	AuthorID    int64     `json:"author_id" example:"123456789"`                                 // 作者ID
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required" example:"1"` // 社区ID
	Status      int32     `json:"status" example:"1"`                                            // 帖子状态
	Title       string    `json:"title" db:"title" binding:"required" example:"这是一个帖子标题"`        // 标题
	Content     string    `json:"content" db:"content" binding:"required" example:"这是帖子内容..."`   // 内容
	CreateTime  time.Time `json:"create_time" example:"2024-01-01T00:00:00Z"`                    // 创建时间
}

// ApiPostDetail 帖子详情响应结构体
type ApiPostDetail struct {
	AuthorName       string                   `json:"author_name" example:"zhangsan"` // 作者名称
	VoteNum          int64                    `json:"vote_num" example:"10"`          // 投票数
	*Post                                     // 嵌入帖子结构体
	*CommunityDetail `json:"CommunityDetail"` // 嵌入社区信息
}
