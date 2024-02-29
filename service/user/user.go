package user

import (
	"IM-Server/global"
	"IM-Server/model/system"
	"context"
	"log"
	"time"
)

type PublicService struct {
}

func (s *PublicService) Login(id int, pwd string) bool {
	rowsAffected := global.DB.Where("id = ? and password = ?", id, pwd).First(&system.User{}).RowsAffected
	if rowsAffected == 1 {
		return true
	}
	return false
}

// IsUserExists 检查数据库中是否存在具有给定账号（邮箱）的用户
func IsUserExists(account string) bool {
	var count int64
	global.DB.Model(&system.User{}).Where("account = ?", account).Count(&count)
	return count > 0
}

//检查数据库中是否存在账户

// GetUserByAccount 根据账号获取用户信息
func GetUserByAccount(account string) (*system.User, error) {
	var user system.User
	result := global.DB.Where("account = ?", account).Limit(1).Find(&user)
	if result.Error != nil {
		// 处理数据库查询错误
		return nil, result.Error
	}

	// 用户不存在的情况下，返回nil
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

// GetUserByUserID  根据用户ID获取用户信息
func GetUserByUserID(userID int64) (*system.User, error) {
	var user system.User
	result := global.DB.Where("	id = ?", userID).Limit(1).Find(&user)
	if result.Error != nil {
		// 处理数据库查询错误
		return nil, result.Error
	}

	// 用户不存在的情况下，返回nil
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

// InsertUser 将新用户插入数据库
func InsertUser(user system.User) error {
	result := global.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// StoreUserAndCodeInRedis 在 Redis 中存储用户信息和验证码
func StoreUserAndCodeInRedis(account, code string) error {
	// 假设 global.Redis 是 *redis.Client 的实例
	_, err := global.Redis.Set(context.Background(), "vCode:"+account, code, 5*time.Minute).Result()
	if err != nil {
		log.Println("在 Redis 中存储用户和验证码时出错:", err)
		return err
	}
	return nil
}

// VerifyVerificationCode 检查提供的验证码是否与存储在 Redis 中的验证码匹配
func VerifyVerificationCode(account, code string) bool {
	// 假设 global.Redis 是你的 Redis 客户端实例
	storedCode, err := global.Redis.Get(context.Background(), "vCode:"+account).Result()
	if err != nil {
		log.Println("从 Redis 中检索验证码时出错:", err)
		return false
	}
	// 将提供的验证码与存储的验证码进行比较
	return storedCode == code
}
