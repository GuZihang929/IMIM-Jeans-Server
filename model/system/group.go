package system

import "time"

type Group struct {
	UserId     int64     `json:"user_id"`
	GroupId    int64     `json:"group_id"`
	JoinInTime time.Time `json:"join_in_time"`
	Identity   int64     `json:"identity"` // 0:普通群员 1:群主 2:管理员
}

func (Group) TableName() string {
	return "group_user"
}
