package system

import (
	"IM-Server/global"
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (GroupController) JoinGroup(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	groupID := c.Query("group_id")
	var info system.GroupUser
	info.UserId = claims.UserID
	info.JoinInTime = time.Now()
	info.GroupId, _ = strconv.ParseInt(groupID, 10, 64)
	info.Identity = 0
	//检查群友是否已经在群中（不可重复加群）
	if user.IsUserInGroup(info.UserId, info.GroupId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已在群中"})
		return
	}
	// 将新群信息插入数据库
	err = user.InsertGroup(info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "群加入失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "群加入成功"})
}

func (GroupController) GetGroup(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var gp []system.GroupUser
	var infos []system.Group

	global.DB.Where("user_id = ?", claims.UserID).Find(&gp)
	infos = make([]system.Group, len(gp))

	for i, groupUser := range gp {
		global.DB.Where("g_id = ?", groupUser.GroupId).First(&infos[i])
	}

	c.JSON(200, gin.H{
		"code":    200,
		"data":    infos,
		"massage": "群列表",
	})

}
