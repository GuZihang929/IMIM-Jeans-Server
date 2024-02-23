package router

import (
	"IM-Server/router/system"
)

type RouterGroup struct {
	System system.RouterGroup
	Friend system.RouterGroup
	Group  system.GroupRouter
}

var RouterGroupApp = new(RouterGroup)
