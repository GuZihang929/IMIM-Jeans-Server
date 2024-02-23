package system

import (
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (GroupController) DeleteGroupUser(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// 获取请求中的群ID和要被删除的用户ID
	userID := c.Query("user_id")
	groupID := c.Query("group_id")
	if err := c.ShouldBindQuery(&userID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	if err := c.ShouldBindQuery(&groupID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	userid, _ := strconv.ParseInt(userID, 10, 64)
	groupid, _ := strconv.ParseInt(groupID, 10, 64)
	// 验证是否有权操作此群
	isOwner, err := user.IsGroupGM(claims.UserID, groupid)
	if err != nil || !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此群"})
		return
	}
	// 将新群信息插入数据库
	err = user.DeleteGroupUser(groupid, userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除群员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除群员成功"})
}
