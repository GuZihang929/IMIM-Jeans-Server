package system

import (
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetFriendList 获取用户的好友列表
func (FriendController) GetFriendList(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId := claims.UserID

	friendList, err := user.GetFriendList(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": friendList})
}
