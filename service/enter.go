package service

import (
	"IM-Server/service/user"
)

type Service struct {
	SystemServiceGroup user.SysGroup
}

var ServiceApp = new(Service)
