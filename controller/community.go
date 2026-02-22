package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"wild_bluebell/logic"
)

// ----- 跟社区相关的内容 -------

// CommunityHandler 获取所有社区
// @Summary 获取社区列表
// @Description 获取所有社区列表，返回社区ID和名称
// @Tags 社区
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ResponseData{data=[]models.Community} "获取成功"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name）以列表（切片）的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 获取社区详情
// @Summary 获取社区详情
// @Description 根据社区ID获取社区详细信息
// @Tags 社区
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "社区ID"
// @Success 200 {object} ResponseData{data=models.CommunityDetail} "获取成功"
// @Failure 400 {object} ResponseData "请求参数错误"
// @Failure 1005 {object} ResponseData "服务器繁忙"
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区Id
	idStr := c.Param("id") // 获取URL参数，这里的key 要与路由一致
	// 将获取到的Id字符串转换成数字
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 通过ID查询社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}
