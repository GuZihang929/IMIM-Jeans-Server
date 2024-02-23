package system

import (
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (GroupController) DeleteGroup(c *gin.Context) {
	token := c.GetHeader("token")
	fmt.Print(token)
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	//token获取群主的id（这里我设想用token来获取用户id检验用户id是否具有群主权限）
	groupOwnerID := claims.UserID
	groupID := c.Query("group_id")
	if err := c.ShouldBindQuery(&groupOwnerID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	if err := c.ShouldBindQuery(&groupID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	groupid, _ := strconv.ParseInt(groupID, 10, 64)
	// 验证群主是否有权操作此群
	isOwner, err := user.IsGroupOwner(groupOwnerID, groupid)
	if err != nil || !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此群"})
		return
	}
	// 将新群信息插入数据库
	err = user.DeleteGroup(groupid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解散群失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "解散群成功"})
}
