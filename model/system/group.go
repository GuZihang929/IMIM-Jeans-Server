package system

import "time"

type Group struct {
	Id        int64     `gorm:"id" json:"id"`
	GId       int64     `gorm:"g_id" json:"g_id,omitempty"`
	GName     string    `gorm:"g_name" json:"g_name,omitempty"`
	GUrl      string    `gorm:"g_url" json:"g_url,omitempty"`
	GCreatId  int64     `gorm:"g_creat_id" json:"g_creat_id,omitempty"`
	CreatTime time.Time `gorm:"creat_time" json:"creat_time"`
}
