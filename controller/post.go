package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"wild_bluebell/logic"
	"wild_bluebell/models"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建新帖子，需要登录，传入社区ID、标题和内容
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Post true "帖子信息"
// @Success 200 {object} ResponseData "创建成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1007 {object} ResponseData "需要登录"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从 c 中获取当前发请求的用户Id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 2. 创建帖子（存入数据库）
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据帖子ID获取帖子详细信息
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} ResponseData{data=models.ApiPostDetail} "获取成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子Id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 根据Id取出帖子数据（查数据库）
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表
// @Summary 获取帖子列表（基础版）
// @Description 获取帖子列表，支持分页
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Success 200 {object} ResponseData{data=[]models.ApiPostDetail} "获取成功"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 获取帖子列表（升级版）
// @Description 根据排序参数动态获取帖子列表，支持按时间(time)或分数(score)排序
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Param order query string false "排序方式，time(时间)或score(分数)，默认time"
// @Success 200 {object} ResponseData{data=[]models.ApiPostDetail} "获取成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数（query string）：/api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
// @Summary 获取社区帖子列表
// @Description 根据社区ID获取该社区的帖子列表，支持分页和排序
// @Tags 帖子
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id query int true "社区ID"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Param order query string false "排序方式，time(时间)或score(分数)，默认time"
// @Success 200 {object} ResponseData{data=[]models.ApiPostDetail} "获取成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /community/posts [get]
func GetCommunityPostListHandler(c *gin.Context) {
	// 初始化结构体参数
	p := &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
	}
	// c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
