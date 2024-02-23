package middleware

import (
	"IM-Server/res"
	"IM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
)

// JWTAuth
// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息
// 这里前端需要把token存储到cookie或者本地localStorage1
// 可以约定刷新令牌或者重新登录
//func JWTAuth() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.Request.Header.Get("Authorization")
//		if authHeader == "" {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2003,
//				"msg":  "请求头中auth为空",
//			})
//			c.Abort()
//			return
//		}
//
//		parts := strings.SplitN(authHeader, " ", 2)
//		if !(len(parts) == 2 && parts[0] == "Bearer") {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2004,
//				"msg":  "请求头中auth格式有误",
//			})
//			c.Abort()
//			return
//		}
//		mc, err := utils.ParseToken(parts[1])
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2005,
//				"msg":  "无效的Token",
//			})
//			c.Abort()
//			return
//		}
//
//		c.Set("UserId", mc.UserId)
//
//		c.Next()
//	}
//
//}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort()
			return
		}
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort()
			return
		}
		// 登录的用户
		c.Set("claims", claims)
		c.Next()
	}
}

func JwtAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("token")
		if token == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort()
			return
		}
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
