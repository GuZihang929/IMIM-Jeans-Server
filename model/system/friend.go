package system

import "time"

type Friend struct {
	UserId     int64     `json:"user_id"`
	FriendId   int64     `json:"friend_id"`
	Status     int64     `json:"status"` // 0: 已删除 1: 已添加
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (Friend) TableName() string {
	return "im_user_friend"
}
