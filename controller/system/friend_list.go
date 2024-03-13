package system

import (
	"IM-Server/model/system"
	"IM-Server/res"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetFriendList 获取用户的好友列表
//func (FriendController) GetFriendList(c *gin.Context) {
//	token := c.GetHeader("token")
//	claims, err := jwts.ParseToken(token)
//	if err != nil {
//		log.Println("Token验证出错:", err)
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
//		return
//	}
//	userId := claims.UserID
//
//	friendList, err := user.GetFriendList(userId)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友列表失败"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"data": friendList})
//}

// GetFriendList 获取用户的好友列表及好友详情
func (FriendController) GetFriendList(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId := claims.UserID

	// 获取好友ID列表
	friendIDs, err := user.GetFriendIDs(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取好友ID列表失败"})
		return
	}

	// 初始化存储好友详细信息的切片
	var friendDetails []FriendWithUserInfo
	for _, friendID := range friendIDs {
		// 根据每个好友ID获取好友详细信息
		friendInfo, err := user.GetUserByUserID(friendID)
		if err != nil {
			// 如果查询单个好友信息时出错，可以选择记录错误并跳过，或直接返回错误
			log.Printf("获取用户ID为 %d 的好友信息时出错: %v", friendID, err)
			continue
		}

		// 将获取到的好友信息添加到最终结果中
		friendDetails = append(friendDetails, FriendWithUserInfo{UserInfo: *friendInfo})
	}

	//登录成功，返回用户信息及生成的token
	res.OkWithData(map[string]interface{}{
		"user": friendDetails,
	}, c)
}

type FriendWithUserInfo struct {
	UserInfo system.User
}
