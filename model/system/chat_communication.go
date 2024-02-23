package system

import "time"

type Communication struct {
	ID        uint      `gorm:"primarykey"`
	FromID    int64     `gorm:"from_id;not null;comment:'发送人id'"`
	FromName  string    `gorm:"from_name;not null;comment:'发送人name'"`
	ToID      int64     `gorm:"to_id;not null;comment:'接收人id，不适用于群消息'"`
	ToName    string    `gorm:"to_name;not null;comment:'接收人name不适用于群消息'"`
	Content   string    `gorm:"content;not null;comment:'消息内容'"`
	Time      time.Time `gorm:"time;not null;comment:'时间'"`
	GroupID   int64     `gorm:"group_id;comment:'群id'"`
	GroupName int64     `gorm:"group_name;comment:'群名称'"`
	IsRead    int64     `gorm:"is_read;default:0;comment:'是否已读，不适用于群消息'"`
	Type      int64     `gorm:"type;default:1;comment:'消息类型：1是普通文本，2是图片，3是语音'"`
	Class     string    `gorm:"class;comment:'消息类：1是用户聊天，2是群组聊天'"`
	Seq       int64     `gorm:"seq"`
}

func (Communication) TableName() string {
	return "chat_communication"
}
