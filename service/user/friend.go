package user

import (
	"IM-Server/global"
	"IM-Server/model/system"
)

// InsertFriend  将新用户插入数据库
func InsertFriend(friend system.Friend) error {
	result := global.DB.Create(&friend)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteFriend1  删除好友关系
func DeleteFriend1(userId int64, friendId int64) error {
	result := global.DB.Where("user_id = ? AND friend_id = ?", userId, friendId).Delete(&system.Friend{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteFriend(userId int64, friendId int64) error {
	result := global.DB.Model(&system.Friend{}).Where("user_id = ? AND friend_id = ?", userId, friendId).Update("status", 0)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CheckFriendship 检查好友关系是否存在
func CheckFriendship(userId int64, friendId int64) (bool, error) {
	var count int64
	result := global.DB.Model(&system.Friend{}).Where("user_id = ? AND friend_id = ? AND status = 1", userId, friendId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// GetFriendList 根据用户ID获取其好友列表
func GetFriendList(userId int64) ([]system.Friend, error) {
	var friendInfos []system.Friend
	result := global.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&friendInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	return friendInfos, nil
}

// GetFriendIDs 根据用户ID获取其好友列表中的friend_id
func GetFriendIDs(userId int64) ([]int64, error) {
	var friendIDs []int64
	result := global.DB.Table("im_user_friend").
		Select("friend_id"). // 只查询friend_id列
		Where("user_id = ? AND status = ?", userId, 1).
		Pluck("friend_id", &friendIDs) // 使用Pluck方法直接将查询结果填充到int64切片中

	if result.Error != nil {
		return nil, result.Error
	}
	return friendIDs, nil
}
