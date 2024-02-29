package system

import (
	"IM-Server/controller"
	"github.com/gin-gonic/gin"
)

type GroupRouter struct {
}

func (*GroupRouter) InitGroupRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	groupRouter := Router.Group("group")

	//获取路由函数
	var GroupController = controller.ApiGroupApp.GroupGroup
	{
		groupRouter.POST("group_create", GroupController.CreateGroup)
		groupRouter.POST("group_join", GroupController.JoinGroup)
		groupRouter.GET("group_get", GroupController.GetGroup)
		groupRouter.PUT("group_update", GroupController.UpdateGroup)
		groupRouter.DELETE("group_delete", GroupController.DeleteGroup)
		groupRouter.DELETE("group_delete_user", GroupController.DeleteGroupUser)
		groupRouter.GET("/get_userbyid", GroupController.GetUserByUserID)
	}
	return groupRouter
}
