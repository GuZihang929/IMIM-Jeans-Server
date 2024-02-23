package user

import (
	"IM-Server/global"
	"IM-Server/model/system"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// InsertGroup  将新的群数据插入数据库
func InsertGroup(info system.Group) error {
	result := global.DB.Create(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IsUserInGroup 检查用户是否已经在群中
func IsUserInGroup(userID int64, groupID int64) bool {
	var existingUser system.Group
	result := global.DB.Where("user_id = ? AND group_id = ?", userID, groupID).First(&existingUser)
	return result.Error == nil
}

// IsGroupOwner 检查用户是否是群主
func IsGroupOwner(userID int64, groupID int64) (bool, error) {
	var groupMember system.Group
	result := global.DB.Where("user_id = ? AND group_id = ? AND identity = ?", userID, groupID, 1).First(&groupMember)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil // 表示用户不是该群组的群主
		}
		return false, result.Error // 查询过程中出现其他错误
	}

	return true, nil // 用户是该群组的群主
}

func UpdateGroup(groupId int64, userId int64) error {
	var group system.Group
	// 根据群号和用户ID查询群组信息
	result := global.DB.Where("group_id = ? AND user_id = ?", groupId, userId).First(&group)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("用户与群组对应关系不存在")
	}

	// 修改用户在群组中的身份标识为2（管理员）
	group.Identity = 2

	// 更新数据库记录
	err := global.DB.Where("group_id = ? AND user_id = ?", groupId, userId).Save(&group).Error
	if err != nil {
		return fmt.Errorf("更新群组信息失败: %v", err)
	}

	return nil
}

// DeleteGroup 解散群组
func DeleteGroup(groupId int64) error {
	result := global.DB.Where(" group_id = ?", groupId).Delete(&system.Group{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteGroupUser  群主或者管理删除群员
func DeleteGroupUser(groupId int64, userId int64) error {
	result := global.DB.Where(" group_id = ? and user_id = ?", groupId, userId).Delete(&system.Group{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func IsGroupGM(userID int64, groupID int64) (bool, error) {
	var groupMember system.Group

	// 查询群组成员信息
	result := global.DB.Where("user_id = ? AND group_id = ? AND identity > 0", userID, groupID).First(&groupMember)

	// 处理查询结果
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil // 用户不是该群组的群主或管理员
	} else if result.Error != nil {
		return false, result.Error // 查询过程中出现其他错误
	}

	return true, nil // 用户拥有管理权限
}
