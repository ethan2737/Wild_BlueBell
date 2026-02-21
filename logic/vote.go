package logic

import (
	"go.uber.org/zap"
	"strconv"
	"wild_bluebell/dao/redis"
	"wild_bluebell/models"
)

/*
投票的集中情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票
direction=0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，有两种情况
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票

投票的限制：每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票
	1. 到期之后将Redis中保存的赞成票及反对票数存储到MySQL表中
	2. 到期之后删除 KeyPostVotedZSetPF
*/

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
