package system

import (
	"IM-Server/controller"
	"github.com/gin-gonic/gin"
)

type PublicRouter struct {
}

func (s *PublicRouter) InitPublicRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	sysRouter := Router.Group("public")

	//获取路由函数
	var publicController = controller.ApiGroupApp.SystemApiGroup
	{
		sysRouter.POST("/login", publicController.Login)
		sysRouter.POST("/register", publicController.RegisterUser)
		sysRouter.POST("/sendEmail", publicController.SendEmail)
		sysRouter.POST("/visitorLogin", publicController.VisitorLogin)
		sysRouter.GET("/visitorLogout", publicController.VisitorLogout)

	}
	return sysRouter
}
