package system

import "time"

type ChatGroups struct {
	Id      int64     `gorm:"id" json:"id,omitempty"`
	GId     int64     `gorm:"g_id" json:"g_id"`
	SId     int64     `gorm:"s_id" json:"s_id,omitempty"`
	Context string    `gorm:"context" json:"context,omitempty"`
	Time    time.Time `gorm:"time" json:"time"`
	Type    int64     `gorm:"type" json:"type,omitempty"`
	Seq     int64     `gorm:"seq" json:"seq"`
}

func (ChatGroups) TableName() string {
	return "chat_Groups"
}
