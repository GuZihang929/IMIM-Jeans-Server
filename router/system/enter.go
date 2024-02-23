package system

type RouterGroup struct {
	PublicRouter
	AuthRouter
	FriendRouter
	GroupRouter
}
