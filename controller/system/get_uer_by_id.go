package system

import (
	"IM-Server/service/user"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetUserByUserID 用于根据用户ID获取用户信息
func (GroupController) GetUserByUserID(c *gin.Context) {
	// 从 URL 参数中获取用户 ID
	userIDStr := c.Query("userID")
	if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
		// 查询数据库获取用户信息
		userinfo, err := user.GetUserByUserID(userID)
		if err != nil {
			// 数据库查询错误
			c.JSON(500, gin.H{"error": "Failed to get user information from the database"})
			return
		}

		// 用户不存在
		if userinfo == nil {
			c.JSON(404, gin.H{"message": "User not found"})
			return
		}

		// 返回用户信息
		c.JSON(200, userinfo)
	} else {
		// 用户ID参数解析错误
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
	}
}
