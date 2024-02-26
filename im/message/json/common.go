package json

import (
	"IM-Server/im/message/model"
	"encoding/json"
	"fmt"
)

type Data struct {
	Des interface{} `json:"Des"`
}

func NewData(d interface{}) Data {
	return Data{
		Des: d,
	}
}

func NewEtmData() Data {
	return Data{}
}

func (d *Data) Data() interface{} {
	return d.Des
}

func (d *Data) UnmarshalJSON(bytes []byte) error {
	d.Des = bytes
	return nil
}

func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Des)
}

func (d *Data) bytes() []byte {
	bytes, ok := d.Des.([]byte)
	if ok {
		return bytes
	}
	marshalJSON, err := d.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return marshalJSON
}

func (d *Data) Deserialize(i interface{}) error {
	s, ok := d.Des.([]byte)
	if ok {
		return json.Unmarshal(s, i)
	}
	return nil
}

type ComMessage struct {
	Sender   int64             `json:"Sender"`   //发送者
	Receiver int64             `json:"Receiver"` //接收者（个人长连接通道/群号）
	Ver      int64             `json:"Ver"`      //单发/群发，0/1
	Seq      int64             `json:"Seq"`      //序列号
	Type     int64             `json:"Type"`     //消息类型，1为普通文本，2为图像，3为语音。
	Action   string            `json:"Action"`   //信号类型, 0为消息，1为心跳检测，2为通知
	Message  string            `json:"Message"`  //消息
	Time     int64             `json:"Time"`     //消息时间
	Extra    map[string]string `json:"Extra"`    //额外信息
	Session  model.Sessions    `json:"Sessions"`
}

func NewEmptyMessage() *ComMessage {
	return &ComMessage{
		Sender:   0,
		Receiver: 0,
		Ver:      0,
		Seq:      0,
		Type:     0,
		Action:   "",
		Message:  "",
		Time:     0,
		Extra:    nil,
	}
}

func (d *ComMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(*d)
}

//func (d *ComMessage) DataToString() string {
//	p := make(map[string]interface{})
//	err := d.Data.Deserialize(&p)
//	if err != nil {
//		fmt.Println(err)
//	}
//	data, _ := json.Marshal(p["Des"])
//	return string(data)
//
//}

//func (d *ComMessage) ToData(b []byte) {
//	err := json.Unmarshal(b, d.Data.Des)
//	if err != nil {
//		fmt.Println(err, "-=-=-=-=")
//	}
//}
