package controller

import "IM-Server/controller/system"

type ApiGroup struct {
	SystemApiGroup system.SystemControllerGroup
	FriendGroup    system.SystemControllerGroup
	GroupGroup     system.SystemControllerGroup
}

var ApiGroupApp = new(ApiGroup)
