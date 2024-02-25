package im

import "fmt"

const (
	//-------------------------------系统------------------------------------//
	redisKeyOnlineNum       = "online_num"      // 在线人数
	redisKeyUserSessionNum  = "session:%d_num"  // 用户离线消息数目
	redisKeyUserSessionMess = "session:%d_mess" // 用户离线最新消息

	//-------------------------------user------------------------------------//
	redisKeyMain      = "gamer.{%d}:main"      // user主要信息
	redisKeyGroup     = "gamer.{%d}:group"     // user加入的群聊
	redisKeyNotifyNum = "gamer.{%d}:NotifyNum" // user不在线通知数目

	//-------------------------------group------------------------------------//
	redisKeyGroupAllUser    = "group.{%d}:all"     // group用户
	redisKeyGroupOnlineUser = "group.{%d}:online"  // group在线用户
	redisKeyGroupChannel    = "group.{%d}:channel" // group频道
	redisKeyGroupAdmin      = "Group.{%d}:Admin"   //group管理员

)

// 在线人数
func GetRedisKeyOnlineNum() string {
	return redisKeyOnlineNum
}

// 用户离线消息数目
func GetRedisKeyUserSessionNum(Id int64) string {
	return fmt.Sprintf(redisKeyUserSessionNum, Id)
}

// 用户离线最新消息
func GetRedisKeyUserSessionMess(Id int64) string {
	return fmt.Sprintf(redisKeyUserSessionMess, Id)
}

// user主要信息key值
func GetRedisKeyMain(userid int64) string {
	return fmt.Sprintf(redisKeyMain, userid)
}

// user加入的群聊
func GetRedisKeyGroup(userid int64) string {
	return fmt.Sprintf(redisKeyGroup, userid)
}

// user不在线通知数目
func GetRedisKeyNotifyNum(userid int64) string {
	return fmt.Sprintf(redisKeyNotifyNum, userid)
}

// group用户
func GetRedisKeyGroupAllUser(id int64) string {
	return fmt.Sprintf(redisKeyGroupAllUser, id)
}

// group在线用户
func GetRedisKeyGroupOnlineUser(id int64) string {
	return fmt.Sprintf(redisKeyGroupOnlineUser, id)
}

// group频道
func GetRedisKeyGroupChannel(id int64) string {
	return fmt.Sprintf(redisKeyGroupChannel, id)
}

// group管理员
func GetRedisKeyGroupAdmin(id int64) string {
	return fmt.Sprintf(redisKeyGroupAdmin, id)
}
