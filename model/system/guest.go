package system

import "time"

type Guest struct {
	ID       int64     `json:"uid"`       // 用户ID
	Nickname string    `json:"nickname"`  // 用户昵称
	Avatar   string    `json:"avatar"`    // 用户头像
	CreateAt time.Time `json:"create_at"` // 创建时间
	UpdateAt time.Time `json:"update_at"` // 更新时间
}
