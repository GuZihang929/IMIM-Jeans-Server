package model

import "IM-Server/im/message/json"

type Session struct {
	json.Data
	Id      int64       `json:"id,omitempty"` //会话对象id
	Name    string      `json:"name,omitempty"`
	Avatar  string      `json:"avatar,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Num     int64       `json:"num,omitempty"`
	Time    int64       `json:"time,omitempty"`
	Ver     int64       `json:"ver"`
}

type Sessions []Session

func (s Sessions) Len() int {

	return len(s)
}

func (s Sessions) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

func (s Sessions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewEmptySessions() Sessions {
	return make([]Session, 0)
}

func NewSessions(s Session) Sessions {
	return append(make([]Session, 0), s)
}
