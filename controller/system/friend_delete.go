package system

import (
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// DeleteFriend 删除朋友关系
func (FriendController) DeleteFriend(c *gin.Context) {
	var info system.Friend

	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	friendIdStr := c.Query("friend_id")
	if err := c.ShouldBindQuery(&friendIdStr); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	info.UserId = claims.UserID
	info.FriendId, _ = strconv.ParseInt(friendIdStr, 10, 64)

	// 验证好友关系是否存在
	exists, err := user.CheckFriendship(info.UserId, info.FriendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "好友关系不存在"})
		return
	}

	// 删除好友关系
	err = user.DeleteFriend(info.UserId, info.FriendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除好友失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "好友删除成功"})
}
