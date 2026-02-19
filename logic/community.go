package logic

import (
	"wild_bluebell/dao/mysql"
	"wild_bluebell/models"
)

// GetCommunityList 获取社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 数据库查询 查找到所有的 community 并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
