package system

import (
	"IM-Server/global"
	"IM-Server/model/system"
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

	type Body struct {
		UsersId []int64 `json:"Users_id"`
	}

	var body Body

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  200,
			"error": "body格式出错",
		})
		return
	}
	var info system.GroupUser

	begin := global.DB.Begin()
	group := GenerateRandomGroupID()
	info.UserId = claims.UserID
	info.JoinInTime = time.Now()
	info.GroupId = group
	info.Identity = 1
	// 将新群信息插入数据库
	err = begin.Create(&info).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "群创建失败"})
		begin.Rollback()
		return
	}
	for _, i2 := range body.UsersId {
		info.UserId = i2
		info.JoinInTime = time.Now()
		info.GroupId = group
		info.Identity = 0
		err = begin.Create(&info).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "群创建失败"})
			begin.Rollback()
			return
		}
	}
	begin.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "群创建成功"})
}

// GenerateRandomGroupID 生成一个随机的群组ID
func GenerateRandomGroupID() int64 {
	rand.Seed(time.Now().UnixNano())  // 设置随机种子以确保每次运行都得到不同的随机数
	return rand.Int63n((1 << 20) - 1) // 假设群号为48位整数（最大值约为2^48-1），可根据实际需求调整范围
}
