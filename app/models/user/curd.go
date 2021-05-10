package user

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/model"
	"GoBlog/pkg/password"
	"GoBlog/pkg/types"
)

// Get 通过 ID 获取用户信息
func Get(uidStr string) (User, error) {
	var user User
	uid := types.StringToInt(uidStr)
	if err := model.DB.First(&user, uid).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// GetByEmail 通过 Email 获取用户信息
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		logger.LogError(err)
		return User{}, err
	}
	return user, nil
}

// HasUserByEmail 通过 Email 判断用户是否存在，存在返回 TRUE，不存在返回 FALSE
func HasUserByEmail(email string) bool {
	var user User
	var count int64
	model.DB.Where("email = ?", email).First(&user).Count(&count)
	return count != 0
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Update 更新文章，如果文章不存在则自动创建
func (user *User) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&user)
	if err = model.DB.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// ComparePassword 对比密码是否匹配
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

// GetStringID 用于获取字符串类型的ID
func (user *User) GetStringID() string {
	return types.Uint64ToString(user.ID)
}
