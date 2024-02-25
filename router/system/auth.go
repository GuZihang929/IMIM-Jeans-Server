package system

import (
	"IM-Server/controller"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
}

func (s *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	//获取路由函数
	var authController = controller.ApiGroupApp.SystemApiGroup
	{
		Router.POST("/info", authController.Info)

	}
	return Router
}
