package system

import (
	"IM-Server/global"
	"IM-Server/im"
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
