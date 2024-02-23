package system

import (
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FriendController struct {
}

// CreateFriend 创建朋友关系
func (FriendController) CreateFriend(c *gin.Context) {
	var info system.Friend

	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	friendId := c.Query("friend_id")
	if err := c.ShouldBind(&friendId); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	info.UserId = claims.UserID
	info.FriendId, _ = strconv.ParseInt(friendId, 10, 64)
	info.Status = 1 // 0 已删除 1 添加好友
	info.CreateTime = time.Now()
	info.UpdateTime = time.Now()

	// 查询用户信息
	userInfo, err := user.GetUserByUserID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}
	if userInfo == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号不存在"})
		return
	}

	// 查询用户信息
	friendInfo, err := user.GetUserByUserID(info.FriendId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}
	if friendInfo == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "添加好友的账号不存在"})
		return
	}
	// 将新用户插入数据库
	err = user.InsertFriend(info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "好友建立失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "好友创建成功"})
}
