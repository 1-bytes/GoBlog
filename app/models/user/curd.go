package user

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/model"
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
		return User{}, err
	}
	return user, nil
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// ComparePassword 对比密码是否匹配
func (user *User) ComparePassword(password string) bool {
	return user.Password == password
}

// GetStringID 用于获取字符串类型的ID
func (user *User) GetStringID() string {
	return types.Uint64ToString(user.ID)
}
