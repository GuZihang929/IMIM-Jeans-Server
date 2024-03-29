package browser

import (
	"IM-Server/global"
	"IM-Server/im"
	"IM-Server/im/conn"
	_json "IM-Server/im/message/json"
	"IM-Server/model/system"
	"IM-Server/utils"
	"context"
	"encoding/json"
	"fmt"

	"strconv"
	"sync"
	"time"
)

type DefaultBrowserManager struct {
	browsers      *browsers
	browserOnline int64
	messageSent   int64
	maxOnline     int64
	startAt       int64
}

var DefaultManager *DefaultBrowserManager

func NewDefaultManager() *DefaultBrowserManager {
	ret := new(DefaultBrowserManager)
	ret.browsers = newBrowsers()
	ret.startAt = time.Now().Unix()
	return ret
}

func SetDefaultManager(manager *DefaultBrowserManager) {
	DefaultManager = manager
}

type browsers struct {
	m    sync.RWMutex
	bros map[int64]*Browser
}

func newBrowsers() *browsers {
	ret := new(browsers)
	ret.m = sync.RWMutex{}
	ret.bros = make(map[int64]*Browser)
	return ret
}

func newBrowser(conn conn.Connection) *Browser {
	bro := new(Browser)
	bro.conn = conn
	bro.state = stateRunning
	bro.messages = make(chan *_json.ComMessage, 40)
	bro.connectAt = time.Now()

	bro.hbR = time.NewTimer(HeartbeatDuration)
	bro.hbW = time.NewTimer(HeartbeatDuration)
	return bro
}

func (c *DefaultBrowserManager) BrowserConnected(ctx context.Context, conn conn.Connection) int64 {

	// 获取一个临时 uid 标识这个连接
	//snowflake := until.NewSnowflake()
	//var connUid int64 = snowflake.NextVal()
	id := ctx.Value("userId").(string)
	connUid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		global.Logger.Error("长连接，请求头数据格式有问题")
	}

	//将userId保存到redis，value为 connUid

	ret := newBrowser(conn)
	ret.SetID(connUid)
	c.browsers.bros[connUid] = ret

	//建立长连接的信息初始化
	InitWebSocket(connUid)

	//处理离线消息,创建会话列表
	ret.OfflineHandel(connUid)
	// 开始处理连接的消息
	ret.Run()
	return connUid
}

type MessageHandler func(message *_json.ComMessage) error

// messageHandleFunc 所有客户端消息都传递到该函数处理
var messageHandleFunc MessageHandler = nil

func SetMessageHandler(handler MessageHandler) {
	messageHandleFunc = func(message *_json.ComMessage) error {
		err := handler(message)
		if err != nil {
			fmt.Println("Message 设置handel出错：", err)
		}
		return err
	}
}

func (c *DefaultBrowserManager) GetBrowser(count int64) *Browser {

	c.browsers.m.RLock()
	bro := c.browsers.bros[count]
	c.browsers.m.RUnlock()
	return bro
}

func Handel(message *_json.ComMessage) error {

	if message.Action == "1" {
		getBrowser := DefaultManager.GetBrowser(message.Sender)

		mess := &_json.ComMessage{Action: "1", Ver: 1, Message: ""}
		getBrowser.messages <- mess

		return nil
	}

	switch message.Ver {
	case 0: //一对一

		//检测用户是否在线，在线直接发送，并且保存数据，不在线，创建接收者不在线map，
		//保存此时数据库中与发送者最大的seq，并保存数据。

		//检测用户是否在线
		exists, err := global.Redis.Exists(context.Background(), im.GetRedisKeyMain(message.Receiver)).Result()
		if err != nil {
			global.Logger.Error("用户在线检测出错，err:" + err.Error())
			return err
		}

		//需要创建session结构体，放入redis
		marshalJSON, _ := message.MarshalJSON()
		err2 := global.Redis.HSet(context.Background(), im.GetRedisKeyUserSessionMess(message.Receiver), message.Sender, string(marshalJSON)).Err()
		if err2 != nil {
			global.Logger.Error("消息放入队列出错，err:" + err2.Error())
			return err2
		}

		_, err = global.Redis.HIncrBy(context.Background(), im.GetRedisKeyUserSessionNum(message.Receiver), strconv.Itoa(int(message.Sender)), 1).Result()
		if err != nil {
			global.Logger.Error("redis自增，err:" + err.Error())
			return err
		}

		result, err := global.Redis.HGet(context.Background(), im.GetRedisKeyUserSessionMess(message.Sender), strconv.Itoa(int(message.Receiver))).Result()
		var mess _json.ComMessage
		err = json.Unmarshal([]byte(result), &mess)
		mess.Message = message.Message
		mess.Time = message.Time

		bytes, err := mess.MarshalJSON()
		if err != nil {
			panic(err)
		}
		err2 = global.Redis.HSet(context.Background(), im.GetRedisKeyUserSessionMess(message.Sender), message.Receiver, string(bytes)).Err()
		if err2 != nil {
			global.Logger.Error("消息放入队列出错，err:" + err2.Error())
			return err2
		}

		if exists == 1 {
			//在线

			//将消息发送到接收者通道中
			getBrowser := DefaultManager.GetBrowser(message.Receiver)
			getBrowser.messages <- message

		}

		comm := &system.ChatPrivate{
			SRId:    utils.MergeId(message.Sender, message.Receiver),
			SId:     message.Sender,
			RId:     message.Receiver,
			Context: message.Message,
			Time:    message.Time,
			Seq:     message.Seq,
			Type:    message.Type,
		}

		//处理数据
		global.DB.Create(&comm)

	case 1: //群聊

		jsonBytes, err := json.Marshal(*message)
		if err != nil {
			global.Logger.Error("JSON 序列化错误:" + err.Error())
			return err
		}

		switch message.Action {

		case "0":
			//用户不在线，将消息保存到redis中（redisKeyUserOfflineMessNum）
			offlineUser, err := global.Redis.SDiff(context.Background(), im.GetRedisKeyGroupAllUser(message.Receiver), im.GetRedisKeyGroupOnlineUser(message.Receiver)).Result()
			if err != nil {
				global.Logger.Error("获取群不在线用户：" + err.Error())
				return err
			}
			//用户量大时的效率怎么样？
			for _, s := range offlineUser {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					global.Logger.Error("数据转换出错，err:" + err.Error())
					return err
				}
				err2 := global.Redis.HSet(context.Background(), im.GetRedisKeyUserSessionMess(i), message.Sender, string(jsonBytes)).Err()
				if err2 != nil {
					global.Logger.Error("消息放入队列出错，err:" + err2.Error())
					return err2
				}

				_, err = global.Redis.HIncrBy(context.Background(), im.GetRedisKeyUserSessionNum(i), strconv.Itoa(int(message.Sender)), 1).Result()
				if err != nil {
					global.Logger.Error("redis自增，err:" + err.Error())
					return err
				}
			}

			//获取群号：message.Receiver,获取此群的频道，发布消息。
			global.Redis.Publish(context.Background(), im.GetRedisKeyGroupChannel(message.Receiver), jsonBytes)

		case "2":

			allAdmin, err := global.Redis.SMembers(context.Background(), im.GetRedisKeyGroupAdmin(message.Receiver)).Result()
			if err != nil {
				global.Logger.Error("获取群All管理员：" + err.Error())
				return err
			}

			//获取群管理员id，发送消息
			admin, err := global.Redis.SInter(context.Background(), im.GetRedisKeyGroupOnlineUser(message.Receiver), im.GetRedisKeyGroupAdmin(message.Receiver)).Result()
			if err != nil {
				global.Logger.Error("获取群在线管理员：" + err.Error())
				return err
			}
			temp := make(map[string]struct{})

			for _, s := range admin {
				temp[s] = struct{}{}
			}

			for _, s := range allAdmin {
				i, err := strconv.ParseInt(s, 10, 64)

				if _, ok := temp[s]; ok {
					if err != nil {
						global.Logger.Error("数据转换出错，err:" + err.Error())
						return err
					}

					//将消息发送到接收者通道中
					getBrowser := DefaultManager.GetBrowser(i)
					getBrowser.messages <- message
				} else {
					//不在线管理员

					err2 := global.Redis.HSet(context.Background(), im.GetRedisKeyUserSessionMess(i), message.Sender, string(jsonBytes)).Err()
					if err2 != nil {
						global.Logger.Error("消息放入队列出错，err:" + err2.Error())
						return err2
					}

					_, err = global.Redis.HIncrBy(context.Background(), im.GetRedisKeyUserSessionNum(i), strconv.Itoa(int(message.Sender)), 1).Result()
					if err != nil {
						global.Logger.Error("redis自增，err:" + err.Error())
						return err
					}
				}

			}
		}
		comm := &system.Communication{
			FromID:  message.Sender,
			GroupID: message.Receiver,
			Content: string(jsonBytes),
			Time:    time.Now().Unix(),
			Seq:     message.Seq,
			Type:    message.Type,
			Class:   "0",
		}
		//处理数据
		global.DB.Create(&comm)

	}

	return nil
}
