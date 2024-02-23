package system

import (
	"IM-Server/model/system"
	"IM-Server/service/user"
	"IM-Server/utils/pwd"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserCreateRequest struct {
	NickName string `json:"nick_name" binding:"required" msg:"请输入昵称"` // 昵称
	UserName string `json:"account" binding:"required" msg:"请输入用户名"`  // 用户名（邮箱号）
	Password string `json:"password" binding:"required" msg:"请输入密码"`  // 密码
	Avatar   string `json:"avatar" msg:"请输入头像"`
	Code     string `json:"code" binding:"required" msg:"请输入验证码"` //验证码
}

func (s SysPublicController) RegisterUser(c *gin.Context) {
	var newUser UserCreateRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	NewUser := system.User{
		Account:  newUser.UserName,
		Nickname: newUser.NickName,
		Password: newUser.Password,
		Avatar:   newUser.Avatar,
	}

	// 检查是否存在相同账号（邮箱）的用户
	if user.IsUserExists(NewUser.Account) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已存在相同邮箱的用户"})
		return
	}

	// 在存储之前对用户的密码进行哈希处理
	// 在注册时对密码进行哈希处理
	hashedPassword, err := pwd.HashPwd(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希处理失败"})
		return
	}
	//newUser.Password = hashedPassword
	// 设置创建时间和更新时间戳
	now := time.Now()
	NewUser.CreateAt = now
	NewUser.UpdateAt = now

	NewUser.Password = hashedPassword
	// 假设 verificationCode 是用户在注册过程中提供的验证码
	if !user.VerifyVerificationCode(newUser.UserName, newUser.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}
	// 将新用户插入数据库
	err = user.InsertUser(NewUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户注册失败"})
		return
	}

	// 返回注册成功的信息
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功。"})
}
