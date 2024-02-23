package json

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	des interface{} `json:"des"`
}

func NewData(d interface{}) Data {
	return Data{
		des: d,
	}
}

func NewEtmData() Data {
	return Data{}
}

func (d *Data) Data() interface{} {
	return d.des
}

func (d *Data) UnmarshalJSON(bytes []byte) error {
	d.des = bytes
	return nil
}

func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.des)
}

func (d *Data) bytes() []byte {
	bytes, ok := d.des.([]byte)
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
	s, ok := d.des.([]byte)
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
	Action   string            `json:"Action"`   //信号类型, 0为消息，1为心跳检测
	Data     Data              `json:"Data"`     //消息数据
	Extra    map[string]string `json:"Extra"`    //额外信息
}

func NewMessage(seq int64, action string, data interface{}) *ComMessage {
	return &ComMessage{
		Ver:    0,
		Seq:    seq,
		Action: action,
		Data:   NewData(data),
	}
}

func NewEmptyMessage() *ComMessage {
	return &ComMessage{
		Ver:    0,
		Seq:    0,
		Action: "",
		Data:   NewEtmData(),
	}
}

func (d *ComMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
}
