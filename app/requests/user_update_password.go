package requests

import (
	"GoBlog/app/models/user"
	"github.com/thedevsaddam/govalidator"
	"net/url"
)

// ValidateUpdatePasswordForm 验证表单，返回 errs 长度等于零即通过
func ValidateUpdatePasswordForm(data user.User) url.Values {
	// 1. 定制认证规则
	rules := govalidator.MapData{
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"email":            []string{"required", "min:4", "max:30", "email", "exists:users,email"},
	}

	// 2. 定制错误信息
	message := govalidator.MapData{
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
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
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入的密码不匹配")
	}
	return errs
}
