package browser

import (
	"IM-Server/global"
	"IM-Server/im"
	_json "IM-Server/im/message/json"
	"context"
	"fmt"
	"strconv"
)

func InitWebSocket(id int64) {

	//上线人数++
	err := global.Redis.IncrBy(context.Background(), im.GetRedisKeyOnlineNum(), 1).Err()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("在线人数添加出错，err:%s", err.Error()))
	}
	//订阅用户的群，并将用户放入在线集合，
	fmt.Println("1111111")
	groups, err := global.Redis.SMembers(context.Background(), im.GetRedisKeyGroup(id)).Result()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("获取%d用户的所有群id错误，err:%d", id, err))
	}
	fmt.Println(groups)
	fmt.Println("1111111")

	for _, group := range groups {
		fmt.Println(group)

		groupId, _ := strconv.ParseInt(group, 10, 64)
		channel := global.Redis.Subscribe(context.Background(), im.GetRedisKeyGroupChannel(groupId))
		_, err = global.Redis.SAdd(context.Background(), im.GetRedisKeyGroupOnlineUser(groupId), id).Result()
		if err != nil {
			global.Logger.Error("将用户添加到群在线集合中出错，err:" + err.Error())
		}
		ch := channel.Channel()
		//go messageHandel.ListeningGroupHandel(&messages, b.User, b.Group)
		go func() {
			bro := DefaultManager.GetBrowser(id)
			for {
				select {
				case msg := <-ch:
					// 发送到bro中
					cm := &_json.ComMessage{
						Sender:   id,
						Receiver: groupId,
						Ver:      1,
						Seq:      0,
						Action:   "",
						Data:     _json.NewData(msg),
						Extra:    nil,
					}
					bro.GetMessageChan() <- cm
				}
			}
		}()
	}

}

func FinishWebSocket(id int64) {
	//下线人数--
	err := global.Redis.IncrBy(context.Background(), im.GetRedisKeyOnlineNum(), -1).Err()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("在线人数减少出错，err:%s", err.Error()))
	}

	global.Redis.Del(context.Background(), im.GetRedisKeyMain(id))

	//用户下线，删除在线群的id

	groups, err := global.Redis.SMembers(context.Background(), im.GetRedisKeyGroup(id)).Result()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("获取%d用户的所有群id错误，err:%d", id, err))
	}

	for _, group := range groups {
		groupId, _ := strconv.ParseInt(group, 10, 64)
		_, err = global.Redis.SRem(context.Background(), im.GetRedisKeyGroupOnlineUser(groupId), id).Result()
		if err != nil {
			global.Logger.Error("将用户从群在线集合取消出错，err:" + err.Error())
		}
	}
}

func SendMessage(sender, receiver, ver, seq int64, action string, data interface{}) {

	message := &_json.ComMessage{
		Sender:   sender,
		Receiver: receiver,
		Ver:      ver,
		Seq:      seq,
		Action:   action,
		Data:     data.(_json.Data),
		Extra:    nil,
	}

	browser := DefaultManager.GetBrowser(message.Receiver)

	browser.GetMessageChan() <- message

}
