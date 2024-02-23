package system

import (
	"IM-Server/global"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s SysPublicController) VisitorLogout(c *gin.Context) {
	// 从请求头或cookie中获取JWT token
	token := c.GetHeader("Authorization") // 假设token以Bearer方式放在Authorization头中，如"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
	if token == "" {
		// 如果未在请求头中找到，则尝试从cookie或其他位置获取token
		// 这里省略了从cookie获取token的代码，具体取决于您的实现
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
		return
	}

	// 删除Redis中与该token关联的游客信息
	err := global.Redis.Del(context.Background(), token).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
