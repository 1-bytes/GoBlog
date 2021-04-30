package requests

import (
	"GoBlog/app/models/user"
	"github.com/thedevsaddam/govalidator"
	"net/url"
)

// ValidateLostPasswordForm 验证表单，返回 errs 长度等于零即通过
func ValidateLostPasswordForm(data user.User) url.Values {
	// 1. 定制认证规则
	rules := govalidator.MapData{
		"email": []string{"required", "min:4", "max:30", "email", "exists:users,email"},
	}

	// 2. 定制错误信息
	message := govalidator.MapData{
		"email": []string{
			"required:邮箱为必填项",
			"min:邮箱长度需大于 4",
			"max:邮箱长度需小于 30",
			"email:邮箱格式不正确，请提供有效的邮箱地址",
		},
	}

	// 3. 配初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      message,
	}

	// 4. 开始认证
	errs := govalidator.New(opts).ValidateStruct()
	return errs
}
