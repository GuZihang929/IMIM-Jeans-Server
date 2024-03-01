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
		Router.POST("/del_session", authController.DelSession)
		Router.POST("/del_session_num", authController.DelSessionNum)
		Router.POST("/get_his_news", authController.GetHistoricalNew)
		Router.POST("/get_group_user", authController.GetGroupAndUser)

		//Router.GET("/world", authController.World)

	}
	return Router
}
