package system

import (
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type GroupController struct {
}

func (GroupController) CreateGroup(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		log.Println("Token验证出错:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var info system.Group
	info.UserId = claims.UserID
	info.JoinInTime = time.Now()
	info.GroupId = GenerateRandomGroupID()
	info.Identity = 1
	// 将新群信息插入数据库
	err = user.InsertGroup(info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "群创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "群创建成功"})
}

// GenerateRandomGroupID 生成一个随机的群组ID
func GenerateRandomGroupID() int64 {
	rand.Seed(time.Now().UnixNano())  // 设置随机种子以确保每次运行都得到不同的随机数
	return rand.Int63n((1 << 48) - 1) // 假设群号为48位整数（最大值约为2^48-1），可根据实际需求调整范围
}
