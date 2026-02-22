package logic

import (
	"go.uber.org/zap"
	"wild_bluebell/dao/mysql"
	"wild_bluebell/dao/redis"
	"wild_bluebell/models"
	"wild_bluebell/pkg/snowflake"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post ID
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
	// 3. 返回
}

// GetPostByID 根据帖子Id查询帖子详情数据
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID() failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	// 根据作者Id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	// 根据社区Id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者Id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区Id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 升级版获取帖子列表处理函数
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去Redis查询 ID 列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	// 如果 ids 为空，直接返回
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return
	}
	// 3. 根据 id 去 mysql 数据库查询帖子详细信息
	// 返回的数据还要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者Id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区Id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList 根据社区获取帖子列表
func GetCommunityPostList(p *models.ParamCommunityPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去Redis查询 ID 列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	// 如果 ids 为空，直接返回
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return
	}
	// 3. 根据 id 去 mysql 数据库查询帖子详细信息
	// 返回的数据还要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者Id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区Id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         votedata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
