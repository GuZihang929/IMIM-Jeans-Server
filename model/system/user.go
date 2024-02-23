package system

import "time"

type User struct {
	UserID   int64     `json:"uid" gorm:"column:uid;primaryKey"` // 用户ID
	Account  string    `json:"account"`                          // 用户账号
	Nickname string    `json:"nickname"`                         // 用户昵称
	Password string    `json:"password"`                         // 用户密码
	Avatar   string    `json:"avatar"`                           // 用户头像
	CreateAt time.Time `json:"create_at"`                        // 创建时间
	UpdateAt time.Time `json:"update_at"`                        // 更新时间
}

func (User) TableName() string {
	return "im_user"
}
