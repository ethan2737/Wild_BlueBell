package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"wild_bluebell/logic"
	"wild_bluebell/models"
)

// PostVoteController 投票处理函数
// @Summary 帖子投票
// @Description 对帖子进行投票，赞成票(1)、反对票(-1)或取消投票(0)
// @Tags 投票
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ParamVoteData true "投票参数"
// @Success 200 {object} ResponseData "投票成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1007 {object} ResponseData "需要登录"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	// 校验参数
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errMap := removeTopStruct(errs.Translate(trans)) // 翻译并去除掉错误提示中的结构体标识
		errMsg := ""
		for _, v := range errMap {
			errMsg = v
			break
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errMsg)
		return
	}
	// 获取当前用户的ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 处理投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}
