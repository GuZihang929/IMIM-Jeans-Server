package browser

import (
	"IM-Server/global"
	"IM-Server/im"
	"IM-Server/im/conn"
	_json "IM-Server/im/message/json"
	"IM-Server/im/message/model"
	"IM-Server/model/system"
	"IM-Server/utils/timingwheel"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

var tw = timingwheel.NewTimingWheel(time.Millisecond*500, 3, 20)

// HeartbeatDuration 心跳间隔
const HeartbeatDuration = time.Second * 20

const (
	_ = iota
	stateRunning
	stateClosing
	stateClosed
)

type Browser struct {
	conn conn.Connection

	// id 唯一标识
	id int64

	connectAt time.Time
	// state browser 状态
	state int32

	// messages 带缓冲的下行消息管道, 缓冲大小40
	messages chan *_json.ComMessage

	// rCloseCh 关闭或写入则停止读
	rCloseCh   chan struct{}
	readClosed int32

	// hbR 心跳倒计时
	hbR    *time.Timer
	hbLost int

	hbW *time.Timer
}

func NewBrowser(conn conn.Connection) *Browser {
	browser := new(Browser)
	browser.conn = conn

	browser.state = stateRunning
	browser.connectAt = time.Now()
	browser.messages = make(chan *_json.ComMessage, 40)

	browser.hbR = time.NewTimer(HeartbeatDuration)
	browser.hbW = time.NewTimer(HeartbeatDuration)
	return browser
}

// SetID 设置 id 标识及设备标识
func (c *Browser) SetID(id int64) {
	atomic.StoreInt64(&c.id, id)
}

func (c *Browser) IsRunning() bool {
	return atomic.LoadInt32(&c.state) == stateRunning
}

func (c *Browser) GetID() int64 {
	return atomic.LoadInt64(&c.id)
}

func (c *Browser) Run() {
	go c.readMessage()
	go c.writeMessage()
}

// Exit 退出客户端
func (c *Browser) Exit() {
	s := atomic.LoadInt32(&c.state)
	if s == stateClosed || s == stateClosing {
		return
	}
	atomic.StoreInt32(&c.state, stateClosing)

	if atomic.LoadInt32(&c.readClosed) != 1 {
		c.rCloseCh <- struct{}{}
	}

}

func (c *Browser) GetMessageChan() chan *_json.ComMessage {
	return c.messages
}

// readMessage 开始从 Connection 中读取消息
func (c *Browser) readMessage() {
	reader := defaultReader{}
	readChan, done := reader.ReadCh(c.conn)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for {
		select {

		case <-c.rCloseCh:
			close(c.rCloseCh)
			goto STOP

		case <-c.hbR.C:

			c.hbLost++
			if c.hbLost > 3 {
				goto STOP
			}
			fmt.Printf("心跳检测超时%d次\n", c.hbLost)
			c.hbR.Reset(HeartbeatDuration)
			//心跳检测返回消息

		case msg := <-readChan:

			if msg.err != nil {
				if !c.IsRunning() || msg.err.Error() == "closed" {
					// 连接断开或致命错误中断读消息
					goto STOP
				}
				global.Logger.Error("readChan通道接收消息出错：" + msg.err.Error())
				continue
			}

			c.hbLost = 0
			c.hbR.Reset(HeartbeatDuration)

			//id := c.id
			// 统一处理消息函数

			messageHandleFunc(msg.m)
			msg.Recycle()
		}
	}
STOP:
	c.hbR.Stop()
	close(done)
	id := c.GetID()
	//下线，取消订阅，删除用户信息
	FinishWebSocket(id)
	global.Logger.Info(fmt.Sprintf("client read closed, id=%d", id))
}

// writeMessage 开始向 Connection 中写入消息队列中的消息
func (c *Browser) writeMessage() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for {
		select {

		case <-c.hbW.C:
			if !c.IsRunning() {
				global.Logger.Error(fmt.Sprintf("read closed, down msg queue timeout, close write now, uid=%d", c.id))
				goto STOP
			}

			c.hbW.Reset(HeartbeatDuration)

		case m := <-c.messages:

			b, err := json.Marshal(*m)
			if err != nil {
				fmt.Println(err)
				continue
			}

			c.hbLost = 0
			c.hbW.Reset(HeartbeatDuration)
			fmt.Println("消息写入")
			fmt.Println(b)
			err = c.conn.Write(b)
		}
	}
STOP:
	c.Exit()
	atomic.StoreInt32(&c.state, stateClosed)
	close(c.messages)
	_ = c.conn.Close()
}

// OfflineHandel 处理离线消息
func (c *Browser) OfflineHandel(key int64) {
	// 获取redis中用户的会话哈希
	result, err := global.Redis.HGetAll(context.Background(), im.GetRedisKeyUserSessionMess(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return
		}
		global.Logger.Error("获取离线消息出错，err" + err.Error())
	}
	fmt.Println("redis中的会话哈希，", result)

	// 创建会话列表

	sessions := model.Sessions{}

	// 解析并将此结构按照时间戳降序排列。
	for _, value := range result {
		message := &_json.ComMessage{}
		err = json.Unmarshal([]byte(value), message)
		if err != nil {
			global.Logger.Error("消息反序列化失败，err:" + err.Error())
			break
		}
		s, err := global.Redis.Get(context.Background(), im.GetRedisKeyMain(message.Receiver)).Result()
		if err != nil {
			if err == redis.Nil {
				return
			}
			global.Logger.Error("获取用户信息出错，err" + err.Error())
		}
		user := &system.User{}
		err = json.Unmarshal([]byte(s), user)
		if err != nil {
			global.Logger.Error("消息反序列化失败，err:" + err.Error())
			break
		}

		num, err := global.Redis.Get(context.Background(), im.GetRedisKeyUserSessionNum(message.Receiver)).Result()
		if err != nil {
			if err == redis.Nil {
				return
			}
			global.Logger.Error("获取消息数目出错，err" + err.Error())
		}

		i, err := strconv.ParseInt(num, 10, 64)

		session := model.Session{
			Id:      user.UserID,
			Name:    user.Nickname,
			Avatar:  user.Avatar,
			Message: message.Message,
			Num:     i,
			Time:    message.Time,
		}
		sessions = append(sessions, session)

	}

	// 将消息按时间戳排序
	sort.Sort(sessions)

	mess := &_json.ComMessage{
		Session: sessions,
	}

	c.messages <- mess
}
