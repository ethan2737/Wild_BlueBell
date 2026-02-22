package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wild_bluebell/controller"
	"wild_bluebell/logger"
	"wild_bluebell/middlewares"

	_ "wild_bluebell/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 路由
func SetupRouter(mode string) *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册 swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 注册路由组
	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)
	// 应用jwt认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 帖子
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		// 根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)

		// 投票
		v1.POST("/vote", controller.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
