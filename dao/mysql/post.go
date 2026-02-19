package mysql

import "wild_bluebell/models"

// CreatePost 插入帖子数据
func CreatePost(p *models.Post) (err error) {
	// sql语句
	sqlStr := `insert into post (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)`
	// 执行sql语句
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return

}

// GetPostByID 通过Id查询帖子数据
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	// sql 语句
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	// 查询数据库
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
