package system

// SystemControllerGroup
// Controller公共组件
type SystemControllerGroup struct {
	SysPublicController
	WebSocketController
	AuthController
	FriendController
	GroupController
}

//var userService = service.ServiceApp.SystemServiceGroup.UserService
//var captchaService = service.ServiceApp.SystemServiceGroup.CaptchaService
//var jwtService = service.ServiceApp.SystemServiceGroup.JwtService
