package models

import "time"

// Post 内存对齐：在定义结构体时，相同类型的字段要放在一起，这样会节省内存空间
type Post struct {
	ID          int64     `json:"id,string" db:"post_id"` // json:"id,string 可以解决在前后端传递值时出现int类型失真的问题（超出js的数字类型范围）
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"` // 通过binding绑定参数，在gin框架上下文中传递
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string                   `json:"author_name"`
	*Post                                     // 嵌入帖子结构体
	*CommunityDetail `json:"CommunityDetail"` // 嵌入社区信息
}
