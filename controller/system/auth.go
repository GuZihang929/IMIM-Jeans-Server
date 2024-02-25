package system

import (
	"IM-Server/model/system"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func (c AuthController) Info(ctx *gin.Context) {
	//访问service

	user := system.User{
		UserID:   5667123,
		Account:  "499@qq.com",
		Nickname: "滴滴滴",
		Avatar:   "https://www.freeimg.cn/i/2024/01/30/65b84d987ee13.png",
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": user,
		"msg":  "上传成功",
	})
}

type Body struct {
	Group int64 `json:"group"`
	User  int64 `json:"user"`
}

//func (c AuthController) World(ctx *gin.Context) {
//
//	// 游客订阅world频道
//	// 订阅频道
//
//	b := Body{}
//
//	err := ctx.ShouldBindJSON(&b)
//	if err != nil {
//		global.Logger.Error("世界频道id解析出错:" + err.Error())
//	}
//	fmt.Println(b)
//
//	channel := global.Redis.Subscribe(context.Background(), strconv.FormatInt(b.Group, 10))
//	fmt.Println("订阅频道：", strconv.FormatInt(b.Group, 10))
//	_, err = global.Redis.SAdd(context.Background(), im.GetRedisKeyGroupOnlineUser(b.Group), b.User).Result()
//	if err != nil {
//		global.Logger.Error("将用户添加到群在线集合中出错，err:" + err.Error())
//	}
//
//	//channelId, err := strconv.ParseInt(world, 10, 64)
//	//user, err := strconv.ParseInt(userId, 10, 64)
//	if err != nil {
//		global.Logger.Error("世界频道id格式出错")
//	}
//	ch := channel.Channel()
//	//go messageHandel.ListeningGroupHandel(&messages, b.User, b.Group)
//	go func() {
//		bro := browser.DefaultManager.GetBrowser(b.User)
//		for {
//			select {
//			case msg := <-ch:
//				// 发送到bro中
//				cm := &_json.ComMessage{
//					Sender:   b.User,
//					Receiver: b.Group,
//					Ver:      1,
//					Seq:      0,
//					Action:   "",
//					Message:     "",
//					Extra:    nil,
//				}
//				bro.GetMessageChan() <- cm
//			}
//		}
//	}()
//	ctx.JSON(200, gin.H{
//		"code": 200,
//		"data": "",
//		"msg":  "上传成功",
//	})
//}
