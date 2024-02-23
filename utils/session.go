package utils

import (
	"crypto/rand"
	"encoding/hex"
)

//type Session struct {
//	SessionId int    `json:"sessionId" gorm:"session_id"`
//	Token     string `json:"token" gorm:"token"`
//}

func SetSession(length int) (string, error) {
	// 计算需要生成的字节数
	byteLength := length / 2
	if length%2 != 0 {
		byteLength++
	}

	// 生成随机字节序列
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将随机字节序列转换为十六进制字符串
	session := hex.EncodeToString(randomBytes)[:length]

	return session, nil
}
