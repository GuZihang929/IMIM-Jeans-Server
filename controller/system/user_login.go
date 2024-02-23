package system

import (
	"IM-Server/global"
	"IM-Server/res"
	"IM-Server/service/user"
	"IM-Server/utils/jwts"
	"IM-Server/utils/pwd"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SysPublicController struct {
}

type LoginRequest struct {
	UserName string `json:"account" binding:"required" msg:"请输入用户名（注册邮箱号）"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
}

func (SysPublicController) Login(c *gin.Context) {
	var cr LoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	// 查询用户信息
	userInfo, err := user.GetUserByAccount(cr.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "数据库查询失败"})
		return
	}

	if userInfo == nil {
		c.JSON(http.StatusOK, gin.H{"error": "账号不存在"})
		return
	}

	isCheck := pwd.CheckPwd(userInfo.Password, cr.Password)
	if !isCheck {
		global.Logger.Warn("用户名或密码错误")
		c.JSON(http.StatusOK, gin.H{"error": "用户名或者密码错误"})
		return
	}

	// 登录成功，生成token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: userInfo.Nickname,
		UserID:   userInfo.UserID,
		Avatar:   userInfo.Avatar,
		Username: userInfo.Account,
	})

	if err != nil {
		global.Logger.Error(err.Error())
	}

	res.OkWithData(token, c)
}
