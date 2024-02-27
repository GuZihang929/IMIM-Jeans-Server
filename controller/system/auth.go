package system

import (
	"IM-Server/global"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func (c AuthController) DelSession(ctx *gin.Context) {

	id := 0
	// 获取会话id
	err := ctx.ShouldBindJSON(&id)
	if err != nil {
		global.Logger.Error("")
	}
}
