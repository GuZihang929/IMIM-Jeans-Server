package system

type ChatPrivate struct {
	Id      int64  `gorm:"id" json:"id,omitempty"`
	SRId    int64  `gorm:"sr_id" json:"s_r_id,omitempty"`
	SId     int64  `gorm:"s_id" json:"s_id,omitempty"`
	RId     int64  `gorm:"r_id" json:"r_id,omitempty"`
	Context string `gorm:"context" json:"context,omitempty"`
	Time    int64  `gorm:"time" json:"time"`
	Type    int64  `gorm:"type" json:"type,omitempty"`
	Seq     int64  `gorm:"seq" json:"seq"`
}

func (ChatPrivate) TableName() string {
	return "chat_private"
}
