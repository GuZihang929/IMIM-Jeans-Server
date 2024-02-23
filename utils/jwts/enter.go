package jwts

import (
	"github.com/dgrijalva/jwt-go/v4"
)

// JwtPayLoad jwt中payload数据
type JwtPayLoad struct {
	Username string `json:"account"`   // 用户名就是账户名
	NickName string `json:"nick_name"` // 昵称
	UserID   int64  `json:"uid"`       // 用户id
	Avatar   string `json:"avatar"`    //用户头像
}

var MySecret []byte

type CustomClaims struct {
	JwtPayLoad
	jwt.StandardClaims
}
