package model

import _json "IM-Server/im/message/json"

type Session struct {
	_json.Data
	Id      int64  `json:"id,omitempty"` //会话对象id
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Message string `json:"message,omitempty"`
	Num     int64  `json:"num,omitempty"`
	Time    int64  `json:"time,omitempty"`
	Ver     int64  `json:"ver"`
}
