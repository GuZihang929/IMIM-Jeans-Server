package initialize

import (
	"IM-Server/middleware"
	"IM-Server/router"
	"github.com/gin-gonic/gin"
)

// 初始化总路由
// todo 初始化总路由
func Routers() *gin.Engine {
	Router := gin.Default()

	Router.Use()

	systemRouter := router.RouterGroupApp.System
	//    公开路有组,不做权限鉴定
	PublicGroup := Router.Group("api")
	PublicGroup.Use()
	{
		systemRouter.InitPublicRouter(PublicGroup)
	}
	//    私有路有组有拦截
	PrivateGroup := Router.Group("auth")
	PrivateGroup.Use(middleware.JwtAuth())
	{
		systemRouter.InitAuthRouter(PublicGroup)
	}

	friendRouter := router.RouterGroupApp.Friend
	//    私有路有组有拦截
	FriendGroup := Router.Group("apis")
	FriendGroup.Use(middleware.JwtAuth())
	{
		friendRouter.InitFriendRouter(FriendGroup)
	}

	groupRouter := router.RouterGroupApp.Group
	//    私有路有组有拦截
	GroupGroup := Router.Group("apis")
	GroupGroup.Use(middleware.JwtAuth())
	{
		groupRouter.InitGroupRouter(GroupGroup)
	}

	return Router
}
