package system

import (
	"IM-Server/global"
	"IM-Server/im"
	"IM-Server/model/system"
	"IM-Server/utils"
	"IM-Server/utils/jwts"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type AuthController struct {
}

func (c AuthController) DelSession(ctx *gin.Context) {

	token := ctx.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := struct {
		Id int `json:"id"`
	}{}
	// 获取会话id
	err = ctx.ShouldBindJSON(&userId)
	if err != nil {
		global.Logger.Error("解析参数出错，" + err.Error())
	}

	_, err = global.Redis.HDel(context.Background(), im.GetRedisKeyUserSessionMess(claims.UserID), strconv.Itoa(userId.Id)).Result()
	if err != nil {
		global.Logger.Error("删除会话列表，" + err.Error())
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "删除失败",
		})
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func (c AuthController) DelSessionNum(ctx *gin.Context) {

	token := ctx.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := struct {
		Id int `json:"id"`
	}{}
	// 获取会话id
	err = ctx.ShouldBindJSON(&userId)
	if err != nil {
		global.Logger.Error("解析参数出错，" + err.Error())
	}

	_, err = global.Redis.HDel(context.Background(), im.GetRedisKeyUserSessionNum(claims.UserID), strconv.Itoa(userId.Id)).Result()
	if err != nil {
		global.Logger.Error("删除会话数目，" + err.Error())
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "删除失败",
		})
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

func (c AuthController) GetHistoricalNew(ctx *gin.Context) {

	token := ctx.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	message := struct {
		Id    int64 `json:"id"`
		Ver   int64 `json:"ver"`
		Seq   int64 `json:"seq"`
		Time  int64 `json:"time"`
		Total int   `json:"total"`
	}{}
	// 获取id
	err = ctx.ShouldBindJSON(&message)
	if err != nil {
		global.Logger.Error("解析参数出错，" + err.Error())
	}

	if message.Ver == 0 {
		// 查询数据库中聊天记录 10条。

		SRId := utils.MergeId(claims.UserID, message.Id)

		cp := []system.ChatPrivate{}
		global.DB.Where("sr_id = ? and time <= ?", SRId, message.Time).Limit(message.Total).Find(&cp)

		ctx.JSON(200, gin.H{
			"code": 200,
			"data": cp,
		})
	} else if message.Ver == 1 {
		// 群消息
	}

}
