package models

import "time"

// Community 社区基本信息
type Community struct {
	ID   int64  `json:"id" db:"community_id" example:"1"`        // 社区ID
	Name string `json:"name" db:"community_name" example:"Go语言"` // 社区名称
}

// CommunityDetail 社区详细信息
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id" example:"1"`                              // 社区ID
	Name         string    `json:"name" db:"community_name" example:"Go语言"`                       // 社区名称
	Introduction string    `json:"introduction,omitempty" db:"introduction" example:"Go语言学习交流社区"` // 社区介绍
	CreateTime   time.Time `json:"create_time" db:"create_time" example:"2024-01-01T00:00:00Z"`   // 创建时间
}
