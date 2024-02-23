package messageHandel

import (
	"IM-Server/im/browser"
	_json "IM-Server/im/message/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func ListeningGroupHandel(ch *<-chan *redis.Message, userId, channelId int64) {
	bro := browser.DefaultManager.GetBrowser(channelId)
	fmt.Println("监听开始")
	for {
		fmt.Println("------")
		select {
		case msg := <-*ch:
			fmt.Println("就收到消息")
			// 发送到bro中
			cm := &_json.ComMessage{
				Sender:   userId,
				Receiver: channelId,
				Ver:      1,
				Seq:      0,
				Action:   "",
				Data:     _json.NewData(msg),
				Extra:    nil,
			}
			bro.GetMessageChan() <- cm
		}
	}
}
