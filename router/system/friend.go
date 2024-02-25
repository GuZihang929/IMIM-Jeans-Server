package system

import (
	"IM-Server/controller"
	"github.com/gin-gonic/gin"
)

type FriendRouter struct {
}

func (*FriendRouter) InitFriendRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	friendRouter := Router.Group("friend")

	//获取路由函数
	var FriendController = controller.ApiGroupApp.FriendGroup
	{
		friendRouter.POST("friend_create", FriendController.CreateFriend)
		friendRouter.DELETE("friend_delete", FriendController.DeleteFriend)
		friendRouter.GET("friend_list", FriendController.GetFriendList)

	}
	return friendRouter
}
